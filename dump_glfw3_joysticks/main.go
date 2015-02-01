package main

import (
	"fmt"
	"time"

	glfw "github.com/shurcooL/glfw3"
	"github.com/shurcooL/go-goon"
)

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
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
