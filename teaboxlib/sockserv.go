package teaboxlib

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TeaboxSocketListener struct {
	addr    string
	conn    net.Listener
	actions []func(*TeaboxAPICall)
}

func NewTeaboxSocketListener(pth string) *TeaboxSocketListener {
	tsl := new(TeaboxSocketListener)
	tsl.addr = pth
	tsl.actions = []func(*TeaboxAPICall){}
	return tsl
}

func (tsl *TeaboxSocketListener) AddActions(action ...func(*TeaboxAPICall)) *TeaboxSocketListener {
	tsl.actions = append(tsl.actions, action...)
	return tsl
}

// Cleanup all the same unix socket addresses prior or after
func (tsl *TeaboxSocketListener) Cleanup() error {
	return os.RemoveAll(tsl.addr)
}

func (tsl *TeaboxSocketListener) Connect() error {
	if tsl.conn != nil {
		return fmt.Errorf("Cannot connect twice to the same socket!")
	}

	var err error
	tsl.conn, err = net.Listen("unix", tsl.addr)
	return err
}

func (tsl *TeaboxSocketListener) Start() error {
	for {
		if tsl.conn == nil {
			break
		}

		bind, err := tsl.conn.Accept()
		if err != nil {
			return err
		}

		// Go over all registered calls and send them the API instruction calls
		go func() {
			var call *TeaboxAPICall
			var data string

			buff := bytes.NewBufferString(data)
			if _, err := io.Copy(buff, bind); err == nil {
				call = NewTeaboxAPICall(buff.Bytes())
			}

			for _, a := range tsl.actions {
				a(call)
			}
		}()
	}

	return nil
}

func (tsl *TeaboxSocketListener) Terminate() error {
	if err := tsl.conn.Close(); err != nil {
		return err
	}

	tsl.conn = nil
	return nil
}

/*
Unix Socket Server, runs the listener above and terminates it when required.
This server also registers new listeners in Start(path).
*/
type TeaboxSocketServer struct {
	mtx           bool
	listener      *TeaboxSocketListener
	localActions  []func(*TeaboxAPICall)
	globalActions []func(*TeaboxAPICall)
}

func NewTeaboxSocketServer() *TeaboxSocketServer {
	tss := new(TeaboxSocketServer)
	tss.localActions = []func(*TeaboxAPICall){}
	tss.globalActions = []func(*TeaboxAPICall){}
	return tss
}

// AddLocalActions are actions that are added per a specific widget implementation. They define their specific APIs.
func (tss *TeaboxSocketServer) AddLocalAction(action func(call *TeaboxAPICall)) *TeaboxSocketServer {
	for {
		if !tss.mtx {
			break
		}
		time.Sleep(time.Millisecond)
	}

	tss.mtx = true
	// Merge action after start
	if tss.listener != nil {
		tss.listener.AddActions(action)
	} else {
		tss.localActions = append(tss.localActions, action)
	}
	tss.mtx = false

	return tss
}

// AddGlobalAction adds an action to the Unix socket server, that will be persisting per entire instance.
// This action is accessed by all forms at all runtime span. This is useful at application start to define
// all common actions.
func (tss *TeaboxSocketServer) AddGlobalAction(action func(*TeaboxAPICall)) *TeaboxSocketServer {
	tss.globalActions = append(tss.globalActions, action)
	return tss
}

// Start the Unix socket Server
func (tss *TeaboxSocketServer) Start(pth string) error {
	if len(tss.localActions) == 0 && len(tss.globalActions) == 0 {
		return fmt.Errorf("No any actions were assigned yet")
	}
	tss.listener = NewTeaboxSocketListener(pth).AddActions(tss.globalActions...).AddActions(tss.localActions...)
	if err := tss.listener.Cleanup(); err != nil {
		return err
	}
	if err := tss.listener.Connect(); err != nil {
		return err
	}

	go tss.listener.Start()

	return nil
}

// Stop the Unix socket server
func (tss *TeaboxSocketServer) Stop() error {
	defer func() {
		tss.listener = nil
		tss.localActions = []func(*TeaboxAPICall){}
	}()
	if err := tss.listener.Terminate(); err != nil {
		return err
	}
	return tss.listener.Cleanup()
}

// IsRunning checks if Unix socket server is running
func (tss *TeaboxSocketServer) IsRunning() bool {
	return tss.listener != nil
}
