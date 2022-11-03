package teaboxui

import (
	"fmt"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox"
)

/*
Form Logger is a text view that is handling \r symbols.
It is used to show the output of a called script.
*/

type TeaSTDOUTWindow struct {
	sb *crtview.TextView
	w  *crtview.TextView
	*crtview.Flex
}

func NewTeaSTDOUTWindow() *TeaSTDOUTWindow {
	c := &TeaSTDOUTWindow{
		Flex: crtview.NewFlex(),
	}

	c.SetDirection(crtview.FlexRow)

	// Add top label
	label := crtview.NewTextView()
	label.SetBackgroundColor(tcell.ColorDarkGrey)
	label.SetTextColor(tcell.ColorLightGray)
	label.SetText("STDOUT output:")

	c.AddItem(label, 1, 0, false)

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
	c.sb = crtview.NewTextView()
	c.sb.SetBackgroundColor(tcell.ColorDarkGrey)
	c.sb.SetTextColor(tcell.ColorLightGray)
	c.sb.SetText("some status here")
	c.AddItem(c.sb, 1, 0, false)

	return c
}

func (tsw *TeaSTDOUTWindow) GetWindow() *crtview.TextView {
	return tsw.w
}

func (tsw *TeaSTDOUTWindow) Action(callback, cmdpath, cmdargs string) error {
	if callback != "" {
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
