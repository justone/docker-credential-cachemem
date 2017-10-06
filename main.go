package main

import (
	"os"

	"github.com/docker/docker-credential-helpers/credentials"
)

func main() {
	mode := "grpc"

	cm := NewCacheMem(mode)

	if len(os.Args) > 1 && os.Args[1] == "daemon" {
		cm.Run()
	} else if len(os.Args) > 1 && os.Args[1] == "stop" {
		cm.Stop()
	} else if len(os.Args) > 1 && os.Args[1] == "bump" {
		cm.Bump()
	} else {
		credentials.Serve(CredHandler{cm})
	}
}
