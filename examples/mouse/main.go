package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/stephen-fox/winuserio"
)

func main() {
	user32, err := winuserio.LoadUser32DLL()
	if err != nil {
		log.Fatalf("failed to load user32.dll - %s", err.Error())
	}

	fn := func(event winuserio.LowLevelMouseEvent) {
		log.Printf("mouse event: 0x%X", event.WParam)
	}

	listener, err := winuserio.NewLowLevelMouseListener(fn, user32)
	if err != nil {
		log.Fatalf("failed to create mouse listener - %s", err.Error())
	}

	log.Println("now listening for mouse events - press Ctrl+C to stop")

	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	select {
	case err := <-listener.OnDone():
		log.Fatalf("keyboard listener stopped unexpectedly - %v", err)
	case <-interrupts:
	}

	listener.Release()
}
