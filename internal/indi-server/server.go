package indiserver

import (
	"log"
	"log/slog"
	"os/exec"
	"sync"
	"syscall"
)

var (
	cmd     *exec.Cmd
	cmdLock sync.Mutex
)

func Start(drivers []string) error {
	if IsRunning() {
		return nil
	}

	slog.Info("Starting INDI server", "drivers", drivers)

	cmdLock.Lock()
	defer cmdLock.Unlock()

	cmd = exec.Command("indiserver", drivers...)
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

	slog.Info("Stopped INDI server")
	return nil
}

func IsRunning() bool {
	return cmd != nil && cmd.Process != nil
}
