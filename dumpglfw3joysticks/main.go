// dumpglfw3joysticks dumps state of attached joysticks.
package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/shurcooL/go-goon"
)

var gamepadFlag = flag.Bool("gamepad", false, "Dump state of attached gamepads.")

func init() { runtime.LockOSThread() }
func main() {
	flag.Parse()

	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	glfw.InitHint(glfw.JoystickHatButtons, glfw.False)
	err := glfw.Init()
	if err != nil {
		return err
	}
	defer glfw.Terminate()

	if *gamepadFlag {
		for {
			var present bool
			for j := glfw.Joystick1; j <= glfw.JoystickLast; j++ {
				if !j.IsGamepad() {
					continue
				}
				present = true
				goon.DumpExpr(j)
				goon.DumpExpr(j.GetGamepadName())
				goon.DumpExpr(j.GetGamepadState())
			}
			if !present {
				fmt.Println("no gamepads present")
			}

			time.Sleep(time.Second)
		}
	}

	for {
		var present bool
		for j := glfw.Joystick1; j <= glfw.JoystickLast; j++ {
			if !j.Present() {
				continue
			}
			present = true
			goon.DumpExpr(j)
			goon.DumpExpr(j.GetName())
			goon.DumpExpr(j.GetGUID())
			fmt.Println(len(j.GetButtons()), "buttons:")
			goon.DumpExpr(j.GetButtons())
			goon.DumpExpr(j.GetHats())
			goon.DumpExpr(j.GetAxes())
		}
		if !present {
			fmt.Println("no joysticks present")
		}

		time.Sleep(time.Second)
	}
}
