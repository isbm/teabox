package teaboxui

import "github.com/isbm/crtview"

type TeaboxWindow interface {
	// GetWindow returns the current window instance
	GetWidget() crtview.Primitive

	// SetApp sets the current application reference. This allows app control from the current window.
	SetApp(app *crtview.Application)

	// Init UI. This is called from the final class, which will be actually constructing the window UI
	Init() TeaboxWindow
}

type TeaboxBaseWindow struct {
	appref *crtview.Application
}

func (tbw *TeaboxBaseWindow) SetApp(app *crtview.Application) {
	tbw.appref = app
}
