package teawidgets

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox/teaboxlib"
)

type TeaboxArgsMainWindow struct {
	cmdId, title, subtitle string
	flags                  []string
	argset                 map[string]string // map of strings for named arguments.
	argindex               []string          // an array of named arguments for args ordering.

	*crtview.Form
}

func NewTeaboxArgsMainWindow(title, subtitle string) *TeaboxArgsMainWindow {
	return (&TeaboxArgsMainWindow{
		Form:     crtview.NewForm(),
		title:    title,
		subtitle: subtitle,
		flags:    []string{},
		argset:   map[string]string{},
		argindex: []string{},
	}).init()
}

func (tmw *TeaboxArgsMainWindow) init() *TeaboxArgsMainWindow {
	// XXX: ID (Set/Get by ID) needs to be re-thought
	tmw.SetTitle(fmt.Sprintf("%s - %s", tmw.title, tmw.subtitle))
	tmw.cmdId = fmt.Sprintf("%s - %s", tmw.title, tmw.subtitle)

	tmw.SetBorder(true)

	// Colors
	tmw.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	tmw.SetFieldTextColor(tcell.ColorWhite)
	tmw.SetFieldBackgroundColor(teaboxlib.WORKSPACE_HEADER)
	tmw.SetFieldBackgroundColorFocused(tcell.ColorGreenYellow)
	tmw.SetFieldTextColorFocused(teaboxlib.WORKSPACE_BACKGROUND)

	tmw.SetButtonBackgroundColor(teaboxlib.FORM_BUTTON_BACKGROUND)
	tmw.SetButtonBackgroundColorFocused(teaboxlib.FORM_BUTTON_BACKGROUND_SELECTED)
	tmw.SetButtonTextColor(teaboxlib.FORM_BUTTON_TEXT)
	tmw.SetButtonTextColorFocused(teaboxlib.FORM_BUTTON_TEXT_SELECTED)

	// Buttons align
	tmw.SetButtonsAlign(crtview.AlignRight)
	tmw.SetButtonsToBottom()

	return tmw
}

func (tmw *TeaboxArgsMainWindow) GetId() string {
	return tmw.cmdId
}

// AddFlag adds a flag to the CLI command per a form.
func (tmw *TeaboxArgsMainWindow) AddFlag(formid, flag string) *TeaboxArgsMainWindow {
	if flag == "" {
		return tmw
	}

	for _, f := range tmw.flags {
		if f == flag { // already set
			return tmw
		}
	}

	// Add a flag
	tmw.flags = append(tmw.flags, flag)

	return tmw
}

// RemoveFlag removes a flag from the CLI command per a form.
func (tmw *TeaboxArgsMainWindow) RemoveFlag(formid, flag string) *TeaboxArgsMainWindow {
	nf := []string{}
	for _, f := range tmw.flags {
		if flag != f {
			nf = append(nf, f)
		}
	}

	tmw.flags = nf

	return tmw
}

func (tmw *TeaboxArgsMainWindow) SetStaticFlags(cmd *teaboxlib.TeaConfModCommand) *TeaboxArgsMainWindow {
	tmw.flags = append(tmw.flags, cmd.GetStaticFlags()...)
	return tmw
}

func (tmw *TeaboxArgsMainWindow) GetFlags() []string {
	return tmw.flags
}

// AddArgument adds an argument to the CLI command per a form. Repeating this function call
// will override the previous value (update).
func (tmw *TeaboxArgsMainWindow) AddArgument(formid, argname, argvalue string) *TeaboxArgsMainWindow {
	isNew := true
	for _, x := range tmw.argindex {
		if x == argname {
			isNew = false
			break
		}
	}

	if isNew {
		tmw.argindex = append(tmw.argindex, argname)
	}

	tmw.argset[argname] = argvalue

	return tmw
}

// RemoveArgument sets an argument to the CLI command per a form
func (tmw *TeaboxArgsMainWindow) RemoveArgument(formid, argname string) *TeaboxArgsMainWindow {
	keys := []string{}
	for _, key := range tmw.argindex {
		if key != argname {
			keys = append(keys, key)
		}
	}

	tmw.argindex = keys
	delete(tmw.argset, argname)

	return tmw
}

