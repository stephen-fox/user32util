package main

import (
	"fmt"
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

	fn := func(event user32util.LowLevelKeyboardEvent) {
		if event.KeyboardButtonAction() == user32util.WMKeyDown {
			fmt.Printf("%q (%d) down\n", event.HookStruct().VirtualKeyCode(), event.HookStruct().VkCode)
		} else if event.KeyboardButtonAction() == user32util.WMKeyUp {
			fmt.Printf("%q (%d) up\n", event.HookStruct().VirtualKeyCode(), event.HookStruct().VkCode)
		}
	}

	listener, err := user32util.NewLowLevelKeyboardListener(fn, user32)
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
