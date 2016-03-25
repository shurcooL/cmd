// dump_glfw3_joysticks dumps state of attached joysticks.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/shurcooL/go-goon"
)

func main() {
	err := glfw.Init()
	if err != nil {
		log.Fatalln(err)
	}
	defer glfw.Terminate()

	for {
		var present bool
		for joy := glfw.Joystick1; joy <= glfw.JoystickLast; joy++ {
			if glfw.JoystickPresent(joy) {
				present = true
				goon.DumpExpr(joy)
				goon.DumpExpr(glfw.GetJoystickName(joy))
				goon.DumpExpr(glfw.GetJoystickButtons(joy))
				goon.DumpExpr(glfw.GetJoystickAxes(joy))
			}
		}
		if !present {
			fmt.Println("no joysticks present")
		}

		time.Sleep(time.Second)
	}
}
