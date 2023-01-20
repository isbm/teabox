package teabox

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
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

func DumpToFile(fn string, obj ...interface{}) {
	ioutil.WriteFile(fn, []byte(spew.Sdump(obj...)), 0644)
}

// AddToFile a text. Used for quick debug purposes where using a real debugger
// is too tedious and time consuming. :-)
func AddToFile(fn string, data string) error {
	f, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()
	for _, chunk := range []string{"--------------", time.Now().Format("2006-01-02 15:04:05"), data} {
		if _, err := f.WriteString(fmt.Sprintf("%v\n", strings.TrimSpace(chunk))); err != nil {
			return err
		}
	}

	return nil
}

// DebugToFile writes a data to a teabox-trace.log in the current directory.
// Since screen is locked, printing to the STDOUT is not an option, so
// this prints to a file, which one can watch with "tail -f" elsewhere.
func DebugToFile(data string) error {
	return AddToFile("teabox-trace.log", data)
}
