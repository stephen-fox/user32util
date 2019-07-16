package main

import (
	"github.com/stephen-fox/winuserio"
	"log"
)

func main() {
	devices, err := winuserio.RawInputDevices()
	if err != nil {
		log.Fatalln(err.Error())
	}

	for i := range devices {
		log.Println(devices[i])
	}
}
