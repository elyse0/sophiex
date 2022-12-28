package main

import (
	"io"
	"sophiex/internal/downloader"
	"sophiex/internal/output"
	"sync"
)

func DownloadSingleHlsUrlToPlayer(url string) {
	downloadManager := &sync.WaitGroup{}

	pipeReader, pipeWriter := io.Pipe()

	hlsDownloader := downloader.CreateHlsDownloader(url, pipeWriter)
	downloadManager.Add(1)

	go hlsDownloader.Download(downloadManager)

	playerDone := make(chan bool)
	player := output.CreateStreamPlayer()
	player.PlayFrom(pipeReader, playerDone)

	downloadManager.Wait()

	<-playerDone
}

func DownloadSingleHttpUrlToPlayer(url string) {
	var httpService = downloader.CreateHttpService()

	response, _ := httpService.Get(url, downloader.HttpRequestConfig{
		Headers: map[string]string{
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0",
			"Accept":          "*/*",
			"Accept-Encoding": "gzip, deflate, br",
			"Connection":      "keep-alive",
		},
	})

	playerDone := make(chan bool)
	player := output.CreateStreamPlayer()
	player.PlayFrom(response.Body, playerDone)

	<-playerDone
}

func DownloadMultipleHttpUrlsToPlayer(urls []string) {
	var httpService = downloader.CreateHttpService()

	downloadManager := &sync.WaitGroup{}

	var namedPipes []*output.StreamOutput
	for _, url := range urls {
		namedPipe := output.CreateNamedPipe()
		namedPipe.Open()

		namedPipes = append(namedPipes, namedPipe)

		downloadManager.Add(1)
		_url := url
		go func() {
			response, _ := httpService.Get(_url, downloader.HttpRequestConfig{
				Headers: map[string]string{
					"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0",
					"Accept":          "*/*",
					"Accept-Encoding": "gzip, deflate, br",
					"Connection":      "keep-alive",
				},
			})

			io.Copy(namedPipe.Stream, response.Body)

			downloadManager.Done()
		}()

	}

	pipeReader, pipeWriter := io.Pipe()

	var inputs []output.OsPath
	for _, namedPipe := range namedPipes {
		inputs = append(inputs, namedPipe)
	}

	muxerDone := make(chan bool)
	muxer := output.CreateMuxer()
	muxer.WriteTo(inputs, pipeWriter, muxerDone)

	playerDone := make(chan bool)
	player := output.CreateStreamPlayer()
	player.PlayFrom(pipeReader, playerDone)

	downloadManager.Wait()

	for _, namedPipe := range namedPipes {
		namedPipe.Close()
	}

	<-muxerDone
	<-playerDone
}

func DownloadMultipleHlsUrlsToPlayer(urls []string) {
	downloadManager := &sync.WaitGroup{}

	var namedPipes []*output.StreamOutput
	for _, url := range urls {
		namedPipe := output.CreateNamedPipe()
		namedPipe.Open()

		namedPipes = append(namedPipes, namedPipe)

		hlsDownloader := downloader.CreateHlsDownloader(url, namedPipe.Stream)
		downloadManager.Add(1)
		go hlsDownloader.Download(downloadManager)
	}

	pipeReader, pipeWriter := io.Pipe()

	var inputs []output.OsPath
	for _, namedPipe := range namedPipes {
		inputs = append(inputs, namedPipe)
	}

	muxerDone := make(chan bool)
	muxer := output.CreateMuxer()
	muxer.WriteTo(inputs, pipeWriter, muxerDone)

	playerDone := make(chan bool)
	player := output.CreateStreamPlayer()
	player.PlayFrom(pipeReader, playerDone)

	downloadManager.Wait()

	for _, namedPipe := range namedPipes {
		namedPipe.Close()
	}

	<-muxerDone
	<-playerDone
}

func main() {
	// hlsUrls := []string{
	// "http://localhost:8000/stream_3/stream_3.m3u8",
	// 	"http://localhost:8000/stream_0/stream_0.m3u8",
	// }

	// DownloadSingleHlsUrlToPlayer("http://localhost:8000/stream_3/stream_3.m3u8")
	// DownloadMultipleHlsUrlsToPlayer([]string{
	// 	"http://localhost:8000/stream_3/stream_3.m3u8",
	// 	"http://localhost:8000/stream_0/stream_0.m3u8",
	// })

	// DownloadSingleHttpUrlToPlayer("http://localhost:8000/test-audio.mp4")
	DownloadMultipleHttpUrlsToPlayer([]string{
		"http://localhost:8000/test-audio.mp4",
		"http://localhost:8000/test-video.mp4",
	})
}
