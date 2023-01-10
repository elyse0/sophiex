package downloader

import (
	"fmt"
	"io"
	"net/http"
	"sophiex/internal/crypto"
	sophiexHttp "sophiex/internal/downloader/http"
	"sophiex/internal/logger"
	"sophiex/internal/output"
	"sophiex/internal/parser"
	"sophiex/internal/utils"
	"sync"
)

type Request struct {
	Index    int
	Fragment parser.HlsFragment
}

type Response struct {
	Fragment parser.HlsFragment
	Response *http.Response
}

type WorkerPool struct {
	manager   sync.WaitGroup
	requests  chan Request
	responses chan utils.OrderedFragment[Response]
}

func (workerPool *WorkerPool) initialize(fragments []parser.HlsFragment) {
	for index, _fragment := range fragments {
		request := Request{
			Index:    index,
			Fragment: _fragment,
		}
		workerPool.requests <- request
	}
	close(workerPool.requests)
}

func (workerPool *WorkerPool) worker() {
	for request := range workerPool.requests {
		logger.Log.Debug("Request url: %s\n", request.Fragment.Url)
		headers := map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0",
		}

		if !request.Fragment.ByteRange.IsEmpty() {
			headers["Range"] = fmt.Sprintf(
				"bytes=%d-%d",
				request.Fragment.ByteRange.Start+1,
				request.Fragment.ByteRange.End)
		}

		response, err := sophiexHttp.HttpService.Get(request.Fragment.Url, sophiexHttp.RequestConfig{
			Headers: headers,
		})
		logger.Log.Debug("Response url: %s\n", request.Fragment.Url)
		if err != nil {
			panic("Http error")
		}

		fragmentResponse := utils.OrderedFragment[Response]{
			Index: request.Index,
			Payload: Response{
				Fragment: request.Fragment,
				Response: response,
			},
		}

		workerPool.responses <- fragmentResponse
	}
	workerPool.manager.Done()
}

func (workerPool *WorkerPool) run(numberOfWorkers int) {
	for i := 0; i < numberOfWorkers; i++ {
		logger.Log.Debug("Creating worker no. %d\n", i)
		workerPool.manager.Add(1)
		go workerPool.worker()
	}
	workerPool.manager.Wait()
	close(workerPool.responses)
}

type HlsDownloader struct {
	initialization parser.HlsInitialization
	fragments      []parser.HlsFragment
	output         output.StreamWriter
}

func CreateHlsDownloader(manifestUrl string, stream output.StreamWriter) *HlsDownloader {
	response, _ := sophiexHttp.HttpService.Get(manifestUrl, sophiexHttp.RequestConfig{
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0",
		},
	})
	manifest, _ := io.ReadAll(response.Body)

	hlsMediaManifest := parser.HlsMediaManifest{
		ManifestUrl:  manifestUrl,
		Manifest:     string(manifest),
		IsLivestream: false,
	}

	parsedManifest, _ := hlsMediaManifest.Parse()

	// fmt.Println(parsedManifest.Fragments)

	return &HlsDownloader{
		initialization: parsedManifest.Initialization,
		fragments:      parsedManifest.Fragments,
		output:         stream,
	}
}

func (downloader *HlsDownloader) Download(streamManager *sync.WaitGroup) {
	var workerPool = WorkerPool{
		requests:  make(chan Request, 10),
		responses: make(chan utils.OrderedFragment[Response], 10),
	}

	go workerPool.initialize(downloader.fragments)
	go workerPool.run(4)

	fragmentOrderedQueue := utils.CreateFragmentOrderedQueue[Response](len(downloader.fragments))

	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0",
	}

	if !downloader.initialization.IsEmpty() {
		if !downloader.initialization.ByteRange.IsEmpty() {
			headers["Range"] = fmt.Sprintf(
				"bytes=%d-%d",
				downloader.initialization.ByteRange.Start,
				downloader.initialization.ByteRange.End)
		}

		initializationResponse, _ := sophiexHttp.HttpService.Get(
			downloader.initialization.Url,
			sophiexHttp.RequestConfig{
				Headers: headers,
			})
		initialization, _ := io.ReadAll(initializationResponse.Body)
		downloader.output.Write(initialization)
	}

	for response := range workerPool.responses {
		fragmentOrderedQueue.Enqueue(response)

		dequeueFragments, hasFinished := fragmentOrderedQueue.Dequeue()
		for _, dequeueFragment := range dequeueFragments {
			_fragment := dequeueFragment

			fragmentContent, _ := io.ReadAll(_fragment.Payload.Response.Body)
			_fragment.Payload.Response.Body.Close()

			if !_fragment.Payload.Fragment.Decryption.IsEmpty() {
				keyResponse, _ := sophiexHttp.HttpService.Get(
					_fragment.Payload.Fragment.Decryption.Uri,
					sophiexHttp.RequestConfig{})
				key, _ := io.ReadAll(keyResponse.Body)

				fragmentContent, _ = crypto.AesDecrypt(
					fragmentContent,
					key,
					_fragment.Payload.Fragment.Decryption.IV,
				)
			}

			downloader.output.Write(fragmentContent)

			// io.Copy(downloader.output, fragmentContent)
			// downloader.output.PlayFrom(dequeueFragment.Response.Body)
		}

		if hasFinished {
			break
		}
	}

	downloader.output.Close()
	streamManager.Done()
}
