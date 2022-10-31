package main

import (
	"os"
	"path"

	"gitlab.com/isbm/teabox/teaboxlib"
	"gitlab.com/isbm/teabox/teaboxlib/teaboxui"
)

func main() {
	appname := path.Base(os.Args[0])

	app := teaboxlib.GetTeaboxApp()
	app.SetRoot(teaboxui.NewTeaboxMainWindow(app, teaboxlib.NewTeaConf(appname)).GetContent(), true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
