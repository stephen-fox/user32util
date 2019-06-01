package main

import (
	"fmt"
	"log"
	"time"

	"github.com/stephen-fox/winuserio"
)

func main() {
	fn := func(event winuserio.LowLevelKeyboardEvent) {
		if event.KeyboardButtonAction() == winuserio.WMKeyDown {
			fmt.Printf("%q (%d) down\n", event.HookStruct().VirtualKeyCode(), event.HookStruct().VkCode)
		} else if event.KeyboardButtonAction() == winuserio.WMKeyUp {
			fmt.Printf("%q (%d) up\n", event.HookStruct().VirtualKeyCode(), event.HookStruct().VkCode)
		}
	}

	listener, err := winuserio.NewLowLevelKeyboardListener(fn)
	if err != nil {
		log.Fatalf("failed to create listener - %s", err.Error())
	}
	defer listener.Release()

	for {
		time.Sleep(1 * time.Second)
	}
}
