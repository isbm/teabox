package teaboxui

import "github.com/isbm/crtview"

type TeaboxWindow interface {
	// GetWindow returns the current window instance
	GetWidget() crtview.Primitive

	// Init UI. This is called from the final class, which will be actually constructing the window UI
	Init() TeaboxWindow
}

type TeaboxBaseWindow struct {
}
