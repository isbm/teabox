package teawidgets

import (
	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox"
	"gitlab.com/isbm/teabox/teaboxlib"
)

type TeaProgressWindowLander struct {
	steps        int
	stepsOffset  int
	lookupPrefix string
	lookupRegex  string

	action      func(call *teaboxlib.TeaboxAPICall)
	eventBar    *crtview.TextView    // Like a status bar, but shows a chunk of the progress
	generalInfo *crtview.TextView    // General text information (static per module)
	progressBar *crtview.ProgressBar // Progressbar itself

	*crtview.Flex
}

func NewTeaProgressWindowLander() *TeaProgressWindowLander {
	return (&TeaProgressWindowLander{
		eventBar:    crtview.NewTextView(),
		generalInfo: crtview.NewTextView(),
		progressBar: crtview.NewProgressBar(),
	}).init()
}

func (pl *TeaProgressWindowLander) init() *TeaProgressWindowLander {
	// Construct UI here

	// Define API action
	pl.action = func(call *teaboxlib.TeaboxAPICall) {
		switch call.GetClass() {
		case teaboxlib.PROGRESS_EVENT:
			pl.eventBar.SetText(call.GetString())
		case teaboxlib.PROGRESS_NEXT:
			pl.stepsOffset++
			if pl.stepsOffset >= pl.steps {
				pl.stepsOffset = pl.steps
			}
			pl.progressBar.SetProgress((100 / pl.steps) * pl.stepsOffset)
		case teaboxlib.PROGRESS_ALLOCATE:
			pl.steps = call.GetInt()
		case teaboxlib.PROGRESS_SET:
			pl.progressBar.SetProgress(call.GetInt())
		case teaboxlib.PROGRESS_LOOKUP_PREFIX:
			pl.lookupPrefix = call.GetString()
		case teaboxlib.PROGRESS_LOOKUP_REGEX:
			pl.lookupRegex = call.GetString()
		}

		teabox.GetTeaboxApp().Draw()
	}

	return pl
}
