// +build !windows
package pipes

import (
	"io/ioutil"
	"path/filepath"
	"syscall"
)

type NamedPipe struct {
	P1BoundPipeName string
	P2BoundPipeName string
}

func NewNamedPipe() NamedPipe {
	tmpDir, _ := ioutil.TempDir("", "named-pipes")
	np := &NamedPipe{
		P1BoundPipeName: filepath.Join(tmpDir, "parentBound"),
		P2BoundPipeName: filepath.Join(tmpDir, "childBound"),
	}

	// Create named pipe

	syscall.Mkfifo(np.p1BoundPipeName, 0600)
	syscall.Mkfifo(np.p2BoundPipeName, 0600)
}
