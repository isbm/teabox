package teawidgets

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox"
	"gitlab.com/isbm/teabox/teaboxlib"
)

/*
Form Logger is a text view that is handling \r symbols.
It is used to show the output of a called script.
*/

type TeaLoggerWindowLander struct {
	action    func(call *teaboxlib.TeaboxAPICall)
	statusBar *crtview.TextView
	titleBar  *crtview.TextView
	w         *crtview.TextView

	*teaCommonBaseWindowLander
	*crtview.Flex
}

func NewTeaLoggerWindowLander() *TeaLoggerWindowLander {
	c := &TeaLoggerWindowLander{
		Flex: crtview.NewFlex(),
	}

	c.SetDirection(crtview.FlexRow)

	// Add top label
	c.titleBar = crtview.NewTextView()
	c.titleBar.SetBackgroundColor(tcell.NewRGBColor(0x88, 0x88, 0x88))
	c.titleBar.SetTextColor(tcell.ColorBlack)

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
	c.AddItem(c.statusBar, 1, 0, false)

	// Action definition
	c.action = func(call *teaboxlib.TeaboxAPICall) {
		switch call.GetClass() {
		case teaboxlib.LOGGER_STATUS:
			c.statusBar.SetText(call.GetString())
		case teaboxlib.LOGGER_TITLE:
			c.titleBar.SetText(call.GetString())
		}
		teabox.GetTeaboxApp().Draw()
	}

	c.Reset()

	return c
}

// Reset all content to the initial values
func (tsw *TeaLoggerWindowLander) Reset() {
	tsw.statusBar.SetText("")
	tsw.titleBar.SetText("")
	tsw.w.SetText("")
}

func (tsw *TeaLoggerWindowLander) AsWidgetPrimitive() crtview.Primitive {
	var w TeaboxLandingWindow = tsw
	return w.(crtview.Primitive)
}

func (tsw *TeaLoggerWindowLander) GetWindow() *crtview.TextView {
	return tsw.w
}

func (tsw *TeaLoggerWindowLander) GetWindowAction() func(call *teaboxlib.TeaboxAPICall) {
	return tsw.action
}

func (tsw *TeaLoggerWindowLander) Action(cmdpath string, cmdargs ...string) error {
	cmd := exec.Command(cmdpath, cmdargs...)
	cmd.Stdout = tsw.GetWindow()
	cmd.Stderr = tsw.GetWindow()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(fmt.Sprintf("Error: command \"%s %s\" quit as %s", cmdpath, strings.Join(cmdargs, " "), err.Error()))
	}

	return nil
}
