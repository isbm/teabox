package main

import (
	"os"
	"path"

	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox/teaboxlib"
	"gitlab.com/isbm/teabox/teaboxlib/teaboxui"
)

func main() {
	app := crtview.NewApplication()
	appname := path.Base(os.Args[0])
	app.SetRoot(teaboxui.NewTeaboxMainWindow(app, teaboxlib.NewTeaConf(appname)).GetContent(), true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
