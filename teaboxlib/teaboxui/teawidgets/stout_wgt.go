package teawidgets

import (
	"fmt"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox"
	"gitlab.com/isbm/teabox/teaboxlib"
)

/*
Form Logger is a text view that is handling \r symbols.
It is used to show the output of a called script.
*/

type TeaSTDOUTWindow struct {
	action    func(call *teaboxlib.TeaboxAPICall)
	statusBar *crtview.TextView
	titleBar  *crtview.TextView
	w         *crtview.TextView
	*crtview.Flex
}

func NewTeaSTDOUTWindow() *TeaSTDOUTWindow {
	c := &TeaSTDOUTWindow{
		Flex: crtview.NewFlex(),
	}

	c.SetDirection(crtview.FlexRow)

	// Add top label
	c.titleBar = crtview.NewTextView()
	c.titleBar.SetBackgroundColor(tcell.NewRGBColor(0x88, 0x88, 0x88))
	c.titleBar.SetTextColor(tcell.ColorBlack)
	c.titleBar.SetText("STDOUT output:")

	c.AddItem(c.titleBar, 1, 0, false)

	// Create STDOUT logger
	c.w = crtview.NewTextView()
	c.w.SetBorder(false)
	c.w.SetBorderPadding(1, 1, 2, 2)
	c.w.SetBackgroundColor(tcell.ColorLightGray)
	c.w.SetTextColor(tcell.ColorBlack)
	c.w.SetSkipCursorReturn(true)
	c.w.SetChangedFunc(func() {
		teabox.GetTeaboxApp().Draw()
	})
	c.AddItem(c.w, 0, 1, true)

	// Add bottom status
	c.statusBar = crtview.NewTextView()
	c.statusBar.SetBackgroundColor(tcell.ColorDarkGrey)
	c.statusBar.SetTextColor(tcell.ColorBlack)
	c.statusBar.SetText("some status here")
	c.AddItem(c.statusBar, 1, 0, false)

	// Action definition
	c.action = func(call *teaboxlib.TeaboxAPICall) {
		switch call.GetClass() {
		case "LOGGER-STATUS":
			c.statusBar.SetText(call.GetString())
			teabox.GetTeaboxApp().Draw()
		case "LOGGER-TITLE":
			c.titleBar.SetText(call.GetString())
			teabox.GetTeaboxApp().Draw()
		}
	}

	return c
}

func (tsw *TeaSTDOUTWindow) AsWidgetPrimitive() crtview.Primitive {
	var w TeaboxLandingWindow = tsw
	return w.(crtview.Primitive)
}

func (tsw *TeaSTDOUTWindow) GetWindow() *crtview.TextView {
	return tsw.w
}

func (tsw *TeaSTDOUTWindow) GetWindowAction() func(call *teaboxlib.TeaboxAPICall) {
	return tsw.action
}

func (tsw *TeaSTDOUTWindow) StopListener() error {
	// Stop Unix socket
	if teabox.GetTeaboxApp().GetCallbackServer().IsRunning() {
		return teabox.GetTeaboxApp().GetCallbackServer().Stop()
	}
	return nil
}

func (tsw *TeaSTDOUTWindow) Action(cmdpath string, cmdargs ...string) error {
	cmd := exec.Command(cmdpath, cmdargs...)
	cmd.Stdout = tsw.GetWindow()
	cmd.Stderr = tsw.GetWindow()

	if err := cmd.Run(); err != nil {
		teabox.GetTeaboxApp().Stop()
		fmt.Println("Error:", err)
	}

	return nil
}
