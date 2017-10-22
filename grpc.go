package main

import (
	"fmt"
	"log"
	"time"

	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/valyala/gorpc"
)

func init() {
	gorpc.RegisterType(Request{})
	gorpc.RegisterType(map[string]string{})
}

type GRPCClient struct {
	client *gorpc.Client
}

func (gc *GRPCClient) Add(url, user, secret string) error {
	_, err := gc.Send(Request{"add", url, user, secret})
	return err
}

func (gc *GRPCClient) Delete(url string) error {
	_, err := gc.Send(Request{Command: "delete", ServerURL: url})
	return err
}

func (gc *GRPCClient) Get(url string) (string, string, error) {
	ret, err := gc.Send(Request{Command: "get", ServerURL: url})
	if err != nil {
		return "", "", err
	}

	retMap := ret.(map[string]string)

	if _, ok := retMap["Username"]; !ok {
		return "", "", credentials.NewErrCredentialsNotFound()
	}

	return retMap["Username"], retMap["Secret"], nil
}

func (gc *GRPCClient) List() (map[string]string, error) {
	ret, err := gc.Send(Request{Command: "list"})
	if err != nil {
		return nil, fmt.Errorf("error listing credentials")
	}

	return ret.(map[string]string), nil
}

func (gc *GRPCClient) Send(r Request) (interface{}, error) {
	return gc.client.Call(r)
}

type GRPCTransport struct {
	System
	server *gorpc.Server
}

func (gt *GRPCTransport) Client(d *CacheMem) (Client, error) {
	gorpc.SetErrorLogger(gorpc.NilErrorLogger)
	client := gorpc.NewUnixClient("/tmp/transport.sock")
	client.RequestTimeout = 100 * time.Millisecond
	client.Start()

	return &GRPCClient{client}, nil
}

func (gt *GRPCTransport) Initialize(d *CacheMem) error {
	gt.server = gorpc.NewUnixServer("/tmp/transport.sock", func(clientAddr string, request interface{}) interface{} {

		// keep daemon alive
		d.Alive()

		credStore := d.Creds()

		switch request.(type) {
		case string:
			return request
		case Request:
			req := request.(Request)

			log.Println("GRPC request:", req.Command)

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
					d.Done()
				}()
			}
		}
		return nil
	})

	log.Println("starting GRPC server")
	if err := gt.server.Start(); err != nil {
		log.Fatalf("Cannot start rpc server: %s", err)
	}

	return nil
}

func (gt *GRPCTransport) Shutdown(d *CacheMem) error {
	gt.server.Stop()

	return nil
}
