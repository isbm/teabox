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

	return c
}

func (tsw *TeaSTDOUTWindow) GetWindow() *crtview.TextView {
	return tsw.w
}

func (tsw *TeaSTDOUTWindow) Action(callback, cmdpath, cmdargs string) error {
	if callback != "" {
		// Setup local action for the future instance
		teabox.GetTeaboxApp().GetCallbackServer().AddLocalAction(func(call *teaboxlib.TeaboxAPICall) {
			switch call.GetClass() {
			case "LOGGER-STATUS":
				tsw.statusBar.SetText(call.GetString())
				teabox.GetTeaboxApp().Draw()
			case "LOGGER-TITLE":
				tsw.titleBar.SetText(call.GetString())
				teabox.GetTeaboxApp().Draw()
			}
		})

		// Run the Unix server instance
		if err := teabox.GetTeaboxApp().GetCallbackServer().Start(callback); err != nil {
			teabox.GetTeaboxApp().Stop()
			fmt.Println(err) // That would be a general system problem
		}
	}

	cmd := exec.Command(cmdpath, cmdargs)
	cmd.Stdout = tsw.GetWindow()
	cmd.Stderr = tsw.GetWindow()

	if err := cmd.Run(); err != nil {
		teabox.GetTeaboxApp().Stop()
		fmt.Println("Error:", err)
	}

	// Stop Unix socket
	if teabox.GetTeaboxApp().GetCallbackServer().IsRunning() {
		return teabox.GetTeaboxApp().GetCallbackServer().Stop()
	}

	return nil
}
