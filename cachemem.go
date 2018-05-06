package main

import (
	"fmt"
	"log"
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
				log.Println("activity detected, resetting timeout")
			case <-time.After(4 * 60 * time.Minute):
				log.Println("timeout reached, signaling 'done'")
				gd.done <- true
				return
			}
		}
	}()

	log.Println("daemon started, waiting for requests")
	<-gd.done
	log.Println("done signal received, shutting down...")
	gd.Transport.Shutdown(gd)

	log.Println("shutdown complete")
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
		fmt.Println("error sending stop signal")
	}
}

func (gd *CacheMem) Bump() {
	cl, _ := gd.Client()

	if _, err := cl.Send(Request{Command: "bump"}); err != nil {
		fmt.Println("error sending bump signal")
	}
}
