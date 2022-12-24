package main

import (
	"fmt"
	"sophiex/internal/logger"
	"sophiex/internal/output"
	"sync"
)

func getAudioUrls() []string {
	var urls []string
	for i := 0; i <= 650; i++ {
		urls = append(urls, fmt.Sprintf("http://localhost:8000/stream_3/data%02d.ts", i))
	}

	logger.Log.Debug("%d audio urls", len(urls))
	return urls
}

func getVideoUrls() []string {
	var urls []string
	for i := 0; i <= 649; i++ {
		urls = append(urls, fmt.Sprintf("http://localhost:8000/stream_0/data%02d.ts", i))
	}

	logger.Log.Debug("%d audio urls", len(urls))
	return urls
}

func main() {
	streams := []*output.StreamDownloader{
		output.CreateHlsStreamDownloader(getAudioUrls(), output.CreateNamedPipe()),
		output.CreateHlsStreamDownloader(getVideoUrls(), output.CreateNamedPipe()),
	}

	// streamPlayer := output.CreateStreamPlayer()
	// streamPlayer.Launch([]*output.NamedPipe{audioOutput, videoOutput})

	outputWriter := output.CreateMuxer()
	outputWriter.Launch(streams, "muxed.mkv")
	//outputWriter := output.CreateStreamPlayer()
	//outputWriter.Launch(streams)

	streamManager := &sync.WaitGroup{}
	for _, stream := range streams {
		streamManager.Add(1)
		go stream.Downloader.Download(streamManager)
	}

	streamManager.Wait()

	for _, stream := range streams {
		stream.Output.Close()
	}
}
