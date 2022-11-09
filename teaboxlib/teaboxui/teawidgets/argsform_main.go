package teawidgets

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox/teaboxlib"
)

type TeaboxArgsMainWindow struct {
	cmdId, title, subtitle string
	*crtview.Form
}

func NewTeaboxArgsMainWindow(title, subtitle string) *TeaboxArgsMainWindow {
	return (&TeaboxArgsMainWindow{
		Form:     crtview.NewForm(),
		title:    title,
		subtitle: subtitle,
	}).init()
}

func (tmw *TeaboxArgsMainWindow) init() *TeaboxArgsMainWindow {
	// XXX: ID (Set/Get by ID) needs to be re-thought
	tmw.SetTitle(fmt.Sprintf("%s - %s", tmw.title, tmw.subtitle))
	tmw.cmdId = fmt.Sprintf("%s - %s", tmw.title, tmw.subtitle)

	tmw.SetBorder(true)

	// Colors
	tmw.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	tmw.SetFieldTextColor(tcell.ColorWhite)
	tmw.SetFieldBackgroundColor(teaboxlib.WORKSPACE_HEADER)
	tmw.SetFieldBackgroundColorFocused(tcell.ColorGreenYellow)
	tmw.SetFieldTextColorFocused(teaboxlib.WORKSPACE_BACKGROUND)

	tmw.SetButtonBackgroundColor(teaboxlib.FORM_BUTTON_BACKGROUND)
	tmw.SetButtonBackgroundColorFocused(teaboxlib.FORM_BUTTON_BACKGROUND_SELECTED)
	tmw.SetButtonTextColor(teaboxlib.FORM_BUTTON_TEXT)
	tmw.SetButtonTextColorFocused(teaboxlib.FORM_BUTTON_TEXT_SELECTED)

	// Buttons align
	tmw.SetButtonsAlign(crtview.AlignRight)
	tmw.SetButtonsToBottom()

	return tmw
}

func (tmw *TeaboxArgsMainWindow) GetId() string {
	return tmw.cmdId
}
