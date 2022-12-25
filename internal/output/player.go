package output

import (
	"fmt"
	"os/exec"
)

type MpvPlayer struct {
	bin string
}

func CreateStreamPlayer() StreamPlayer {
	return &MpvPlayer{
		bin: "mpv",
	}
}

func (player *MpvPlayer) Launch(streams []*StreamDownloader, done chan bool) {
	var playerCommand *exec.Cmd
	if len(streams) == 1 {
		playerCommand = exec.Command(player.bin, streams[0].Output.Path)

		go func() {
			err := playerCommand.Run()
			if err != nil {
				fmt.Print(err)
			}
		}()

		close(done)
		return
	}

	playerCommand = exec.Command(player.bin, "-")
	playerStdin, _ := playerCommand.StdinPipe()

	muxer := CreateMuxer()
	muxer.WriteTo(streams, playerStdin, done)

	go func() {
		err := playerCommand.Run()
		if err != nil {
			fmt.Print(err)
		}
	}()
}
