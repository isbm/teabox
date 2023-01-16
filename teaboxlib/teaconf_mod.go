package teaboxlib

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	wzlib_logger "github.com/infra-whizz/wzlib/logger"
)

type TeaConfCmdOption struct {
	label      string
	optionType string
	value      interface{}

	wzlib_logger.WzLogger
}

func NewTeaConfCmdOption(data interface{}) *TeaConfCmdOption {
	tcc := new(TeaConfCmdOption)
	if err := tcc.Parse(data); err != nil {
		tcc.GetLogger().Error(fmt.Sprintf("Unable to parse command option %v: %s", data, err.Error()))
		os.Exit(1)
	}

	return tcc
}

/*
Parse option. Option is written in two distinct forms:

 1. [[LABEL] [TYPE (string|bool)] <VALUE>]
 2. []

In case it contains a label and its type to the option, then is used in complex
widgets, like "dropdown" or "list" or "text input". If there is no label and a type, then
it is usually a toggle, because it has no label to the very option. It is just a checkbox
or an entry text. In this case only value is used as a string.

If it is entirely empty, then it is a checkbox/toggle and used as an ON/OFF flag. For example,
one wants to set "-d" or "--yes" flags, so they write in yaml an empty value of options "- []".
*/
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

	switch len(tokens) {
	case 0:
		tokens = []string{"", "", ""}
	case 1:
		tokens = append([]string{"", "string"}, tokens...)
	case 2:
		tokens = append([]string{""}, tokens...)
	case 3:
	default:
		return fmt.Errorf("unknown format to the option: %v", option)
	}

	for idx, t := range tokens {
		switch idx {
		case 0:
			tco.label = t
		case 1:
			tco.optionType = t
		case 2:
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

// GetValueAsString returns value forced to a string
func (tco *TeaConfCmdOption) GetValueAsString() string {
	return fmt.Sprintf("%v", tco.value)
}

// TeaConfModArg is an argument to the command. It has options of choice and is labeled accordingly.
// For example, "--path" can have different pre-set choices or user can enter his own value.
type TeaConfModArg struct {
	// Name of the argument, like "--path=....", "--arch=..." etc
	name string

	// Widget type, such as list, dropdown, text, toggle etc
	argtype string

	// Label on the form for this widget
	label string

	// Attributes of the argument
	attrs *TeaConfArgAttributes

	// Preset options. They can be also loaded dynamically via socket
	options []*TeaConfCmdOption

	wzlib_logger.WzLogger
}

func NewTeaConfModArg(args map[interface{}]interface{}) *TeaConfModArg {
	a := new(TeaConfModArg)
	a.options = []*TeaConfCmdOption{}

	optbuf := []interface{}{}

	for wn, wd := range args {
		switch wn.(string) {
		case "type":
			a.argtype = wd.(string)
		case "label":
			a.label = wd.(string)
		case "options":
			// Postpone opts parse
			optbuf = wd.([]interface{})
		case "name":
			a.name = wd.(string) // Add as-is. If it is with double-dash, then it is so.
		case "attributes":
			attrs, _ := wd.([]interface{}) // Avoid explicit cast crash. If syntax is wrong, then just skip it by passing nil.
			a.attrs = NewTeaConfArgAttributes(attrs)
		}
	}

	// Parse opts
	if a.argtype != "tabular" {
		for _, opt := range optbuf {
			a.options = append(a.options, NewTeaConfCmdOption(opt))
		}
	} else {
		a.options = NewTeaConfTabularData(optbuf).MakeOptionsData()
	}

	if a.argtype == "" {
		a.GetLogger().Error("No type found for argument.")
		a.GetLogger().Debug(spew.Sdump(args))
		os.Exit(1)
	} else if a.label == "" && a.argtype != "silent" {
		a.GetLogger().Error("No label found for argument.")
		a.GetLogger().Debug(spew.Sdump(args))
		os.Exit(1)
	} else if len(a.options) == 0 {
		a.GetLogger().Error("No default options has been found for argument, but the widget is not dynamic.")
		a.GetLogger().Debug(spew.Sdump(args))
		os.Exit(1)
	} else if a.name == "" {
		a.GetLogger().Error("No name found for argument.")
		a.GetLogger().Debug(spew.Sdump(args))
		os.Exit(1)
	}

	return a
}

func (a *TeaConfModArg) GetWidgetType() string {
	return a.argtype
}

func (a *TeaConfModArg) GetWidgetLabel() string {
	return a.label
}

// GetAttrs returns argument extra attributes
func (a *TeaConfModArg) GetAttrs() *TeaConfArgAttributes {
	return a.attrs
}

func (a *TeaConfModArg) GetOptions() []*TeaConfCmdOption {
	return a.options
}

// GetArgName is a name of an argument target, e.g. "--path".
// Name is unchanged, because various commands can have anything.
func (a *TeaConfModArg) GetArgName() string {
	return a.name
}

// TeaConfModCommand is a command within a module chain. It will call whatever command, specified in its path
// and passed arguments (TeaConfModArg).
type TeaConfModCommand struct {
	path      string
	title     string
	option    string // If not empty, then the command is optional and requires a yes/no before proceed.
	arguments []*TeaConfModArg
	flags     []string
}

func NewTeaConfModCommand(cmd map[interface{}]interface{}) *TeaConfModCommand {
	tcm := new(TeaConfModCommand)
	tcm.parse(cmd)
	return tcm
}

func (tmc *TeaConfModCommand) parse(cmd map[interface{}]interface{}) *TeaConfModCommand {
	tmc.arguments = []*TeaConfModArg{}
	tmc.flags = []string{}

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
			if sk == "flags" {
				for _, flag := range v.([]interface{}) {
					tmc.flags = append(tmc.flags, flag.(string))
				}
			} else if sk == "args" {
				for _, arg := range v.([]interface{}) {
					tmc.arguments = append(tmc.arguments, NewTeaConfModArg(arg.(map[interface{}]interface{})))
				}
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

// GetArguments returns module args
func (tmc *TeaConfModCommand) GetArguments() []*TeaConfModArg {
	return tmc.arguments
}

// GetStaticFlags returns static flags
func (tmc *TeaConfModCommand) GetStaticFlags() []string {
	return tmc.flags
}

// TeaConfModule is a wrapper for UI to shape a correct arguments to a target executable, call it, interact with it
// and provide results for futher processing within a chain.
type TeaConfModule struct {
	socketPath string
	landing    string
	setup      string
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

func (tcf *TeaConfModule) SetSetupCommand(setup string) *TeaConfModule {
	tcf.setup = setup
	return tcf
}

func (tcf *TeaConfModule) GetSetupCommand() string {
	if tcf.setup != "" {
		return strings.Fields(tcf.setup)[0]
	}

	return ""
}

func (tcf *TeaConfModule) GetSetupCommandArgs() []string {
	if tcf.setup != "" {
		return strings.Fields(tcf.setup)[1:]
	}
	return []string{}
}

// SetLandingPageType is one of "logger", "progress", "list" etc.
func (tcf *TeaConfModule) SetLandingPageType(lp string) *TeaConfModule {
	switch lp {
	case "logger", "progress", "list":
		tcf.landing = lp
	case "":
		tcf.landing = "logger"
	default:
		fmt.Printf("Error: unknown landing page ID in module \"%s\"\n", tcf.GetTitle())
		os.Exit(1)
	}
	return tcf
}

func (tcf *TeaConfModule) GetLandingPageType() string {
	return tcf.landing
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
