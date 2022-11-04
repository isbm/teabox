package teaboxui

import (
	"fmt"
	"path"
	"strings"

	"github.com/gdamore/tcell/v2"
	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox"
	"gitlab.com/isbm/teabox/teaboxlib"
	"gitlab.com/isbm/teabox/teaboxlib/teaboxui/teawidgets"
)

// TeaForm is just an enhancement to crtview.Form to set its unique ID to find it out later.
type TeaForm struct {
	cmdId string
	*crtview.Form
}

func NewTeaForm() *TeaForm {
	return &TeaForm{
		Form: crtview.NewForm(),
	}
}

func (tf *TeaForm) SetId(id string) {
	tf.cmdId = strings.ToLower(strings.ReplaceAll(id, " ", "-"))
}

func (tf *TeaForm) GetId() string {
	return tf.cmdId
}

// TeaFormsPanel is a layer of windows, and it contains many TeaForm instances to switch between them.
type TeaFormsPanel struct {
	wout *teawidgets.TeaSTDOUTWindow
	*crtview.Panels
}

func NewTeaFormsPanel() *TeaFormsPanel {
	tfp := &TeaFormsPanel{
		Panels: crtview.NewPanels(),
	}

	tfp.wout = teawidgets.NewTeaSTDOUTWindow()
	tfp.AddPanel("_stdout", tfp.wout, true, false)

	return tfp
}

func (tfp *TeaFormsPanel) ShowStdoutWindow() {
	tfp.SetCurrentPanel("_stdout")
}

func (tfp *TeaFormsPanel) GetStdoutWindow() *teawidgets.TeaSTDOUTWindow {
	return tfp.wout
}

func (tfp *TeaFormsPanel) AddForm(title, subtitle string) *TeaForm {
	f := NewTeaForm()

	f.SetTitle(fmt.Sprintf("%s - %s", title, subtitle))
	f.SetId(title)

	f.SetBorder(true)

	// Colors
	f.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	f.SetFieldTextColor(tcell.ColorWhite)
	f.SetFieldBackgroundColor(teaboxlib.WORKSPACE_HEADER)
	f.SetFieldBackgroundColorFocused(tcell.ColorGreenYellow)
	f.SetFieldTextColorFocused(teaboxlib.WORKSPACE_BACKGROUND)

	f.SetButtonBackgroundColor(teaboxlib.FORM_BUTTON_BACKGROUND)
	f.SetButtonBackgroundColorFocused(teaboxlib.FORM_BUTTON_BACKGROUND_SELECTED)
	f.SetButtonTextColor(teaboxlib.FORM_BUTTON_TEXT)
	f.SetButtonTextColorFocused(teaboxlib.FORM_BUTTON_TEXT_SELECTED)

	// Buttons align
	f.SetButtonsAlign(crtview.AlignRight)
	f.SetButtonsToBottom()

	tfp.AddPanel(f.GetId(), f, true, tfp.GetPanelCount() == 1)

	return f
}

// TeaboxArgsForm contains a layers with TeaForms on it, also their output, intro screen, callback screens etc.
type TeaboxArgsForm struct {
	TeaboxBaseWindow
	layers *crtview.Panels

	/*
		This is not super-nice, but still, per each form (via GetId() method) it holds the following:

		  1. A map of all named arguments with their values, e.g. "--path=/dev/null" etc.
		  2. An array of named arguments without their values in order to put an order to the map above.
		  3. An array of flags, like "-d -v -x -y -z" etc. Currently flags are always passed first.

		Note, these are all arguments to the *module*, so if Teabox is not flexible enough to call a random
		userland command, one can always wrap it with a shell script. :-)

		NOTE: If you still think this is dumb, feel free to make it better and send your PR!
	*/
	argset   map[string]map[string]string // map of forms by id, each has a map of strings for named arguments.
	argindex map[string][]string          // map of forms by id, each has an array of named arguments for their ordering.
	flagset  map[string][]string          // map of forms by id, each has an array of flags. Flags are passed first.

	modCmdIndex map[string]*teaboxlib.TeaConfModCommand

	wzlib_logger.WzLogger
}

func NewTeaboxArgsForm() *TeaboxArgsForm {
	taf := new(TeaboxArgsForm)
	taf.modCmdIndex = map[string]*teaboxlib.TeaConfModCommand{}

	taf.argset = map[string]map[string]string{}
	taf.argindex = map[string][]string{}
	taf.flagset = map[string][]string{}

	return taf
}

func (taf *TeaboxArgsForm) GetWidget() crtview.Primitive {
	return taf.layers
}

