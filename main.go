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
		output.CreateHlsStreamDownloader("http://localhost:8000/stream_3/stream_3.m3u8", output.CreateNamedPipe()),
		output.CreateHlsStreamDownloader("http://localhost:8000/stream_0/stream_0.m3u8", output.CreateNamedPipe()),
	}

	// streamPlayer := output.CreateStreamPlayer()
	// streamPlayer.Launch([]*output.NamedPipe{audioOutput, videoOutput})

	muxerDone := make(chan bool)

	outputWriter := output.CreateMuxer()
	outputWriter.WriteToFile(streams, "muxed.mkv", muxerDone)
	// outputWriter.WriteTo(streams, muxerDone)

	// outputWriter := output.CreateStreamPlayer()
	// outputWriter.Launch(streams, muxerDone)

	streamManager := &sync.WaitGroup{}
	for _, stream := range streams {
		streamManager.Add(1)
		go stream.Downloader.Download(streamManager)
	}

	streamManager.Wait()

	for _, stream := range streams {
		stream.Output.Close()
	}

	<-muxerDone
}
