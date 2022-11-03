package teaboxlib

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TeaConfCmdOption struct {
	name       string
	label      string
	optionType string
	value      interface{}
}

func NewTeaConfCmdOption(data interface{}) *TeaConfCmdOption {
	tcc := new(TeaConfCmdOption)
	if err := tcc.Parse(data); err != nil {
		panic(fmt.Sprintf("Unable to parse command option: %v", data))
	}

	return tcc
}

// Parse option. It is always written in the following form:
// - [<ARG/OPT NAME> <LABEL> <TYPE (string|bool) [VALUE]]
//
// NOTE: Currently method never returns anything else but nil as an error,
// silently yanking wrong values. Maybe it is not the best behaveiour, will see.
func (tco *TeaConfCmdOption) Parse(option interface{}) error {
	var tokens []string
	aopt, ok := option.([]interface{})
	if ok {
		for _, opt := range aopt {
			tokens = append(tokens, fmt.Sprintf("%v", opt)) // force to string at this point
		}
	} else {
		tokens = strings.Fields(strings.TrimSpace(fmt.Sprintf("%v", option)))
	}

	for idx, t := range tokens {
		switch idx {
		case 0:
			tco.name = t
		case 1:
			tco.label = t
		case 2:
			tco.optionType = t
		case 3:
			switch tco.optionType {
			case "bool":
				tco.value = t == "yes" || t == "true"
			case "int":
				v, e := strconv.Atoi(t)
				if e == nil {
					tco.value = v
				} else {
					tco.value = t
				}
			default:
				tco.value = t
			}
		}
	}
	return nil
}

// GetName the option name, like "--help", "--with-somthing" etc.
// This should include full option with all the required dashes
func (tco *TeaConfCmdOption) GetName() string {
	return tco.name
}

// GetLabel of the option as a title to the widget, summry and explanation
func (tco *TeaConfCmdOption) GetLabel() string {
	return tco.label
}

// GetType of the option (bool, int, string)
func (tco *TeaConfCmdOption) GetType() string {
	return tco.optionType
}

// GetValue returns parsed type as boolean, integer or a string
func (tco *TeaConfCmdOption) GetValue() interface{} {
	return tco.value
}

// TeaConfModArg is an argument to the command. It has options of choice and is labeled accordingly.
// For example, "--path" can have different pre-set choices or user can enter his own value.
type TeaConfModArg struct {
	argtype string
	label   string
	options []*TeaConfCmdOption
}

func NewTeaConfModArg(args map[interface{}]interface{}) *TeaConfModArg {
	a := new(TeaConfModArg)
	a.options = []*TeaConfCmdOption{}

	for wn, wd := range args {
		switch wn.(string) {
		case "type":
			a.argtype = wd.(string)
		case "label":
			a.label = wd.(string)
		case "options":
			for _, opt := range wd.([]interface{}) {
				a.options = append(a.options, NewTeaConfCmdOption(opt))
			}
		}
	}
	return a
}

func (a *TeaConfModArg) GetWidgetType() string {
	return a.argtype
}

func (a *TeaConfModArg) GetWidgetLabel() string {
	return a.label
}

func (a *TeaConfModArg) GetOptions() []*TeaConfCmdOption {
	return a.options
}

// TeaConfModCommand is a command within a module chain. It will call whatever command, specified in its path
// and passed arguments (TeaConfModArg).
type TeaConfModCommand struct {
	path      string
	title     string
	option    string // If not empty, then the command is optional and requires a yes/no before proceed.
	arguments []*TeaConfModArg
}

func NewTeaConfModCommand(cmd map[interface{}]interface{}) *TeaConfModCommand {
	tcm := new(TeaConfModCommand)
	tcm.parse(cmd)
	return tcm
}

