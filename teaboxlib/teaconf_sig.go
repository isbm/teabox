package teaboxlib

import (
	"fmt"
	"strings"
)

type TeaConfArgSignalAction struct {
	name string
	args []string
}

// Signal action takes raw string and parses it into an action with the whole verfication
func NewTeaConfArgSignalAction(action string) (*TeaConfArgSignalAction, error) {
	act := new(TeaConfArgSignalAction)
	act.args = []string{}

	if action == "" {
		return act, nil
	}

	err := act.parse(action)
	if err != nil {
		return nil, err
	}

	return act, nil
}

// GetName of the command to be called
func (act *TeaConfArgSignalAction) GetName() string {
	return act.name
}

// GetArguments of the command
func (act *TeaConfArgSignalAction) GetArguments() []string {
	return act.args
}

// Parse the signal action
func (act *TeaConfArgSignalAction) parse(action string) error {
	cmd := []string{}
	for _, c := range strings.Split(action, " ") {
		c = strings.TrimSpace(c)
		if c != "" {
			cmd = append(cmd, c)
		}
	}
	if len(cmd) < 1 {
		return fmt.Errorf("could not parse signal action")
	}

	act.name = cmd[0]
	if len(cmd) > 1 {
		act.args = cmd[1:]
	}

	return nil
}

// TeaConfArgSignals is a container of all signals, defined per a widget argument
type TeaConfArgSignals struct {
	sigs map[string]*TeaConfArgSignalAction
}

// NewTeaConfArgSignlans constructor
func NewTeaConfArgSignals(data map[interface{}]interface{}) (*TeaConfArgSignals, error) {
	tcsig := new(TeaConfArgSignals)
	tcsig.sigs = map[string]*TeaConfArgSignalAction{}

	if data == nil {
		data = map[interface{}]interface{}{}
	}

	for k, v := range data {
		action, err := NewTeaConfArgSignalAction(fmt.Sprintf("%v", v))
		if err != nil {
			return nil, err
		} else {
			tcsig.SetSignal(fmt.Sprintf("%v", k), action)
		}
	}

	return tcsig, nil
}

// GetSignals defined to the specific widget of an argument
func (tcsig *TeaConfArgSignals) GetSignals() []string {
	sigs := []string{}
	for sigdef := range tcsig.sigs {
		sigs = append(sigs, sigdef)
	}

	return sigs
}

// GetSignalValue returns a defined command on a specific signal
func (tcsig *TeaConfArgSignals) GetSignalValue(sigdef string) *TeaConfArgSignalAction {
	sig, ok := tcsig.sigs[sigdef]
	if !ok {
		dummy, _ := NewTeaConfArgSignalAction("")
		return dummy
	}

	return sig
}

// SetSignal sets a signal to a slot. If such signal previously set, the new will be ignored.
func (tcsig *TeaConfArgSignals) SetSignal(sigdef string, value *TeaConfArgSignalAction) *TeaConfArgSignals {
	_, ok := tcsig.sigs[sigdef]
	if !ok { // set only when key does not exist yet
		tcsig.sigs[sigdef] = value
	}

	return tcsig
}
