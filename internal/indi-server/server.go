package indiserver

import (
	"log"
	"os/exec"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

var (
	cmd     *exec.Cmd
	cmdLock sync.Mutex
)

func Start() error {
	if IsRunning() {
		return nil
	}

	cmdLock.Lock()
	defer cmdLock.Unlock()

	fifoFile := "/tmp/indififo"
	unix.Mkfifo(fifoFile, 0o600)

	args := []string{"-vvv", "-f", fifoFile, "-r", "0"}

	cmd = exec.Command("indiserver", args...)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	return cmd.Start()
}

func Stop() error {
	if !IsRunning() {
		return nil
	}

	cmdLock.Lock()
	defer cmdLock.Unlock()

	err := cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}

	_ = cmd.Wait()
	cmd = nil

	log.Println("Stopped INDI server")
	return nil
}

func IsRunning() bool {
	return cmd != nil && cmd.Process != nil
}
