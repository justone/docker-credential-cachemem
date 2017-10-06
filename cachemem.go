package main

import (
	"fmt"

	"github.com/docker/docker-credential-helpers/credentials"
)

type CacheMem struct{}

func (c CacheMem) Add(creds *credentials.Credentials) error {
	cl := newClient()

	if _, err := cl.Call(Request{
		"add",
		creds.ServerURL,
		creds.Username,
		creds.Secret,
	}); err != nil {
		fmt.Println(err)
		return fmt.Errorf("error adding credentials")
	}

	return nil
}

func (c CacheMem) Delete(serverURL string) error {
	cl := newClient()

	if _, err := cl.Call(Request{
		Command:   "delete",
		ServerURL: serverURL,
	}); err != nil {
		return fmt.Errorf("error getting credentials")
	}
	return nil
}

func (c CacheMem) Get(serverURL string) (string, string, error) {
	cl := newClient()

	ret, err := cl.Call(Request{
		Command:   "get",
		ServerURL: serverURL,
	})
	if err != nil {
		fmt.Println(err)
		return "", "", fmt.Errorf("error getting credentials")
	}

	fmt.Println(ret)
	retMap := ret.(map[string]string)

	fmt.Println(retMap)

	if _, ok := retMap["Username"]; !ok {
		return "", "", credentials.NewErrCredentialsNotFound()
	}

	return retMap["Username"], retMap["Secret"], nil
}

func (c CacheMem) List() (map[string]string, error) {
	cl := newClient()

	ret, err := cl.Call(Request{
		Command: "list",
	})
	if err != nil {
		return nil, fmt.Errorf("error listing credentials")
	}

	return ret.(map[string]string), nil
}
