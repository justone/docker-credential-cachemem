package main

import (
	"fmt"
	"time"
)

func NewCacheMem(mode string) *CacheMem {
	system := &DefaultSystem{}

	var transp Transport
	if mode == "grpc" {
		transp = &GRPCTransport{System: system}
	}

	return &CacheMem{
		system,
		transp,
		make(chan bool),
		make(chan bool),
		make(map[string]cred),
	}
}

type Client interface {
	Add(string, string, string) error
	Delete(string) error
	Get(string) (string, string, error)
	List() (map[string]string, error)
	Send(Request) (interface{}, error)
}

type Transport interface {
	Initialize(*CacheMem) error
	Shutdown(*CacheMem) error
	Client(*CacheMem) (Client, error)
}

type CacheMem struct {
	System
	Transport
	done, alive chan bool
	creds       map[string]cred
}

func (gd *CacheMem) Done() {
	gd.done <- true
}

func (gd *CacheMem) Alive() {
	gd.alive <- true
}

func (gd *CacheMem) Creds() map[string]cred {
	return gd.creds
}

func (gd *CacheMem) Client() (Client, error) {
	return gd.Transport.Client(gd)
}

func (gd *CacheMem) Run() {
	gd.Transport.Initialize(gd)

	go func() {
		for {
			select {
			case <-gd.alive:
				fmt.Println("keeping alive longer")
			case <-time.After(60 * time.Second):
				fmt.Println("stopping because of timeout")
				gd.done <- true
				return
			}
		}
	}()

	fmt.Println("waiting for done")
	<-gd.done
	fmt.Println("done received")
	gd.Transport.Shutdown(gd)

	fmt.Println("all done")
}

type cred struct {
	username, secret string
}

type Request struct {
	Command, ServerURL, Username, Secret string
}

func (gd *CacheMem) Stop() {
	cl, _ := gd.Client()

	if _, err := cl.Send(Request{Command: "stop"}); err != nil {
		fmt.Println("error stopping daemon")
	}
}

func (gd *CacheMem) Bump() {
	cl, _ := gd.Client()

	if _, err := cl.Send(Request{Command: "bump"}); err != nil {
		fmt.Println("error bumping daemon")
	}
}
