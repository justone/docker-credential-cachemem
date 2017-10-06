package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/valyala/gorpc"
)

func init() {
	gorpc.RegisterType(Request{})
	gorpc.RegisterType(map[string]string{})
}

func main() {
	mode := "grpc"

	var daemon Daemon
	if mode == "grpc" {
		daemon = &GRPCDaemon{}
	}

	if len(os.Args) > 1 && os.Args[1] == "daemon" {
		daemon.Run()
	} else if len(os.Args) > 1 && os.Args[1] == "stop" {
		daemon.Stop()
	} else if len(os.Args) > 1 && os.Args[1] == "bump" {
		daemon.Bump()
	} else {
		credentials.Serve(CacheMem{})
	}
}

type Daemon interface {
	Run()
	Stop()
	Bump()
}

type GRPCDaemon struct{}

func (gd *GRPCDaemon) Run() {
	fn, done, alive := initDispatcher()
	server := gorpc.NewUnixServer("/tmp/transport.sock", fn)

	fmt.Println("starting")
	if err := server.Start(); err != nil {
		log.Fatalf("Cannot start rpc server: %s", err)
	}

	go func() {
		for {
			select {
			case <-alive:
				fmt.Println("keeping alive longer")
			case <-time.After(60 * time.Second):
				fmt.Println("stopping because of timeout")
				done <- true
				return
			}
		}
	}()

	fmt.Println("waiting for done")
	<-done
	fmt.Println("done received")
	server.Stop()

	fmt.Println("all done")
}

type cred struct {
	username, secret string
}

type Request struct {
	Command, ServerURL, Username, Secret string
}

func initDispatcher() (func(string, interface{}) interface{}, chan bool, chan bool) {
	done := make(chan bool)
	alive := make(chan bool)

	credStore := make(map[string]cred)

	return func(clientAddr string, request interface{}) interface{} {
		fmt.Println("got request", request)

		// keep daemon alive
		alive <- true

		switch request.(type) {
		case string:
			return request
		case Request:
			req := request.(Request)
			fmt.Println(req.Command)
			switch req.Command {
			case "add":
				credStore[req.ServerURL] = cred{req.Username, req.Secret}
			case "delete":
				delete(credStore, req.ServerURL)
			case "get":
				ret := map[string]string{
					"ServerURL": req.ServerURL,
				}

				if cred, ok := credStore[req.ServerURL]; ok {
					ret["Username"] = cred.username
					ret["Secret"] = cred.secret
				}
				fmt.Println("returning", ret)
				return ret
			case "list":
				creds := make(map[string]string)

				for server, cred := range credStore {
					creds[server] = cred.username
				}

				return creds
			case "alive":
				// nothing
			case "stop":
				// nothing
				go func() {
					time.Sleep(time.Second)
					done <- true
				}()
			}
		}
		return nil
	}, done, alive
}

func newClient() *gorpc.Client {

	gorpc.SetErrorLogger(gorpc.NilErrorLogger)
	client := gorpc.NewUnixClient("/tmp/transport.sock")
	client.RequestTimeout = 100 * time.Millisecond
	client.Start()

	return client
}

func (gd *GRPCDaemon) Stop() {
	cl := newClient()

	if _, err := cl.Call("stop"); err != nil {
		fmt.Println("error stopping daemon")
	}
}

func (gd *GRPCDaemon) Bump() {
	cl := newClient()

	if _, err := cl.Call("alive"); err != nil {
		fmt.Println("error bumping daemon")
	}
}
