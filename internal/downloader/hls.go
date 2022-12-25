package downloader

import (
	"os"
	"sophiex/internal/downloader/fragment"
	"sophiex/internal/logger"
	"sophiex/internal/utils"
	"sync"
)

type WorkerPool struct {
	manager   sync.WaitGroup
	requests  chan fragment.FragmentRequest
	responses chan fragment.FragmentResponse
}

var httpService = createHttpService()

func (workerPool *WorkerPool) initialize(urls []string) {
	for index, url := range urls {
		request := fragment.FragmentRequest{
			Index: index,
			Url:   url,
		}
		workerPool.requests <- request
	}
	close(workerPool.requests)
}

func (workerPool *WorkerPool) worker() {
	for request := range workerPool.requests {
		logger.Log.Debug("Request url: %s\n", request.Url)
		response, err := httpService.get(request.Url, HttpRequestConfig{
			Headers: nil,
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
	urls   []string
	output *os.File
}

func CreateHlsDownloader(urls []string, output *os.File) *HlsDownloader {
	return &HlsDownloader{
		urls:   urls,
		output: output,
	}
}

func (downloader *HlsDownloader) Download(streamManager *sync.WaitGroup) {
	var workerPool = WorkerPool{
		requests:  make(chan fragment.FragmentRequest, 10),
		responses: make(chan fragment.FragmentResponse, 10),
	}

	go workerPool.initialize(downloader.urls)
	go workerPool.run(4)

	fragmentOrderedQueue := utils.CreateFragmentOrderedQueue(len(downloader.urls))

	for response := range workerPool.responses {
		fragmentOrderedQueue.Enqueue(response)

		dequeueFragments, hasFinished := fragmentOrderedQueue.Dequeue()
		for _, dequeueFragment := range dequeueFragments {
			downloader.output.ReadFrom(dequeueFragment.Response.Body)
			dequeueFragment.Response.Body.Close()
		}

		if hasFinished {
			break
		}
	}

	streamManager.Done()
}