func (taf *TeaboxArgsForm) ShowForm(id string) {
	taf.layers.SetCurrentPanel(id)
}

func (taf *TeaboxArgsForm) Init() TeaboxWindow {
	taf.layers = crtview.NewPanels()

	intro := crtview.NewTextView()
	intro.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	intro.SetBorder(true)
	intro.SetText("\n\n\n\nSelect option from the menu on the left")
	intro.SetTextAlign(crtview.AlignCenter)

	taf.layers.AddPanel("_intro-screen", intro, true, true)

	for _, mod := range teabox.GetTeaboxApp().GetGlobalConfig().GetModuleStructure() {
		taf.generateForms(mod)
	}

	return taf
}

// AddFlag adds a flag to the CLI command per a form.
func (taf *TeaboxArgsForm) AddFlag(formid, flag string) *TeaboxArgsForm {
	if flag == "" {
		return taf
	}

	if _, ok := taf.flagset[formid]; !ok {
		taf.flagset[formid] = []string{}
	}

	for _, f := range taf.flagset[formid] {
		if f == flag { // already set
			return taf
		}
	}

	// Add a flag
	taf.flagset[formid] = append(taf.flagset[formid], flag)

	return taf
}

// RemoveFlag removes a flag from the CLI command per a form.
func (taf *TeaboxArgsForm) RemoveFlag(formid, flag string) *TeaboxArgsForm {
	if _, ok := taf.flagset[formid]; !ok {
		taf.flagset[formid] = []string{}
		return taf // nothing to remove
	}

	nf := []string{}
	for _, f := range taf.flagset[formid] {
		if flag != f {
			nf = append(nf, f)
		}
	}

	taf.flagset[formid] = nf

	return taf
}

// AddArgument adds an argument to the CLI command per a form. Repeating this function call
// will override the previous value (update).
func (taf *TeaboxArgsForm) AddArgument(formid, argname, argvalue string) *TeaboxArgsForm {
	if _, ok := taf.argindex[formid]; !ok {
		taf.argindex[formid] = []string{}
		taf.argset[formid] = map[string]string{}
	}

	isNew := true
	for _, x := range taf.argindex[formid] {
		if x == argname {
			isNew = false
			break
		}
	}

	if isNew {
		taf.argindex[formid] = append(taf.argindex[formid], argname)
	}

	taf.argset[formid][argname] = argvalue

	return taf
}

// RemoveArgument sets an argument to the CLI command per a form
func (taf *TeaboxArgsForm) RemoveArgument(formid, argname string) *TeaboxArgsForm {
	if _, ok := taf.argindex[formid]; !ok {
		taf.argindex[formid] = []string{}
		taf.argset[formid] = map[string]string{}
		return taf // nothing to remove
	}
	keys := []string{}
	for _, key := range taf.argindex[formid] {
		if key != argname {
			keys = append(keys, key)
		}
	}
	taf.argindex[formid] = keys
	delete(taf.argset[formid], argname)

	return taf
}

// GetCommandArguments returns an array of strings in a form of a formed command line, like so:
//
//	[]string{"-x", "-y", "-z", "--path=/dev/null"}
//
// All the data is ordered as it is described in the module configuration.
func (taf *TeaboxArgsForm) GetCommandArguments(formid string) []string {
	cargs := []string{}

	// Get flags for this form, if any
	if fset, ok := taf.flagset[formid]; ok {
		cargs = append(cargs, fset...)
	}

	// Get ordered arguments, if any
	if aidx, ok := taf.argindex[formid]; ok { // this assumes argindex always contains all keys in the
		for _, arg := range aidx {
			val := taf.argset[formid][arg]
			if val != "" {
				val = fmt.Sprintf("%s=%s", arg, val)
			} else {
				val = arg
			}
			cargs = append(cargs, val)
		}
	}

	return cargs
}

// ShowIntroScreen hides current form and shows the statrup one
func (taf *TeaboxArgsForm) ShowIntroScreen() {
	taf.layers.SetCurrentPanel("_intro-screen")
}

func (taf *TeaboxArgsForm) onError(err error) {
	if err != nil {
		teabox.GetTeaboxApp().Stop()
		taf.GetLogger().Error(err.Error())
	}
}

