package output

import (
	"fmt"
	"io"
	"os/exec"
)

type MpvPlayer struct {
	bin string
}

func CreateStreamPlayer() *MpvPlayer {
	return &MpvPlayer{
		bin: "mpv",
	}
}

func (player *MpvPlayer) PlayFrom(input io.Reader, done chan bool) {
	playerCommand := exec.Command(player.bin, "-", "--player-operation-mode=pseudo-gui")
	playerCommand.Stdin = input

	go func() {
		err := playerCommand.Run()
		if err != nil {
			panic(err)
		}
		close(done)
	}()
}

func (player *MpvPlayer) Launch(streams []*StreamDownloader, done chan bool) {
	var playerCommand *exec.Cmd
	if len(streams) == 1 {
		// playerCommand = exec.Command(player.bin, streams[0].Output.Path)
		playerCommand = exec.Command(player.bin, "FIXME")

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
	// FIXME
	muxer.WriteTo(nil, playerStdin, done)

	go func() {
		err := playerCommand.Run()
		if err != nil {
			fmt.Print(err)
		}
	}()
}
