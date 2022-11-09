package teawidgets

import (
	"os/exec"

	"github.com/gdamore/tcell/v2"
	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	"github.com/isbm/crtview"
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
	load := &TeaboxArgsLoadingWindow{
		Flex: crtview.NewFlex(),
	}
	load.init()

	return load
}

func (ld *TeaboxArgsLoadingWindow) init() {
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
	ld.statusBar.SetText("Loading...")
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
	ld.progressBar.SetProgress((100 / ld.progressSteps) * ld.progressOffset)
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
func (ld *TeaboxArgsLoadingWindow) SetAction(action func()) {
	ld.action = action
}

func (ld *TeaboxArgsLoadingWindow) Load(cmd string, args ...string) error {
	go func() {
		if err := exec.Command(cmd, args...).Run(); err != nil {
			ld.GetLogger().Panic(err)
		}

		ld.action()
	}()
	return nil
}
