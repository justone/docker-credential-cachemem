package main

import (
	"fmt"
	"os"
	"time"

	xdg "github.com/casimir/xdg-go"
	"github.com/docker/docker-credential-helpers/credentials"
	daemon "github.com/sevlyar/go-daemon"
)

var validCommands = map[string]bool{
	"daemon":  true,
	"stop":    true,
	"bump":    true,
	"store":   true,
	"get":     true,
	"erase":   true,
	"list":    true,
	"version": true,
}

func main() {
	mode := "grpc"

	cm := NewCacheMem(mode)

	if _, ok := validCommands[os.Args[1]]; !ok {
		fmt.Println("bad arg")
		return
	}

	ensureDaemon(cm)

	if len(os.Args) > 1 && os.Args[1] == "stop" {
		cm.Stop()
	} else if len(os.Args) > 1 && os.Args[1] == "bump" {
		cm.Bump()
	} else {
		credentials.Serve(CredHandler{cm})
	}
}

func ensureDaemon(cm *CacheMem) {

	app := xdg.App{Name: "cachemem"}

	os.MkdirAll(app.DataPath(""), 0755)
	cntxt := &daemon.Context{
		PidFileName: app.DataPath("daemon.pid"),
		PidFilePerm: 0644,
		LogFileName: app.DataPath("daemon.log"),
		LogFilePerm: 0640,
		WorkDir:     app.DataPath(""),
	}

	d, err := cntxt.Reborn()
	if err != nil {
		if err != daemon.ErrWouldBlock {
			fmt.Println(err)
		}
		return
	}
	if d != nil {
		// we're in the parent process, check if spun-off daemon is working yet
		for i := 0; i < 10; i++ {
			cl, _ := cm.Client()

			if _, err := cl.Send(Request{Command: "bump"}); err == nil {
				break
			}

			time.Sleep(100 * time.Millisecond)
		}
		return
	}

	defer cntxt.Release()

	cm.Run()
}
