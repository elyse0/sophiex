package downloader

import (
	"fmt"
	"io"
	"sophiex/internal/downloader/fragment"
	"sophiex/internal/logger"
	"sophiex/internal/output"
	"sophiex/internal/parser"
	"sophiex/internal/utils"
	"sync"
)

type WorkerPool struct {
	manager   sync.WaitGroup
	requests  chan fragment.FragmentRequest
	responses chan fragment.FragmentResponse
}

var httpService = CreateHttpService()

func (workerPool *WorkerPool) initialize(fragments []parser.HlsFragment) {
	for index, _fragment := range fragments {
		request := fragment.FragmentRequest{
			Index: index,
			Url:   _fragment.Url,
		}
		workerPool.requests <- request
	}
	close(workerPool.requests)
}

func (workerPool *WorkerPool) worker() {
	for request := range workerPool.requests {
		logger.Log.Debug("Request url: %s\n", request.Url)
		response, err := httpService.Get(request.Url, HttpRequestConfig{
			Headers: map[string]string{
				"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0",
				"Accept":          "*/*",
				"Accept-Encoding": "gzip, deflate, br",
				"Connection":      "keep-alive",
			},
		})
		logger.Log.Debug("Response url: %s\n", request.Url)
		if err != nil {
			panic("Http error")
		}

		fragmentResponse := fragment.FragmentResponse{
			Index:    request.Index,
			Response: response,
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
	fragments []parser.HlsFragment
	output    output.StreamWriter
}

func CreateHlsDownloader(manifestUrl string, stream output.StreamWriter) *HlsDownloader {
	response, _ := httpService.Get(manifestUrl, HttpRequestConfig{
		Headers: map[string]string{
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0",
			"Accept":          "*/*",
			"Accept-Encoding": "gzip, deflate, br",
			"Connection":      "keep-alive",
		},
	})
	manifest, _ := io.ReadAll(response.Body)

	hlsMediaManifest := parser.HlsMediaManifest{
		ManifestUrl:  manifestUrl,
		Manifest:     string(manifest),
		IsLivestream: false,
	}

	parsedManifest, _ := hlsMediaManifest.Parse()

	return &HlsDownloader{
		fragments: parsedManifest.Fragments,
		output:    stream,
	}
}

func (downloader *HlsDownloader) Download(streamManager *sync.WaitGroup) {
	var workerPool = WorkerPool{
		requests:  make(chan fragment.FragmentRequest, 10),
		responses: make(chan fragment.FragmentResponse, 10),
	}

	go workerPool.initialize(downloader.fragments)
	go workerPool.run(4)

	fragmentOrderedQueue := utils.CreateFragmentOrderedQueue(len(downloader.fragments))

	for response := range workerPool.responses {
		fragmentOrderedQueue.Enqueue(response)

		dequeueFragments, hasFinished := fragmentOrderedQueue.Dequeue()
		for _, dequeueFragment := range dequeueFragments {
			io.Copy(downloader.output, dequeueFragment.Response.Body)
			// downloader.output.PlayFrom(dequeueFragment.Response.Body)
			dequeueFragment.Response.Body.Close()
		}

		if hasFinished {
			break
		}
	}

	streamManager.Done()
}
