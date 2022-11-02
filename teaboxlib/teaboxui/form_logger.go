package teaboxui

import (
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
