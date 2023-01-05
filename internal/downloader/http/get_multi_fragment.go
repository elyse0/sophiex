package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"sophiex/internal/logger"
	"sophiex/internal/output"
	"sophiex/internal/utils"
	"strconv"
	"sync"
)

type Range struct {
	Start int
	End   int
}

type RangeRequest struct {
	Index int
	Url   string
	Range Range
}

type WorkerPool struct {
	manager   sync.WaitGroup
	requests  chan RangeRequest
	responses chan utils.OrderedFragment[*http.Response]
}

var httpService = CreateHttpService()

func (workerPool *WorkerPool) initialize(url string, ranges []Range) {
	for index, _range := range ranges {
		request := RangeRequest{
			Index: index,
			Url:   url,
			Range: _range,
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
				"Range":           fmt.Sprintf("bytes=%d-%d", request.Range.Start, request.Range.End),
			},
		})
		logger.Log.Debug("Response url: %s\n", request.Url)
		if err != nil {
			panic("Http error")
		}

		rangeResponse := utils.OrderedFragment[*http.Response]{
			Index:   request.Index,
			Payload: response,
		}

		workerPool.responses <- rangeResponse
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

func (httpService *HttpService) GetMultiFragment(
	url string,
	config HttpRequestConfig,
	output output.StreamWriter,
	streamManager *sync.WaitGroup,
) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		panic(err)
	}

	addRequestHeaders(request, config.Headers)

	response, err := httpService.client.Do(request)
	if err != nil {
		return nil, err
	}
	response.Body.Close()

	if response.Header.Get("Accept-Ranges") != "bytes" {
		return nil, errors.New("http: Server doesn't support Byte-Range")
	}

	contentLengthString := response.Header.Get("Content-Length")
	if contentLengthString == "" {
		return nil, errors.New("http: Couldn't read Content-Length")
	}

	contentLength, err := strconv.Atoi(contentLengthString)
	if err != nil {
		return nil, err
	}

	chunkSize := 4000000

	var ranges []Range
	for byteCount := 0; byteCount <= contentLength; byteCount += chunkSize {
		var _range Range
		if (byteCount + chunkSize) <= contentLength {
			_range = Range{
				Start: byteCount,
				End:   (byteCount + chunkSize) - 1,
			}
		} else {
			_range = Range{
				Start: byteCount,
				End:   contentLength + 1,
			}
		}

		ranges = append(ranges, _range)
	}

	// Real download

	var workerPool = WorkerPool{
		requests:  make(chan RangeRequest, 10),
		responses: make(chan utils.OrderedFragment[*http.Response], 10),
	}

	go workerPool.initialize(url, ranges)
	go workerPool.run(4)

	fragmentOrderedQueue := utils.CreateFragmentOrderedQueue[*http.Response](len(ranges))

	for response := range workerPool.responses {
		fragmentOrderedQueue.Enqueue(response)

		dequeueFragments, hasFinished := fragmentOrderedQueue.Dequeue()
		for _, dequeueFragment := range dequeueFragments {
			io.Copy(output, dequeueFragment.Payload.Body)
			// downloader.output.PlayFrom(dequeueFragment.Response.Body)
			dequeueFragment.Payload.Body.Close()
		}

		if hasFinished {
			break
		}
	}

	streamManager.Done()
	output.Close()

	return response, nil
}
