package main

import (
	"fmt"
	"os"
	"path"

	"gitlab.com/isbm/teabox"
	"gitlab.com/isbm/teabox/teaboxlib"
	"gitlab.com/isbm/teabox/teaboxlib/teaboxui"
)

var VERSION = "0.1"

func main() {
	if os.Getenv("TERM") != "xterm-256color" {
		fmt.Println("Terminal should work in 256 color mode.")
		os.Exit(1)
	}

	appname := path.Base(os.Args[0])

	conf, err := teaboxlib.NewTeaConf(appname)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	defer os.Remove(conf.GetSocketPath())

	app := teabox.GetTeaboxApp().SetGlobalConfig(conf)
	app.SetRoot(teaboxui.InitTeaboxMainWindow().GetContent(), true)

	if err := app.Run(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Println(teabox.GetTeaboxQuitMessage())
	}
}
