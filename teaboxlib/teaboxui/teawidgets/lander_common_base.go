package teawidgets

import "gitlab.com/isbm/teabox"

// Common methods mixin
type teaCommonBaseWindowLander struct {
}

// StopListener
func (base *teaCommonBaseWindowLander) StopListener() error {
	// Stop Unix socket
	if teabox.GetTeaboxApp().GetCallbackServer().IsRunning() {
		return teabox.GetTeaboxApp().GetCallbackServer().Stop()
	}
	return nil
}
