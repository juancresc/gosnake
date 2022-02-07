package main

import (
	"github.com/nsf/termbox-go"
	"os"
	"os/exec"
	"runtime"
)

func reset() {
	termbox.Sync() // cosmestic purpose
}

func keyPress(keypress chan int) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
keyPressListenerLoop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break keyPressListenerLoop
			default:
				// we only want to read a single character or one key pressed event
				reset()
				keypress <- int(ev.Ch)

			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
	keypress <- 27
}

var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
