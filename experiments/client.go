package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	xdg "github.com/casimir/xdg-go"
	"github.com/tv42/httpunix"
)

func main() {
	fmt.Println(os.Args)

	app := xdg.App{Name: "cachemem"}
	fmt.Println(app.ConfigPath("foo.yaml"))
	fmt.Println(app.ConfigPath(""))
	fmt.Println(app.DataPath("foo.yaml"))
	fmt.Println(app.DataPath(""))

	return
	// This example shows using a customized http.Client.
	u := &httpunix.Transport{
		DialTimeout:           100 * time.Millisecond,
		RequestTimeout:        1 * time.Second,
		ResponseHeaderTimeout: 1 * time.Second,
	}
	u.RegisterLocation("myservice", "/home/nate/listen.sock")

	var client = http.Client{
		Transport: u,
	}

	resp, err := client.Get("http+unix://myservice/foo")
	if err != nil {
		log.Fatal(err)
	}
	buf, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", buf)
	resp.Body.Close()
}
