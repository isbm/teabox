package teawidgets

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/isbm/crtview"
	"github.com/isbm/crtview/crtwin/crtforms"
	"gitlab.com/isbm/teabox"
	"gitlab.com/isbm/teabox/teaboxlib"
)

// Operations for the widget update. This is internal only here, don't use outside.
const (
	__OP_W_ADD = iota
	__OP_W_SET
	__OP_W_CLR
)

type TeaboxArgsMainWindow struct {
	cmdId, title, subtitle string
	flags                  []string
	argset                 map[string]string                   // map of strings for named arguments.
	argindex               []string                            // an array of named arguments for args ordering.
	labeledArg             map[string]*teaboxlib.TeaConfModArg // map of label to arg object pointer. Used to find argument name by label (same as FormItem)
	namedArg               map[string]*teaboxlib.TeaConfModArg
	skipLoad               bool
	confModCommand         *teaboxlib.TeaConfModCommand

	*crtview.Form
}

func NewTeaboxArgsMainWindow(title, subtitle string) *TeaboxArgsMainWindow {
	return (&TeaboxArgsMainWindow{
		Form:       crtview.NewForm(),
		title:      title,
		subtitle:   subtitle,
		flags:      []string{},
		argset:     map[string]string{},
		argindex:   []string{},
		labeledArg: map[string]*teaboxlib.TeaConfModArg{},
		namedArg:   map[string]*teaboxlib.TeaConfModArg{},
		skipLoad:   false,
	}).init()
}

func (tmw *TeaboxArgsMainWindow) init() *TeaboxArgsMainWindow {
	tmw.SetTitle(fmt.Sprintf("%s - %s", tmw.title, tmw.subtitle))
	tmw.cmdId = fmt.Sprintf("%s - %s", tmw.title, tmw.subtitle)

	tmw.SetBorder(true)

	// Colors
	tmw.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	tmw.SetFieldTextColor(tcell.ColorBlack.TrueColor())
	tmw.SetFieldBackgroundColor(teaboxlib.WORKSPACE_HEADER)
	tmw.SetFieldBackgroundColorFocused(tcell.ColorGreenYellow)
	tmw.SetFieldTextColorFocused(tcell.ColorBlack.TrueColor())

	// Buttons align
	tmw.SetButtonsAlign(crtview.AlignRight)
	tmw.SetButtonsToBottom(true)

	return tmw
}

func (tmw *TeaboxArgsMainWindow) SkipLoad() bool {
	return tmw.skipLoad
}

func (tmw *TeaboxArgsMainWindow) SetSkipLoad(action func(), message string) {
	tmw.skipLoad = true

	tmw.AddFormItem(crtforms.NewFormTextView().SetText(message))
	tmw.AddButton("Close", action)
	tmw.SetFocus(1) // Focus on 2nd element, i.e. button in this case
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
		// Skip argument, it is for view-only
		if tmw.namedArg[arg].GetAttrs().HasOption("view-only") {
			continue
		}

		// Maybe skip argument, depending how on value conditions
		val := tmw.argset[arg]
		if val != "" {
			val = fmt.Sprintf("%s=%s", arg, val)
		} else if tmw.namedArg[arg].GetAttrs() == nil || !tmw.namedArg[arg].GetAttrs().HasOption("skip-empty") {
			val = arg
		}

		// Add agreed argument
		if val != "" {
			cargs = append(cargs, val)
		}
	}

	return cargs
}

// AddArgWidgets adds actual widgets for each argument
func (tmw *TeaboxArgsMainWindow) AddArgWidgets(cmd *teaboxlib.TeaConfModCommand) {
	tmw.confModCommand = cmd
	for _, a := range tmw.confModCommand.GetArguments() {
		tmw.namedArg[a.GetArgName()] = a
		tmw.labeledArg[a.GetWidgetLabel()] = a
		switch a.GetWidgetType() {
		case "dropdown", "list":
			tmw.AddDropDownSimple(a)
		case "text":
			tmw.AddInputField(a)
		case "toggle":
			tmw.AddCheckBox(a)
		case "tabular":
			tmw.AddTabularField(a)
		case "password", "masked":
			tmw.AddPasswordField(a)
		default:
			fmt.Printf("Module config error: Unknown widget definition \"%s\" for command argument \"%s\" at %s\n", a.GetWidgetType(), cmd.GetTitle(), cmd.GetCommandPath())
			os.Exit(1)
		}
	}
}

