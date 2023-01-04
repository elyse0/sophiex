package main

import (
	"io"
	"sophiex/internal/downloader"
	"sophiex/internal/downloader/http"
	"sophiex/internal/output"
	"sophiex/internal/sites_extractor"
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
	var httpService = http.CreateHttpService()

	response, _ := httpService.Get(url, http.HttpRequestConfig{
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
	var httpService = http.CreateHttpService()

	downloadManager := &sync.WaitGroup{}

	var namedPipes []*output.StreamOutput
	for _, url := range urls {
		namedPipe := output.CreateNamedPipe()
		namedPipe.Open()

		namedPipes = append(namedPipes, namedPipe)

		downloadManager.Add(1)
		_url := url
		go func() {
			response, _ := httpService.Get(_url, http.HttpRequestConfig{
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

func DownloadMultipleHttpUrlsUsingMultipartToPlayer() {
	downloadableFormats := sites_extractor.GetDownloadableFormats(
		"https://www.youtube.com/watch?v=SHO3dE-IDWk")

	var urls []string
	for _, format := range downloadableFormats {
		if format.Protocol != sites_extractor.Http {
			panic(format)
		}

		urls = append(urls, format.Url)
	}

	downloadManager := &sync.WaitGroup{}
	httpService := http.CreateHttpService()

	var namedPipes []*output.StreamOutput
	for _, url := range urls {
		namedPipe := output.CreateNamedPipe()
		namedPipe.Open()

		namedPipes = append(namedPipes, namedPipe)

		downloadManager.Add(1)

		_url := url
		go func() {
			_, _ = httpService.GetMultiFragment(
				// "http://localhost:8080/skam-france.mp4",
				_url,
				http.HttpRequestConfig{},
				namedPipe.Stream,
				downloadManager,
			)
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

	<-playerDone
}

func DownloadFormat(format sites_extractor.DownloadableFormat, output io.Writer, downloadManager *sync.WaitGroup) {
	httpService := http.CreateHttpService()

	downloadManager.Add(1)

	switch format.Protocol {
	case sites_extractor.Http:
		_url := format.Url
		go func() {
			_, _ = httpService.GetMultiFragment(
				_url,
				http.HttpRequestConfig{},
				output,
				downloadManager,
			)
		}()
	case sites_extractor.Hls:
		_url := format.Url
		hlsDownloader := downloader.CreateHlsDownloader(_url, output)

		go hlsDownloader.Download(downloadManager)
	default:
		panic("Unhandled format")
	}
}

func DownloadSingleFormatToPlayer(format sites_extractor.DownloadableFormat) {
	downloadManager := &sync.WaitGroup{}

	pipeReader, pipeWriter := io.Pipe()

	DownloadFormat(format, pipeWriter, downloadManager)

	playerDone := make(chan bool)
	player := output.CreateStreamPlayer()
	player.PlayFrom(pipeReader, playerDone)

	downloadManager.Wait()

	<-playerDone
}

func DownloadMultipleFormatsToPlayer(formats []sites_extractor.DownloadableFormat) {
	downloadManager := &sync.WaitGroup{}

	var namedPipes []*output.StreamOutput
	for _, format := range formats {
		namedPipe := output.CreateNamedPipe()
		namedPipe.Open()
		namedPipes = append(namedPipes, namedPipe)

		DownloadFormat(format, namedPipe.Stream, downloadManager)
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

func DownloadFormatsToPlayer(formats []sites_extractor.DownloadableFormat) {
	if len(formats) == 1 {
		DownloadSingleFormatToPlayer(formats[0])
	} else {
		DownloadMultipleFormatsToPlayer(formats)
	}
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
	// DownloadMultipleHttpUrlsToPlayer([]string{
	// 	"http://localhost:8000/test-audio.mp4",
	// 	"http://localhost:8000/test-video.mp4",
	// })

	// downloadableFormats := sites_extractor.GetDownloadableFormats(
	// 	"https://www.youtube.com/watch?v=dQcHxnIeS5w")

	downloadableFormats := sites_extractor.GetDownloadableFormats(
		"https://www.youtube.com/watch?v=zFt0tO4Op14")

	// downloadableFormats := []sites_extractor.DownloadableFormat{
	// 	{
	// 		Id:       "arte",
	// 		Url:      "https://pdvideosdaserste-a.akamaihd.net/int/2022/04/25/e6d7e26d-69b8-4700-8cb4-588072fe6d1f/JOB_92432_sendeton_1920x1080-50p-5000kbit.mp4",
	// 		Protocol: sites_extractor.Http,
	// 	},
	// }

	// downloadableFormats := []sites_extractor.DownloadableFormat{
	// 	{
	// 		Id:       "video",
	// 		Url:      "https://cloudingest.ftven.fr/2b43723d56716/833b322d-321e-4201-a43d-c57f6e763eb6_france-domtom_TA.ism/ZXhwPTE2NzI0OTc0MzJ+YWNsPSUyZjJiNDM3MjNkNTY3MTYlMmY4MzNiMzIyZC0zMjFlLTQyMDEtYTQzZC1jNTdmNmU3NjNlYjZfZnJhbmNlLWRvbXRvbV9UQS5pc20qfmhtYWM9YWUzM2ZiNmM5ZTk5NmU4MGI4NTQyNjRiZDljMDEyODZmYjM5N2NkZDlkNDllMTliZTIwZDE0M2M1NDhjNDYwZQ==/833b322d-321e-4201-a43d-c57f6e763eb6_france-domtom_TA-video=2000000.m3u8",
	// 		Protocol: sites_extractor.Hls,
	// 	},
	// 	{
	// 		Id:       "audio",
	// 		Url:      "https://cloudingest.ftven.fr/2b43723d56716/833b322d-321e-4201-a43d-c57f6e763eb6_france-domtom_TA.ism/ZXhwPTE2NzI0OTc0MzJ+YWNsPSUyZjJiNDM3MjNkNTY3MTYlMmY4MzNiMzIyZC0zMjFlLTQyMDEtYTQzZC1jNTdmNmU3NjNlYjZfZnJhbmNlLWRvbXRvbV9UQS5pc20qfmhtYWM9YWUzM2ZiNmM5ZTk5NmU4MGI4NTQyNjRiZDljMDEyODZmYjM5N2NkZDlkNDllMTliZTIwZDE0M2M1NDhjNDYwZQ==/833b322d-321e-4201-a43d-c57f6e763eb6_france-domtom_TA-audio_fre=96000.m3u8",
	// 		Protocol: sites_extractor.Hls,
	// 	},
	// }

	DownloadFormatsToPlayer(downloadableFormats)

	// if len(downloadableFormats) == 1 {
	// 	DownloadSingleHttpUrlToPlayer(downloadableFormats[0].Url)
	// } else {
	// 	var urls []string
	// 	for _, format := range downloadableFormats {
	// 		urls = append(urls, format.Url)
	// 	}
	// 		DownloadMultipleHttpUrlsToPlayer(urls)
	// }

	// f, err := os.Create("test.mp4")
	// if err != nil {
	// 	 panic(err)
	// }

	// DownloadMultipleHttpUrlsUsingMultipartToPlayer()
}
