package teawidgets

import (
	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox/teaboxlib"
)

type TeaboxArgsIntroWindow struct {
	*crtview.Flex
}

func NewTeaboxArgsIntroWindow() *TeaboxArgsIntroWindow {
	return (&TeaboxArgsIntroWindow{Flex: crtview.NewFlex()}).init()
}

func (intro *TeaboxArgsIntroWindow) init() *TeaboxArgsIntroWindow {
	intro.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	intro.SetBorder(true)
	intro.SetBorderColor(teaboxlib.WORKSPACE_HEADER)
	intro.SetDirection(crtview.FlexRow)

	padding := crtview.NewBox()
	padding.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	padding.SetBorder(false)

	// Left padding
	intro.AddItem(padding, 0, 1, false)

	// Content
	cnt := crtview.NewTextView()
	cnt.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	cnt.SetBorder(false)
	cnt.SetText("◀ Select option from the menu")
	cnt.SetTextAlign(crtview.AlignCenter)
	intro.AddItem(cnt, 1, 1, false)

	// Right padding
	intro.AddItem(padding, 0, 1, false)

	return intro
}
