package output

import (
	"os"
	"sophiex/internal/downloader"
	"sync"
)

type StreamOutput struct {
	Name   string
	Path   string
	Stream *os.File
}

type Downloader interface {
	Download(downloadManager *sync.WaitGroup)
}

type StreamDownloader struct {
	Downloader Downloader
	Output     *StreamOutput
}

func CreateHlsStreamDownloader(urls []string, outputStream *StreamOutput) *StreamDownloader {
	outputStream.Open()
	return &StreamDownloader{
		Downloader: downloader.CreateHlsDownloader(urls, outputStream.Stream),
		Output:     outputStream,
	}
}

type StreamPlayer interface {
	Launch(streams []*StreamDownloader)
}

type StreamWriter interface {
	Launch(streams []*StreamDownloader, outputPath string)
}
