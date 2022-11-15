package teabox

import (
	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox/teaboxlib"
)

// Global reference of the Application. It is accessed via GetTeaboxApp() function.
var __APP_REF *TeaboxApplication
var __MSG_REF string

// TeaboxApplication singleton
type TeaboxApplication struct {
	callback *teaboxlib.TeaboxSocketServer
	config   *teaboxlib.TeaConf

	*crtview.Application
}

// GetCallbackServer, responsible to the Unix socket and modules callbacks
func (ta *TeaboxApplication) GetCallbackServer() *teaboxlib.TeaboxSocketServer {
	return ta.callback
}

// SetGlobalConfig file to the application
func (ta *TeaboxApplication) SetGlobalConfig(conf *teaboxlib.TeaConf) *TeaboxApplication {
	ta.config = conf
	return ta
}

// GetGlobalConfig of the application
func (ta *TeaboxApplication) GetGlobalConfig() *teaboxlib.TeaConf {
	return ta.config
}

// Stop application
func (ta *TeaboxApplication) Stop(message string) {
	__MSG_REF = message
	ta.Application.Stop()
}

// GetTeaboxApp is a factory that returns a singleton of an Application reference
// if it wasn't creted before, it gets created.
func GetTeaboxApp() *TeaboxApplication {
	if __APP_REF == nil {
		__APP_REF = &TeaboxApplication{
			Application: crtview.NewApplication(),
			callback:    teaboxlib.NewTeaboxSocketServer(),
		}
	}
	return __APP_REF
}

func GetTeaboxQuitMessage() string {
	return __MSG_REF
}
