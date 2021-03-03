package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/stephen-fox/user32util"
)

func main() {
	flag.Parse()

	dll, err := user32util.LoadUser32DLL()
	if err != nil {
		log.Fatalf("failed to load user32 dll - %s", err)
	}

	if flag.Arg(0) == "print" {
		log.Println("will show mouse movement")
		getPoints(dll)
	} else {
		if flag.NArg() == 0 {
			log.Fatalf("please specify at least one coordinate pair")
		}

		points := make([]user32util.Point, flag.NArg())
		for i, coordStr := range flag.Args() {
			coordParts := strings.Split(coordStr, ",")
			if len(coordParts) != 2 {
				log.Fatalf("argument number %d is not in the format x,y", i)
			}

			x, err := strconv.Atoi(coordParts[0])
			if err != nil {
				log.Fatalf("failed to parse x coord of argument %d (%s) - %s", i, coordParts[0], err)
			}

			y, err := strconv.Atoi(coordParts[1])
			if err != nil {
				log.Fatalf("failed to parse y coord of argument %d (%s) - %s", i, coordParts[1], err)
			}

			points[i] = user32util.Point{
				X: int32(x),
				Y: int32(y),
			}
		}

		log.Printf("clicking between %+v", points)
		clickBetween(points, dll)
	}
}

func getPoints(dll *user32util.User32DLL) {
	listener, err := user32util.NewLowLevelMouseListener(func(event user32util.LowLevelMouseEvent) {
		log.Printf("point: %+v", event.Struct.Point)
	}, dll)
	if err != nil {
		log.Fatalf("failed to start listner - %s", err)
	}

	interrupts := make(chan os.Signal)
	signal.Notify(interrupts, os.Interrupt)
	select {
	case <-interrupts:
		listener.Release()
	case err := <-listener.OnDone():
		log.Fatalln("listener exited - err is:", err)
	}
}

func clickBetween(sequentialClicks []user32util.Point, dll *user32util.User32DLL) {
	for {
		for _, point := range sequentialClicks {
			log.Printf("moving to point %+v", point)
			_, err := user32util.SetCursorPos(point.X, point.Y, dll)
			if err != nil {
				log.Fatalf("failed to send mouse down input - %s", err)
			}

			log.Println("clicking")
			err = user32util.SendMouseInput(user32util.MouseInput{
				DwFlags: user32util.MouseEventFLeftDown,
			}, dll)
			if err != nil {
				log.Fatalf("failed to send mouse down input - %s", err)
			}

			log.Println("unclick")
			err = user32util.SendMouseInput(user32util.MouseInput{
				DwFlags: user32util.MouseEventFLeftUp,
			}, dll)
			if err != nil {
				log.Fatalf("failed to send mouse up input - %s", err)
			}

			time.Sleep(5*time.Second)
		}
	}
}