func (taf *TeaboxArgsForm) generateForms(c teaboxlib.TeaConfComponent) {
	if c.IsGroupContainer() {
		for _, x := range c.GetChildren() {
			taf.generateForms(x)
		}
	}

	if len(c.GetCommands()) < 1 {
		return
	}

	mod := c.(*teaboxlib.TeaConfModule) // Only module can have at least command
	formPanel := NewTeaFormsPanel()

	for _, cmd := range mod.GetCommands() { // One form can have many tabs!
		f := formPanel.AddForm(mod.GetTitle(), cmd.GetTitle())

		// Update relative path to its absolute
		if !strings.HasPrefix(cmd.GetCommandPath(), "/") {
			cmd.SetCommandPath(path.Join(mod.GetModulePath(), cmd.GetCommandPath()))
		}

		taf.modCmdIndex[f.GetId()] = cmd

		// Add static flags
		for _, flag := range cmd.GetStaticFlags() {
			taf.AddFlag(f.GetId(), flag)
		}

		// Add arguments
		for _, a := range cmd.GetArguments() {
			switch a.GetWidgetType() {
			case "dropdown", "list":
				taf.onError(taf.addDropdownListWidget(mod.GetTitle(), cmd.GetTitle(), f, a))
			case "text":
				taf.onError(taf.addTextWidget(mod.GetTitle(), cmd.GetTitle(), f, a))
			case "toggle":
				taf.onError(taf.addToggleWidget(mod.GetTitle(), cmd.GetTitle(), f, a))
			case "silent":
			}
		}

		// Or next/previous, if not the last form
		f.AddButton("Start", func() {
			formPanel.ShowStdoutWindow()
			go func() {
				if err := formPanel.wout.Action(mod.GetCallbackPath(), taf.modCmdIndex[f.GetId()].GetCommandPath(), taf.GetCommandArguments(f.GetId())...); err != nil {
					teabox.GetTeaboxApp().Stop()
					fmt.Println("Error:", err)
				}
				formPanel.SetCurrentPanel(f.GetId())
			}()
		})

		f.AddButton("Cancel", func() {
			teabox.GetTeaboxApp().SetFocus(GetTeaboxMainWindow().GetMainMenu().GetWidget())
			taf.ShowIntroScreen()
		})

		break // currently we take only a first command
	}

	taf.layers.AddPanel(mod.GetTitle(), formPanel, true, false)
}

func (taf *TeaboxArgsForm) addDropdownListWidget(modName, cmdName string, tf *TeaForm, arg *teaboxlib.TeaConfModArg) error {
	opts := []string{}
	for _, opt := range arg.GetOptions() {
		if v, _ := opt.GetValue().(string); v != "" {
			opts = append(opts, v)
		}
	}
	if len(opts) == 0 {
		return fmt.Errorf("List \"%s\" in command \"%s\" of module \"%s\" has no values.", arg.GetWidgetLabel(), cmdName, modName)
	}

	tf.AddDropDownSimple(arg.GetWidgetLabel(), 0, func(index int, option *crtview.DropDownOption) {
		taf.AddArgument(tf.GetId(), arg.GetArgName(), strings.TrimSpace(option.GetText()))
	}, opts...)
	return nil
}

/*
Text could have only one argument as a default text:

	[DEFAULT_TEXT]

The field can be also completely empty.
*/
func (taf *TeaboxArgsForm) addTextWidget(modName, cmdName string, tf *TeaForm, arg *teaboxlib.TeaConfModArg) error {
	if len(arg.GetOptions()) > 0 {
		val := arg.GetOptions()[0].GetValueAsString()
		if val != "" {
			// register the value, if any
			taf.AddArgument(tf.GetId(), arg.GetArgName(), val)
		}

		tf.AddInputField(arg.GetWidgetLabel(), val, 0, nil, func(text string) {
			taf.AddArgument(tf.GetId(), arg.GetArgName(), strings.TrimSpace(text))
		})
	}

	return nil
}

func (taf *TeaboxArgsForm) addToggleWidget(modName, cmdName string, tf *TeaForm, arg *teaboxlib.TeaConfModArg) error {
	if len(arg.GetOptions()) == 0 {
		return fmt.Errorf("Toggle \"%s\" in command \"%s\" of module \"%s\" should have its default state with at least one option.", arg.GetWidgetLabel(), cmdName, modName)
	}

	state, ok := arg.GetOptions()[0].GetValue().(bool) // This is not the *value* but checked/unchecked state
	if ok && state {
		// Register state
		taf.AddArgument(tf.GetId(), arg.GetArgName(), arg.GetOptions()[0].GetLabel())
	}

	tf.AddCheckBox(arg.GetWidgetLabel(), "", state, func(checked bool) {
		if checked {
			taf.AddArgument(tf.GetId(), arg.GetArgName(), arg.GetOptions()[0].GetLabel())
		} else {
			taf.RemoveArgument(tf.GetId(), arg.GetArgName())
		}
	})

	return nil
}
