package downloader

import (
	"fmt"
	"io"
	"sophiex/internal/crypto"
	"sophiex/internal/downloader/fragment"
	sophiexHttp "sophiex/internal/downloader/http"
	fragment2 "sophiex/internal/fragment"
	"sophiex/internal/ordered-queue"
	"sophiex/internal/output"
	"sophiex/internal/parser/hls"
	"sync"
)

type HlsDownloader struct {
	initialization hls.Initialization
	fragments      []hls.Fragment
	output         output.StreamWriter
}

func CreateHlsDownloader(manifestUrl string, stream output.StreamWriter) *HlsDownloader {
	response, _ := sophiexHttp.HttpService.Get(manifestUrl, sophiexHttp.RequestConfig{
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0",
		},
	})
	manifest, _ := io.ReadAll(response.Body)

	hlsMediaManifest := hls.MediaManifest{
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
	var queue = fragment.Queue{
		Requests:  make(chan fragment.Request, 10),
		Responses: make(chan ordered_queue.OrderedItem[fragment.Response], 10),
	}

	var genericFragments []fragment2.Generic
	for _, frag := range downloader.fragments {
		genericFragments = append(genericFragments, frag.Generic)
	}

	go queue.Initialize(genericFragments)
	go queue.Run(4)

	fragmentOrderedQueue := ordered_queue.CreateOrderedQueue[fragment.Response](len(downloader.fragments))

	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0",
	}

	initialization := downloader.initialization
	if !initialization.IsEmpty() {
		byteRange := initialization.ByteRange
		if !byteRange.IsEmpty() {
			headers["Range"] = fmt.Sprintf("bytes=%d-%d", byteRange.Start, byteRange.End)
		}

		initializationResponse, _ := sophiexHttp.HttpService.Get(
			initialization.Url,
			sophiexHttp.RequestConfig{
				Headers: headers,
			})
		initialization, _ := io.ReadAll(initializationResponse.Body)
		downloader.output.Write(initialization)
	}

	for response := range queue.Responses {
		fragmentOrderedQueue.Enqueue(response)

		dequeueFragments, hasFinished := fragmentOrderedQueue.Dequeue()
		for _, dequeueFragment := range dequeueFragments {
			_fragment := dequeueFragment

			fragmentContent, _ := io.ReadAll(_fragment.Payload.Response.Body)
			_fragment.Payload.Response.Body.Close()

			decryption := _fragment.Payload.Fragment.Decryption
			if !decryption.IsEmpty() {
				uri := decryption.Uri
				keyResponse, _ := sophiexHttp.HttpService.Get(uri, sophiexHttp.RequestConfig{})
				key, _ := io.ReadAll(keyResponse.Body)

				fragmentContent, _ = crypto.AesDecrypt(fragmentContent, key, decryption.IV)
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