// GetCommandArguments returns an array of strings in a form of a formed command line, like so:
//
//	[]string{"-x", "-y", "-z", "--path=/dev/null"}
//
// All the data is ordered as it is described in the module configuration.
func (tmw *TeaboxArgsMainWindow) GetCommandArguments(formid string) []string {
	// Get flags for this form, if any
	cargs := append([]string{}, tmw.flags...) // copy

	// Get ordered arguments, if any
	for _, arg := range tmw.argindex {
		val := tmw.argset[arg]
		if val != "" {
			val = fmt.Sprintf("%s=%s", arg, val)
		} else {
			val = arg
		}
		cargs = append(cargs, val)
	}

	return cargs
}

// AddArgWidgets adds actual widgets for each argument
func (tmw *TeaboxArgsMainWindow) AddArgWidgets(cmd *teaboxlib.TeaConfModCommand) {
	for _, a := range cmd.GetArguments() {
		switch a.GetWidgetType() {
		case "dropdown", "list":
			tmw.AddDropDownSimple(a)
		case "text":
			tmw.AddInputField(a)
		case "toggle":
			tmw.AddCheckBox(a)
		case "silent": // XXX: should be deprecated from the specs
		}
	}
}

func (tmw *TeaboxArgsMainWindow) AddDropDownSimple(arg *teaboxlib.TeaConfModArg) error {
	opts := []string{}
	for _, opt := range arg.GetOptions() {
		if v, _ := opt.GetValue().(string); v != "" {
			opts = append(opts, v)
		}
	}
	if len(opts) == 0 {
		return fmt.Errorf("List \"%s\" in command \"%s\" of module \"%s\" has no values.", arg.GetWidgetLabel(), tmw.subtitle, tmw.title)
	}

	tmw.Form.AddDropDownSimple(arg.GetWidgetLabel(), 0, func(index int, option *crtview.DropDownOption) {
		//taf.AddArgument(tf.GetId(), arg.GetArgName(), strings.TrimSpace(option.GetText()))
	}, opts...)
	return nil
}

/*
Text could have only one argument as a default text:

	[DEFAULT_TEXT]

The field can be also completely empty.
*/
func (tmw *TeaboxArgsMainWindow) AddInputField(arg *teaboxlib.TeaConfModArg) error {
	if len(arg.GetOptions()) > 0 {
		val := arg.GetOptions()[0].GetValueAsString()
		if val != "" {
			// register the value, if any
			tmw.AddArgument(tmw.GetId(), arg.GetArgName(), val)
		}

		tmw.Form.AddInputField(arg.GetWidgetLabel(), val, 0, nil, func(text string) {
			tmw.AddArgument(tmw.GetId(), arg.GetArgName(), strings.TrimSpace(text))
		})
	}

	return nil
}

func (tmw *TeaboxArgsMainWindow) AddCheckBox(arg *teaboxlib.TeaConfModArg) error {
	if len(arg.GetOptions()) == 0 {
		return fmt.Errorf("Toggle \"%s\" in command \"%s\" of module \"%s\" should have its default state with at least one option.", arg.GetWidgetLabel(), tmw.title, tmw.subtitle)
	}

	state, ok := arg.GetOptions()[0].GetValue().(bool) // This is not the *value* but checked/unchecked state
	if ok && state {
		// Register state
		tmw.AddArgument(tmw.GetId(), arg.GetArgName(), arg.GetOptions()[0].GetLabel())
	}

	tmw.Form.AddCheckBox(arg.GetWidgetLabel(), "", state, func(checked bool) {
		if checked {
			tmw.AddArgument(tmw.GetId(), arg.GetArgName(), arg.GetOptions()[0].GetLabel())
		} else {
			tmw.RemoveArgument(tmw.GetId(), arg.GetArgName())
		}
	})

	return nil
}
