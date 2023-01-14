package fragment

import (
	"fmt"
	"net/http"
	sophiexHttp "sophiex/internal/downloader/http"
	"sophiex/internal/fragment"
	"sophiex/internal/logger"
	"sophiex/internal/ordered-queue"
	"sync"
)

type Request struct {
	Index    int
	Fragment fragment.Generic
}

type Response struct {
	Fragment fragment.Generic
	Response *http.Response
}

type Queue struct {
	manager   sync.WaitGroup
	Requests  chan Request
	Responses chan ordered_queue.OrderedItem[Response]
}

func (queue *Queue) Initialize(fragments []fragment.Generic) {
	for index, _fragment := range fragments {
		request := Request{
			Index:    index,
			Fragment: _fragment,
		}
		queue.Requests <- request
	}
	close(queue.Requests)
}

func (queue *Queue) worker() {
	for request := range queue.Requests {
		headers := map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0",
		}

		byteRange := request.Fragment.ByteRange
		if !byteRange.IsEmpty() {
			headers["Range"] = fmt.Sprintf("bytes=%d-%d", byteRange.Start+1, byteRange.End)
		}

		url := request.Fragment.Url
		logger.Log.Debug("Request url: %s\n", url)
		response, err := sophiexHttp.HttpService.Get(url, sophiexHttp.RequestConfig{Headers: headers})
		logger.Log.Debug("Response url: %s\n", url)
		if err != nil {
			panic("Http error")
		}

		fragmentResponse := ordered_queue.OrderedItem[Response]{
			Index: request.Index,
			Payload: Response{
				Fragment: request.Fragment,
				Response: response,
			},
		}

		queue.Responses <- fragmentResponse
	}
	queue.manager.Done()
}

func (queue *Queue) Run(numberOfWorkers int) {
	for i := 0; i < numberOfWorkers; i++ {
		logger.Log.Debug("Creating worker no. %d\n", i)
		queue.manager.Add(1)
		go queue.worker()
	}
	queue.manager.Wait()
	close(queue.Responses)
}