func (tmc *TeaConfModCommand) parse(cmd map[interface{}]interface{}) *TeaConfModCommand {
	tmc.arguments = []*TeaConfModArg{}
	var failed string
	var obj interface{}
	for k, v := range cmd {
		sk, ok := k.(string)
		if !ok {
			panic(fmt.Sprintf("Parse command failure: option %v should be a string type", k))
		}

		if sv, _ := v.(string); sv != "" {
			switch sk {
			case "path":
				tmc.path = sv
			case "title":
				tmc.title = sv
			case "option":
				tmc.option = sv
			}
		} else if varg, _ := v.([]interface{}); varg != nil {
			for _, arg := range v.([]interface{}) {
				tmc.arguments = append(tmc.arguments, NewTeaConfModArg(arg.(map[interface{}]interface{})))
			}
		} else {
			if failed == "" {
				failed = fmt.Sprintf("%v", k)
				obj = v
			}
		}
	}

	if failed != "" {
		fmt.Printf("Unknown syntax on commands configuration near \"%v\" of command \"%s\"\n", failed, tmc.title)
		fmt.Println("Value: ", obj)
		os.Exit(1)
	}

	return nil
}

func (tmc *TeaConfModCommand) GetTitle() string {
	return tmc.title
}

func (tmc *TeaConfModCommand) GetCommandPath() string {
	return tmc.path
}

func (tmc *TeaConfModCommand) SetCommandPath(p string) {
	tmc.path = p
}

// If this returns non-empty string, then the command is optional and this string is the message.
func (tmc *TeaConfModCommand) GetOptionLabel() string {
	return tmc.option
}

func (tmc *TeaConfModCommand) GetArguments() []*TeaConfModArg {
	return tmc.arguments
}

// TeaConfModule is a wrapper for UI to shape a correct arguments to a target executable, call it, interact with it
// and provide results for futher processing within a chain.
type TeaConfModule struct {
	socketPath string
	conditions []map[string]string
	commands   []*TeaConfModCommand

	TeaConfBaseEntity
}

func NewTeaConfModule(title string) *TeaConfModule {
	tcm := new(TeaConfModule)
	tcm.SetTitle(title)
	tcm.etype = "module"

	return tcm
}

// SetCallbackPath sets a physical path on the disk for the Unix socket to communicate between the processes.
func (tcf *TeaConfModule) SetCallbackPath(pt interface{}) *TeaConfModule {
	if v, ok := pt.(string); ok {
		tcf.socketPath = v
	}
	return tcf
}

// GetCallbackPath returns a physical path on the disk for the Unix socket to communicate between the processes.
func (tcf *TeaConfModule) GetCallbackPath() string {
	return tcf.socketPath
}

// SetCondition sets described conditions, under which module is running or not.
func (tcf *TeaConfModule) SetCondition(cond interface{}) *TeaConfModule {
	if cond == nil {
		return tcf
	}

	tcf.conditions = []map[string]string{}
	condset, ok := cond.([]interface{})
	if !ok {
		panic("Wrong configuration of the module: " + tcf.title)
	}

	for _, icnd := range condset {
		cnd := map[string]string{}
		imcnd := icnd.(map[interface{}]interface{})
		for k, v := range imcnd {
			ks, ok := k.(string)
			if !ok {
				panic(fmt.Sprintf("Wrong configuration of the module %s: key %v is not string", tcf.title, k))
			}
			vs, ok := v.(string)
			if !ok {
				panic(fmt.Sprintf("Wrong configuration of the module %s: value %v is not string", tcf.title, v))
			}

			cnd[ks] = vs
		}
		tcf.conditions = append(tcf.conditions, cnd)
	}

	return tcf
}

func (tcf *TeaConfModule) GetCondition() []map[string]string {
	return tcf.conditions
}

func (tcf *TeaConfModule) SetCommands(commands interface{}) *TeaConfModule {
	if commands == nil {
		return tcf
	}
	tcf.commands = []*TeaConfModCommand{}

	cmds, ok := commands.([]interface{})
	if !ok {
		panic(fmt.Sprintf("Wrong configuration of commands at module %s", tcf.title))
	}

	for _, cmddata := range cmds {
		cmd, ok := cmddata.(map[interface{}]interface{})
		if !ok {
			panic(fmt.Sprintf("Wrong configuration of commands at module %s", tcf.title))
		}
		tcf.commands = append(tcf.commands, NewTeaConfModCommand(cmd))
	}

	return tcf
}

func (tcf *TeaConfModule) GetCommands() []*TeaConfModCommand {
	return tcf.commands
}
