package main

import (
	"fmt"
	"log"

	"github.com/lazywhite/sample-device-plugin/pkg/foo"
)



func main() {
	if err := foo.Register(); err != nil {
		log.Fatal(fmt.Errorf("failed to register with kubelet: %s", err))
	}

	stop := make(chan struct{})
	go func() {
		if err := foo.WatchKubeletRestart(stop); err != nil {
			log.Fatal(fmt.Errorf("error watching kubelet restart: %s", err))
		}
	}()

	if err := foo.Serve(); err != nil {
		log.Fatal(fmt.Errorf("error running server: %s", err))
	}
	stop <- struct{}{}
}

