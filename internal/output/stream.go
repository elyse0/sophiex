package output

import (
	"os"
	"sync"
)

type StreamWriter interface {
	Write(data []byte) (n int, err error)
}

type StreamOutput struct {
	name   string
	path   string
	Stream *os.File
}

type Downloader interface {
	Download(downloadManager *sync.WaitGroup)
}

type StreamDownloader struct {
	Downloader Downloader
	Output     StreamWriter
}

type StreamPlayer interface {
	Launch(streams []*StreamDownloader, done chan bool)
}
