package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/stephen-fox/user32util"
)

func main() {
	user32, err := user32util.LoadUser32DLL()
	if err != nil {
		log.Fatalf("failed to load user32.dll - %s", err.Error())
	}

	fn := func(event user32util.LowLevelMouseEvent) {
		log.Printf("mouse event: %+v", event.Struct.Point)
	}

	listener, err := user32util.NewLowLevelMouseListener(fn, user32)
	if err != nil {
		log.Fatalf("failed to create mouse listener - %s", err.Error())
	}

	log.Println("now listening for mouse events - press Ctrl+C to stop")

	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	select {
	case err := <-listener.OnDone():
		log.Fatalf("mouse listener stopped unexpectedly - %v", err)
	case <-interrupts:
	}

	listener.Release()
}
