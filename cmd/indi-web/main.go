package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/sys/unix"
)

var (
	cmd     *exec.Cmd
	cmdLock sync.Mutex
)

func Index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	body, _ := os.ReadFile("web/template/index.html")
	fmt.Fprint(w, string(body))
}

func INDIServer(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	cmdLock.Lock()
	defer cmdLock.Unlock()

	if cmd != nil && cmd.Process != nil {
		err := cmd.Process.Signal(syscall.SIGTERM)
		if err != nil {
			log.Printf("An error occurred while terminating INDI server: %v\n", err)
			return
		}

		_ = cmd.Wait()
		cmd = nil

		fmt.Fprintln(w, "Stopped")
	} else {
		fifoFile := "/tmp/indififo"
		unix.Mkfifo(fifoFile, 0o600)

		args := []string{"-vvv", "-p", "4000", "-f", fifoFile, "-r", "0"}

		cmd = exec.Command("indiserver", args...)
		cmd.Stdout = log.Writer()
		cmd.Stderr = log.Writer()

		if err := cmd.Start(); err != nil {
			log.Printf("An error occurred while starting INDI server: %v\n", err)
			return
		}

		fmt.Fprint(w, "Running")
	}
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/indi/activate", INDIServer)
	router.ServeFiles("/static/*filepath", http.Dir("web/static"))

	log.Fatal(http.ListenAndServe(":8080", router))
}
