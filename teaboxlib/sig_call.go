package teaboxlib

import (
	"fmt"
	"os/exec"
	"path"
)

type SigCall struct {
	modCmd *TeaConfModCommand
}

// NewSigCall creates
func NewSigCall(mc *TeaConfModCommand) *SigCall {
	sc := new(SigCall)
	sc.modCmd = mc
	return sc
}

// CallSignal action of the widget (any)
func (sc *SigCall) CallSignal(act *TeaConfArgSignalAction) {
	if err := sc.call(act); err != nil {
		panic(fmt.Sprintf("Error while calling signal \"%s %v\": %s", act.GetName(), act.GetArguments(), err.Error()))
	}
}

func (sc *SigCall) call(act *TeaConfArgSignalAction) error {
	if sc.modCmd == nil {
		return fmt.Errorf("error signal call: undefined module configuration")
	}

	if act.GetName() == "" { // Undefined call
		return nil
	}

	pth := path.Dir(sc.modCmd.GetCommandPath())
	out, err := exec.Command(path.Join(pth, act.GetName()), act.GetArguments()...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error: %v\n%v", err.Error(), out)
	}

	return nil
}
