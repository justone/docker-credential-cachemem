package main

import "github.com/docker/docker-credential-helpers/credentials"

type CacheMem struct{}

func (c CacheMem) Add(creds *credentials.Credentials) error {
	return nil
}

func (c CacheMem) Delete(serverURL string) error {
	return nil
}

func (c CacheMem) Get(serverURL string) (string, string, error) {
	return "", "", nil
}

func (c CacheMem) List() (map[string]string, error) {
	return map[string]string{}, nil
}
