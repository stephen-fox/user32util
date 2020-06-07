package main

import (
	"fmt"
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

	fn := func(event winuserio.LowLevelKeyboardEvent) {
		if event.KeyboardButtonAction() == winuserio.WMKeyDown {
			fmt.Printf("%q (%d) down\n", event.HookStruct().VirtualKeyCode(), event.HookStruct().VkCode)
		} else if event.KeyboardButtonAction() == winuserio.WMKeyUp {
			fmt.Printf("%q (%d) up\n", event.HookStruct().VirtualKeyCode(), event.HookStruct().VkCode)
		}
	}

	listener, err := winuserio.NewLowLevelKeyboardListener(fn, user32)
	if err != nil {
		log.Fatalf("failed to create listener - %s", err.Error())
	}

	log.Println("now listening for keyboard events - press Ctrl+C to stop")

	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	select {
	case err := <-listener.OnDone():
		log.Fatalf("keyboard listener stopped unexpectedly - %v", err)
	case <-interrupts:
	}

	listener.Release()
}
