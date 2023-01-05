package main

import (
	"io"
	"sophiex/internal/downloader"
	"sophiex/internal/downloader/http"
	"sophiex/internal/output"
	"sophiex/internal/sites_extractor"
	"sync"
)

func DownloadFormat(format sites_extractor.DownloadableFormat, output io.WriteCloser, downloadManager *sync.WaitGroup) {
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

func DownloadSingleFormatToPlayer(format sites_extractor.DownloadableFormat, manager *sync.WaitGroup) io.Reader {
	pipeReader, pipeWriter := io.Pipe()

	DownloadFormat(format, pipeWriter, manager)

	return pipeReader
}

func DownloadMultipleFormatsToPlayer(formats []sites_extractor.DownloadableFormat, downloadManager *sync.WaitGroup) io.Reader {
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

	muxer := output.CreateMuxer()
	muxer.WriteTo(inputs, pipeWriter, downloadManager)

	return pipeReader
}

func DownloadFormatsToPlayer(formats []sites_extractor.DownloadableFormat, manager *sync.WaitGroup) io.Reader {
	if len(formats) == 1 {
		return DownloadSingleFormatToPlayer(formats[0], manager)
	} else {
		return DownloadMultipleFormatsToPlayer(formats, manager)
	}
}

func main() {
	manager := sync.WaitGroup{}

	downloadableFormats := sites_extractor.GetDownloadableFormats(
		"http://localhost:8080/master.m3u8")

	pipeReader := DownloadFormatsToPlayer(downloadableFormats, &manager)

	// outputFile, _ := os.Create("skam-test.mkv")
	// io.Copy(outputFile, pipeReader)

	player := output.CreateStreamPlayer()
	player.PlayFrom(pipeReader)

	manager.Wait()
}
