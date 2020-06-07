package main

import (
	"flag"
	"log"
	"time"

	"github.com/stephen-fox/winuserio"
)

func main() {
	sendAfter := flag.Duration("s", 1*time.Second, "The amount of time in seconds to wait before sending value")
	mouse := flag.Bool("mouse", false, "Send mouse input instead of keyboard input")

	flag.Parse()

	user32, err := winuserio.LoadUser32DLL()
	if err != nil {
		log.Fatalf("failed to load user32.dll - %s", err.Error())
	}

	if *mouse {
		for {
			log.Printf("will send 0x%X in %s", 0x41, sendAfter.String())

			time.Sleep(*sendAfter)

			err := winuserio.SendMouseInput(winuserio.MouseInput{
				DwFlags: winuserio.MouseEventFRightDown,
			}, user32)
			if err != nil {
				log.Fatalf("failed to send input - %s", err.Error())
			}
		}
	} else {
		for {
			log.Printf("will send 0x%X in %s", 0x41, sendAfter.String())

			time.Sleep(*sendAfter)

			err := winuserio.SendKeydbInput(winuserio.KeybdInput{
				WVK: 0x41,
			}, user32)
			if err != nil {
				log.Fatalf("failed to send input - %s", err.Error())
			}
		}
	}
}
