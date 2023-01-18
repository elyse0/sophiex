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
	httpService := http.CreateService()

	downloadManager.Add(1)

	switch format.Protocol {
	case sites_extractor.Http:
		_url := format.Url
		go func() {
			_, err := httpService.GetMultiFragment(
				_url,
				http.RequestConfig{},
				output,
				downloadManager,
			)
			if err != nil {
				panic(err)
			}
		}()
	case sites_extractor.Hls:
		_url := format.Url
		hlsDownloader := downloader.CreateHlsDownloader(_url, output)

		go hlsDownloader.Download(downloadManager)
	default:
		panic("Unhandled format")
	}
}

func DownloadSingleFormat(format sites_extractor.DownloadableFormat, manager *sync.WaitGroup) io.Reader {
	pipeReader, pipeWriter := io.Pipe()

	DownloadFormat(format, pipeWriter, manager)

	return pipeReader
}

func DownloadMultipleFormats(formats []sites_extractor.DownloadableFormat, downloadManager *sync.WaitGroup) io.Reader {
	var namedPipes []*output.StreamOutput
	for _, format := range formats {
		namedPipe, err := output.CreateNamedPipe()
		if err != nil {
			panic(err)
		}

		err = namedPipe.Open()
		if err != nil {
			panic(err)
		}

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

func DownloadFormats(formats []sites_extractor.DownloadableFormat, manager *sync.WaitGroup) io.Reader {
	if len(formats) == 1 {
		return DownloadSingleFormat(formats[0], manager)
	} else {
		return DownloadMultipleFormats(formats, manager)
	}
}

func main() {
	manager := sync.WaitGroup{}

	downloadableFormats, err := sites_extractor.GetDownloadableFormats(os.Args[1])
	if err != nil {
		panic(err)
	}

	pipeReader := DownloadFormats(downloadableFormats, &manager)

	// outputFile, _ := os.Create("skam-test.mkv")
	// io.Copy(outputFile, pipeReader)

	player := output.CreateStreamPlayer()
	player.PlayFrom(pipeReader)

	manager.Wait()
}
