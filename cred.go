package main

import (
	"fmt"

	"github.com/docker/docker-credential-helpers/credentials"
)

type CredHandler struct {
	*CacheMem
}

func (c CredHandler) Add(creds *credentials.Credentials) error {
	cl, _ := c.Client()

	if err := cl.Add(creds.ServerURL, creds.Username, creds.Secret); err != nil {
		return fmt.Errorf("error adding credentials")
	}

	return nil
}

func (c CredHandler) Delete(serverURL string) error {
	cl, _ := c.Client()

	if err := cl.Delete(serverURL); err != nil {
		return fmt.Errorf("error getting credentials")
	}
	return nil
}

func (c CredHandler) Get(serverURL string) (string, string, error) {
	cl, _ := c.Client()

	user, secret, err := cl.Get(serverURL)
	if err != nil {
		return "", "", err
	}

	return user, secret, nil
}

func (c CredHandler) List() (map[string]string, error) {
	cl, _ := c.Client()

	ret, err := cl.List()
	if err != nil {
		return nil, fmt.Errorf("error listing credentials")
	}

	return ret, nil
}
