package main

import (
	"flag"
	"log"
	"time"

	"github.com/stephen-fox/user32util"
)

func main() {
	sendAfter := flag.Duration("s", 1*time.Second, "The amount of time in seconds to wait before sending value")
	mouse := flag.Bool("mouse", false, "Send mouse input instead of keyboard input")

	flag.Parse()

	user32, err := user32util.LoadUser32DLL()
	if err != nil {
		log.Fatalf("failed to load user32.dll - %s", err.Error())
	}

	if *mouse {
		for {
			log.Printf("will send 0x%X in %s", 0x41, sendAfter.String())

			time.Sleep(*sendAfter)

			err := user32util.SendMouseInput(user32util.MouseInput{
				DwFlags: user32util.MouseEventFRightDown,
			}, user32)
			if err != nil {
				log.Fatalf("failed to send input - %s", err.Error())
			}
		}
	} else {
		for {
			log.Printf("will send 0x%X in %s", 0x41, sendAfter.String())

			time.Sleep(*sendAfter)

			err := user32util.SendKeydbInput(user32util.KeybdInput{
				WVK: 0x41,
			}, user32)
			if err != nil {
				log.Fatalf("failed to send input - %s", err.Error())
			}
		}
	}
}
