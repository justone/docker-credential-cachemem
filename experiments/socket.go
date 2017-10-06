package main

// From: https://github.com/golang/go/issues/11822
//
// You have two choices. One is to call syscall.Umask(0077) before creating the
// socket. This will affect your entire process and will make all created
// files, including the socket, disable the lower bits, so that the socket will
// be mode 0700 instead of 0777. The other is to call os.Chmod after creating
// the socket. This way there would still be a window where the socket has the
// 0777 mode, so if you are worried about attackers and not just appeasing
// nginx then that might not be preferable.
//
// You only need either the syscall.Umask or the os.Chmod, not both.
//
// This program demonstrates both:

import (
	"log"
	"net"
	"os"
	"syscall"
)

func main() {
	syscall.Umask(0077)
	l, err := net.Listen("unix", "/tmp/asdf")
	if err != nil {
		log.Fatal(err)
	}
	check()
	if err := os.Chmod("/tmp/asdf", 0700); err != nil {
		log.Fatal(err)
	}
	check()
	l.Close()
}

func check() {
	fi, err := os.Stat("/tmp/asdf")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("mode", fi.Mode())
}
