package output

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"sophiex/internal/logger"
	"syscall"
)

func CreateNamedPipe() *StreamOutput {
	name := fmt.Sprintf("sophiex-%d-%s", os.Getpid(), uuid.NewString())
	path := filepath.Join(os.TempDir(), name)

	err := syscall.Mkfifo(path, 0640)
	if err != nil {
		panic("Couldn't create NamedPipe")
	}

	return &StreamOutput{
		Name:   name,
		Path:   path,
		Stream: nil,
	}
}

func (namedPipe *StreamOutput) Open() {
	stream, err := os.OpenFile(namedPipe.Path, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		panic("Couldn't open NamedPipe")
	}

	namedPipe.Stream = stream
}

func (namedPipe *StreamOutput) Close() {
	logger.Log.Debug("Closing pipe")
	if namedPipe.Stream == nil {
		return
	}

	err := namedPipe.Stream.Close()
	if err != nil {
		panic(err)
		return
	}

	err = os.Remove(namedPipe.Path)
	if err != nil {
		panic("Couldn't close NamedPipe")
	}

	namedPipe.Stream = nil
}