// AddTabularFiled adds a list of tabular data. This is a complex field that has multiple columns.
func (tmw *TeaboxArgsMainWindow) AddTabularField(arg *teaboxlib.TeaConfModArg) error {
	if len(arg.GetOptions()) == 0 {
		return fmt.Errorf("no tabular data found for the field %s", arg.GetArgName())
	}

	rows := [][]string{}

	for _, rowData := range arg.GetOptions()[1:] {
		if rowData.GetType() == "tabular:row" {
			rows = append(rows, rowData.GetValue().(*teaboxlib.TeaConfTabularRow).GetLabels())
		}
	}

	tabular := crtforms.NewFormTabularChoice(arg.GetWidgetLabel(), arg.GetOptions()[0].GetValue().(*teaboxlib.TeaConfTabularRow).GetLabels(),
		rows, arg.GetAttrs().HasOption("selector"), arg.GetAttrs().KeywordValueAsInts("hidden")...).
		SetFieldHeight(arg.GetAttrs().KeywordValueAsInt("height")).
		SetMultiselect(arg.GetAttrs().HasOption("multiselect")).
		SetHasSearch(arg.GetAttrs().HasOption("search")).
		SetValueColumn(arg.GetAttrs().KeywordValueAsInt("value")).
		SetExpandingColumn(arg.GetAttrs().KeywordValueAsInt("expand"))
	tabular.SetFocusedBorderStyle(crtview.BorderSingle)
	tabular.SetTitleWhitespace(true)
	tabular.SetBorderColorFocused(tmw.GetAttributes().FieldBackgroundColorFocused)
	tabular.SetSelectedFunc(func(row, column int) {
		tmw.AddArgument(tmw.GetId(), arg.GetArgName(), tabular.GetValueAt(row-1))
	})
	tabular.Select(1, 1)
	tmw.Form.AddFormItem(tabular)

	// Pre-select first value
	tmw.AddArgument(tmw.GetId(), arg.GetArgName(), tabular.GetValueAt(0))

	return nil
}

func (tmw *TeaboxArgsMainWindow) AddDropDownSimple(arg *teaboxlib.TeaConfModArg) error {
	opts := []string{}
	for _, opt := range arg.GetOptions() {
		if v, _ := opt.GetValue().(string); v != "" {
			opts = append(opts, v)
		}
	}
	if len(opts) == 0 {
		return fmt.Errorf("list \"%s\" in command \"%s\" of module \"%s\" has no values", arg.GetWidgetLabel(), tmw.subtitle, tmw.title)
	}

	tmw.Form.AddDropDownSimple(arg.GetWidgetLabel(), 0, func(index int, option *crtview.DropDownOption) {
		tmw.AddArgument(tmw.GetId(), arg.GetArgName(), strings.TrimSpace(option.GetText()))
	}, opts...)
	return nil
}

