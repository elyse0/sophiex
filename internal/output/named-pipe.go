package output

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"sophiex/internal/logger"
	"syscall"
)

func CreateNamedPipe() (*StreamOutput, error) {
	name := fmt.Sprintf("sophiex-%d-%s", os.Getpid(), uuid.NewString())
	path := filepath.Join(os.TempDir(), name)

	err := syscall.Mkfifo(path, 0640)
	if err != nil {
		return nil, err
	}

	return &StreamOutput{
		name:   name,
		path:   path,
		Stream: nil,
	}, nil
}

func (namedPipe *StreamOutput) Path() string {
	return namedPipe.path
}

func (namedPipe *StreamOutput) Open() error {
	stream, err := os.OpenFile(namedPipe.path, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		return err
	}

	namedPipe.Stream = stream
	return nil
}

func (namedPipe *StreamOutput) Close() error {
	logger.Log.Debug("Closing pipe")
	if namedPipe.Stream == nil {
		return nil
	}

	err := namedPipe.Stream.Close()
	if err != nil {
		return err
	}

	err = os.Remove(namedPipe.path)
	if err != nil {
		return err
	}

	namedPipe.Stream = nil
	return nil
}
