package downloader

import (
	"fmt"
	"net/http"
	"time"
)

type HttpService struct {
	client *http.Client
}

func createHttpService() *HttpService {
	transport := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	httpService := HttpService{
		client: &http.Client{
			Transport: transport,
		},
	}

	return &httpService
}

func (httpService *HttpService) get(url string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	response, err := httpService.client.Do(request)
	if err != nil {
		fmt.Println(response)
		return nil, err
	}

	return response, nil
}