/*
Text could have only one argument as a default text:

	[DEFAULT_TEXT]

The field can be also completely empty.
*/
func (tmw *TeaboxArgsMainWindow) AddInputField(arg *teaboxlib.TeaConfModArg) error {
	// XXX: this will omit the whole widget, if there no default value
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

func (tmw *TeaboxArgsMainWindow) AddPasswordField(arg *teaboxlib.TeaConfModArg) error {
	var val string = ""
	if len(arg.GetOptions()) > 0 {
		val = arg.GetOptions()[0].GetValueAsString()
	}

	if val != "" {
		tmw.AddArgument(tmw.GetId(), arg.GetArgName(), val)
	}

	tmw.Form.AddPasswordField(arg.GetWidgetLabel(), val, 0, '*', func(text string) {
		tmw.AddArgument(tmw.GetId(), arg.GetArgName(), text) // Don't trim space here :)
	})

	return nil
}

func (tmw *TeaboxArgsMainWindow) AddCheckBox(arg *teaboxlib.TeaConfModArg) error {
	if len(arg.GetOptions()) == 0 {
		return fmt.Errorf("toggle \"%s\" in command \"%s\" of module \"%s\" should have its default state with at least one option", arg.GetWidgetLabel(), tmw.title, tmw.subtitle)
	}

	state, ok := arg.GetOptions()[0].GetValue().(bool) // This is not the *value* but checked/unchecked state
	if ok && state {
		// Register state
		tmw.AddArgument(tmw.GetId(), arg.GetArgName(), arg.GetOptions()[0].GetLabel())
	}

	tmw.Form.AddCheckBox(arg.GetWidgetLabel(), "", state, func(checked bool) {
		if checked {
			// Add argument notification
			tmw.AddArgument(tmw.GetId(), arg.GetArgName(), arg.GetOptions()[0].GetLabel())

			// Add selected signal
			if err := tmw.Call(arg.GetSignals().GetSignalValue("selected")); err != nil {
				teabox.DebugToFile(fmt.Sprintf("Error on \"%s\" while calling checkbox \"%s\" hook: %s",
					tmw.GetId(), arg.GetArgName(), err.Error()))
			}
		} else {
			// Remove argument notification
			tmw.RemoveArgument(tmw.GetId(), arg.GetArgName())

			// Add unselected signal
			if err := tmw.Call(arg.GetSignals().GetSignalValue("deselected")); err != nil {
				teabox.DebugToFile(fmt.Sprintf("Error on \"%s\" while calling checkbox \"%s\" hook: %s",
					tmw.GetId(), arg.GetArgName(), err.Error()))
			}
		}
	})

	return nil
}

// Call action
func (tmw *TeaboxArgsMainWindow) Call(act *teaboxlib.TeaConfArgSignalAction) error {
	if act.GetName() == "" { // Undefined call
		return nil
	}

	pth := path.Dir(tmw.confModCommand.GetCommandPath())
	out, err := exec.Command(path.Join(pth, act.GetName()), act.GetArguments()...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error: %v\n%v", err.Error(), out)
	}

	return nil
}

func (tmw *TeaboxArgsMainWindow) GetSocketAcceptAction() func(*teaboxlib.TeaboxAPICall) {
	return func(call *teaboxlib.TeaboxAPICall) {
		switch call.GetClass() {

		// Overwriting with the new values
		case teaboxlib.FORM_SET_BY_LABEL:
			tmw.updateField(call, tmw.GetFormItemByLabel(call.GetKey()), __OP_W_SET)
		case teaboxlib.FORM_SET_BY_ORD:
			tmw.updateField(call, tmw.GetFormItem(call.GetKeyAsInt()), __OP_W_SET)
		case teaboxlib.FORM_SET_TABLE_BY_ORD:
			tmw.updateField(call, tmw.GetFormItem(call.GetKeyAsInt()), __OP_W_SET)

		// Adding/merging new values
		case teaboxlib.FORM_ADD_BY_LABEL:
			tmw.updateField(call, tmw.GetFormItemByLabel(call.GetKey()), __OP_W_ADD)
		case teaboxlib.FORM_ADD_BY_ORD:
			tmw.updateField(call, tmw.GetFormItem(call.GetKeyAsInt()), __OP_W_ADD)
		case teaboxlib.FORM_ADD_TABLE_BY_ORD:
			tmw.updateField(call, tmw.GetFormItem(call.GetKeyAsInt()), __OP_W_SET)

		// Clearing/resetting
		case teaboxlib.FORM_CLR_BY_LABEL:
			tmw.updateField(call, tmw.GetFormItemByLabel(call.GetKey()), __OP_W_CLR)
		case teaboxlib.FORM_CLR_BY_ORD:
			tmw.updateField(call, tmw.GetFormItem(call.GetKeyAsInt()), __OP_W_CLR)
		case teaboxlib.FORM_CLR_TABLE_BY_ORD:
			tmw.updateField(call, tmw.GetFormItem(call.GetKeyAsInt()), __OP_W_SET)
		}
	}
}

func (tmw *TeaboxArgsMainWindow) updateField(call *teaboxlib.TeaboxAPICall, item crtview.FormItem, op int) {
	arg := tmw.labeledArg[item.GetLabel()]

	// Reset all supported fields
	switch field := item.(type) {
	case *crtview.InputField:
		switch op {
		case __OP_W_SET:
			field.SetText(call.GetValue().(string))
		case __OP_W_ADD:
			field.SetText(field.GetText() + call.GetValue().(string))
		case __OP_W_CLR:
			field.SetText("")
		}

	case *crtview.CheckBox:
		if op != __OP_W_CLR {
			field.SetChecked(call.GetBool())
			if call.GetBool() {
				tmw.AddArgument(tmw.GetId(), arg.GetArgName(), arg.GetOptions()[0].GetLabel())
			} else {
				tmw.RemoveArgument(tmw.GetId(), arg.GetArgName())
			}
		} else {
			// Clear sets to false
			field.SetChecked(false)
			tmw.RemoveArgument(tmw.GetId(), arg.GetArgName())
		}

	case *crtview.DropDown:
		opts := []*crtview.DropDownOption{}
		for _, opt := range strings.Split(call.GetString(), "|") {
			opts = append(opts, crtview.NewDropDownOption(strings.TrimSpace(opt)))
		}
		if len(opts) < 1 {
			return
		}

		switch op {
		case __OP_W_ADD:
			field.AddOptions(opts...)
		case __OP_W_SET:
			field.SetOptions(func(index int, option *crtview.DropDownOption) {
				tmw.AddArgument(tmw.GetId(), arg.GetArgName(), strings.TrimSpace(option.GetText()))
			}, opts...)
			// Pre-select a first visible option in a dropdown
			tmw.AddArgument(tmw.GetId(), arg.GetArgName(), opts[0].GetText())
		case __OP_W_CLR:
			field.SetOptions(nil, []*crtview.DropDownOption{}...)
			tmw.RemoveArgument(tmw.GetId(), arg.GetArgName())
		}

	case *crtforms.FormTabularChoice:
		if call.GetType() != "json" {
			teabox.GetTeaboxApp().Stop(fmt.Sprintf("table data requires two dimentional array (tabular)"+
				" in JSON format. Current type: %s", call.GetType()))
		}

		td := call.GetValue()
		if td == nil {
			teabox.GetTeaboxApp().Stop("tabular field requested for the update, but no JSON data was found")
		}

		var rdata []interface{}
		if err := json.Unmarshal([]byte(td.(string)), &rdata); err != nil {
			teabox.GetTeaboxApp().Stop(fmt.Sprintf("unable to parse tabular data: %s. Is the whole JSON payload is in one "+
				"single-quoted string and scalar values are quoted with a double-quotes?", err.Error()))
		}

		rows := [][]string{}
		for _, iRow := range rdata {
			row := []string{}
			for _, c := range iRow.([]interface{}) {
				row = append(row, c.(string))
			}
			rows = append(rows, row)
		}

		switch op {
		case __OP_W_ADD:
			field.AppendContent(rows)
		case __OP_W_CLR:
			field.Clear()
		case __OP_W_SET:
			field.ReplaceContent(rows)
		}

		if op != __OP_W_CLR {
			field.Select(1, 1)
			tmw.AddArgument(tmw.GetId(), arg.GetArgName(), field.GetValueAt(0))
		} else {
			tmw.RemoveArgument(tmw.GetId(), arg.GetArgName())
		}
	}
}
