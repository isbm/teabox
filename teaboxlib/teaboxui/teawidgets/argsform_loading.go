package teawidgets

import (
	"fmt"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox"
	"gitlab.com/isbm/teabox/teaboxlib"
)

type TeaboxArgsLoadingWindow struct {
	statusBar      *crtview.TextView
	progressBar    *crtview.ProgressBar
	progressSteps  int
	progressOffset int

	action func()

	*crtview.Flex
	wzlib_logger.WzLogger
}

func NewTeaboxArgsLoadingWindow() *TeaboxArgsLoadingWindow {
	return (&TeaboxArgsLoadingWindow{
		Flex: crtview.NewFlex(),
	}).init()

}

func (ld *TeaboxArgsLoadingWindow) init() *TeaboxArgsLoadingWindow {
	ld.SetDirection(crtview.FlexColumn)
	ld.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)

	content := crtview.NewFlex()
	content.SetDirection(crtview.FlexRow)
	content.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	content.SetBorder(false)

	// Vertical spacer
	vspacer := crtview.NewBox()
	vspacer.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	vspacer.SetBorder(false)

	// Add a vertical spacer from the top
	content.AddItem(vspacer, 0, 1, false)

	ld.statusBar = crtview.NewTextView()
	ld.statusBar.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	ld.statusBar.SetBorder(false)
	ld.statusBar.SetText("")
	content.AddItem(ld.statusBar, 1, 1, false)

	ld.progressBar = crtview.NewProgressBar()
	ld.progressBar.SetEmptyColor(tcell.ColorDarkGreen)
	ld.progressBar.SetFilledColor(tcell.ColorYellow)
	ld.progressBar.SetProgress(42)
	content.AddItem(ld.progressBar, 1, 1, false)

	// Add a vertical spacer from the bottom
	content.AddItem(vspacer, 0, 1, false)

	// Add a painted box as a padding, otherwise it will be just black space
	padding := crtview.NewBox()
	padding.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	padding.SetBorder(false)

	// Assemble the final padded content
	ld.AddItem(padding, 5, 1, false)
	ld.AddItem(content, 0, 1, false)
	ld.AddItem(padding, 5, 1, false)

	ld.Reset()

	return ld
}

// AllocateProgress allocates how many times increment by one should
// reach its 100%
func (ld *TeaboxArgsLoadingWindow) AllocateProgress(steps int) {
	ld.progressSteps = steps
}

// IncrementProgress increments it by one. One steps represents an amount
// of percentage, previously allocated.
func (ld *TeaboxArgsLoadingWindow) IncrementProgress() {
	ld.progressOffset++

	var p int = 100
	if ld.progressOffset < ld.progressSteps {
		p = (p / ld.progressSteps) * ld.progressOffset
	}

	ld.progressBar.SetProgress(p)
}

// SetProgress in percentage
func (ld *TeaboxArgsLoadingWindow) SetProgress(p int) {
	ld.progressBar.SetProgress(p)
}

// SetStatus of the status bar
func (ld *TeaboxArgsLoadingWindow) SetStatus(status string) {
	ld.statusBar.SetText(status)
}

// SetAction
func (ld *TeaboxArgsLoadingWindow) SetAfterLoadAction(action func()) {
	ld.action = action
}

func (ld *TeaboxArgsLoadingWindow) Reset() {
	ld.SetStatus("Loading...")
	ld.progressOffset = 0
	ld.progressSteps = 0
	ld.SetProgress(0)
}

// GetSocketAcceptAction is a function for Unix socket on action
func (ld *TeaboxArgsLoadingWindow) GetSocketAcceptAction() func(*teaboxlib.TeaboxAPICall) {
	return func(call *teaboxlib.TeaboxAPICall) {
		if !ld.IsVisible() {
			return // Don't update here anything, we are invisible
		}
		switch call.GetClass() {
		case teaboxlib.INIT_SET_PROGRESS:
			ld.SetProgress(call.GetInt())
		case teaboxlib.INIT_INC_PROGRESS:
			ld.IncrementProgress()
		case teaboxlib.INIT_ALLOC_PROGRESS:
			ld.AllocateProgress(call.GetInt())
		case teaboxlib.INIT_SET_STATUS:
			ld.SetStatus(call.GetString())
		case teaboxlib.INIT_RESET:
			ld.Reset()
		}

		teabox.GetTeaboxApp().Draw()
	}
}

func (ld *TeaboxArgsLoadingWindow) Load(cmd string, args ...string) error {
	go func() {
		if output, err := exec.Command(cmd, args...).CombinedOutput(); err != nil {
			teabox.GetTeaboxApp().Stop(fmt.Sprintf("Error: Failure while loading form from setup command.\nCommand: %s\nargs: %v\nSystem exit: %s\nError details: %s", cmd, args, err.Error(), output))
		}

		ld.action()
	}()
	return nil
}
