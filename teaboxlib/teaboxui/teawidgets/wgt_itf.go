package teawidgets

import (
	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox/teaboxlib"
)

const (
	LANDING_W_LOGGER   = "_logger"
	LANDING_W_PROGRESS = "_progress"
	LANDING_W_LIST     = "_list"

	INTRO_W = "_intro"
)

// TeaboxLandingWindow interface.
// Window widgets are the entire panels, that can use the whole window space.
type TeaboxLandingWindow interface {
	// Cast to the primitive, so it can be added to panels etc
	AsWidgetPrimitive() crtview.Primitive

	// Action implementation, which takes:
	//   - "callback" path of the Unix socket to start a listener on it
	//   - "cmdpath" is a command absolute path (preferrably) or env returned
	//   - "cmdargs" are the arguments to the "cmdpath"
	Action(cmdpath string, cmdargs ...string) error
	StopListener() error

	// Return window action on Unix socket calls, specific per this widget
	GetWindowAction() func(call *teaboxlib.TeaboxAPICall)
}
