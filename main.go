package main

import (
	"io"
	"sophiex/internal/downloader"
	"sophiex/internal/downloader/http"
	"sophiex/internal/output"
	"sophiex/internal/sites_extractor"
	"sync"
)

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

	player := output.CreateStreamPlayer()
	player.PlayFrom(pipeReader)

	downloadManager.Wait()
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

	player := output.CreateStreamPlayer()
	player.PlayFrom(pipeReader)

	downloadManager.Wait()

	for _, namedPipe := range namedPipes {
		namedPipe.Close()
	}

	<-muxerDone
}

func DownloadFormatsToPlayer(formats []sites_extractor.DownloadableFormat) {
	if len(formats) == 1 {
		DownloadSingleFormatToPlayer(formats[0])
	} else {
		DownloadMultipleFormatsToPlayer(formats)
	}
}

func main() {
	downloadableFormats := sites_extractor.GetDownloadableFormats(
		"http://localhost:8080/master.m3u8")

	DownloadFormatsToPlayer(downloadableFormats)
}
