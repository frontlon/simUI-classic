package main

import (
	"flag"
	"log"

	"simUI/code/utils/go-sciter"
	"simUI/code/utils/go-sciter/window"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("html file needed")
	}
	rect := sciter.NewRect(300, 300, 300, 400)
	// create window
	w, err := window.New(sciter.DefaultWindowCreateFlag, rect)
	if err != nil {
		log.Fatal(err)
	}

	w.LoadFile(flag.Arg(0))
	w.Show()
	w.Run()
}
