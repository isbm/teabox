package main

import (
	"os"
	"path"

	"gitlab.com/isbm/teabox"
	"gitlab.com/isbm/teabox/teaboxlib"
	"gitlab.com/isbm/teabox/teaboxlib/teaboxui"
)

func main() {
	appname := path.Base(os.Args[0])

	app := teabox.GetTeaboxApp().SetGlobalConfig(teaboxlib.NewTeaConf(appname))
	app.SetRoot(teaboxui.InitTeaboxMainWindow().GetContent(), true)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
