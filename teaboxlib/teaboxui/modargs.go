package teaboxui

import (
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox"
	"gitlab.com/isbm/teabox/teaboxlib"
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
	wout *TeaSTDOUTWindow
	*crtview.Panels
}

func NewTeaFormsPanel() *TeaFormsPanel {
	tfp := &TeaFormsPanel{
		Panels: crtview.NewPanels(),
	}

	tfp.wout = NewTeaSTDOUTWindow()
	tfp.AddPanel("_stdout", tfp.wout, true, false)

	return tfp
}

func (tfp *TeaFormsPanel) ShowStdoutWindow() {
	tfp.SetCurrentPanel("_stdout")
}

func (tfp *TeaFormsPanel) GetStdoutWindow() *TeaSTDOUTWindow {
	return tfp.wout
}

func (tfp *TeaFormsPanel) AddForm(title string) *TeaForm {
	f := NewTeaForm()

	f.SetTitle(title)
	f.SetId(f.GetTitle())

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
	conf *teaboxlib.TeaConf
	TeaboxBaseWindow
	layers *crtview.Panels

	modCmdIndex map[string]*teaboxlib.TeaConfModCommand
}

func NewTeaboxArgsForm(conf *teaboxlib.TeaConf) *TeaboxArgsForm {
	taf := new(TeaboxArgsForm)
	taf.conf = conf
	taf.modCmdIndex = map[string]*teaboxlib.TeaConfModCommand{}

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

	for _, mod := range taf.conf.GetModuleStructure() {
		taf.generateForms(mod)
	}

	return taf
}

// ShowIntroScreen hides current form and shows the statrup one
func (taf *TeaboxArgsForm) ShowIntroScreen() {
	taf.layers.SetCurrentPanel("_intro-screen")
}

func (taf *TeaboxArgsForm) generateForms(c teaboxlib.TeaConfComponent) {
	if c.GetType() != "module" {
		return
	}

	mod := c.(*teaboxlib.TeaConfModule)
	if mod.IsGroupContainer() {
		for _, x := range mod.GetChildren() {
			taf.generateForms(x)
		}
	}

	if len(mod.GetCommands()) < 1 {
		return
	}

	formPanel := NewTeaFormsPanel()

	for _, cmd := range mod.GetCommands() { // One form can have many tabs!
		f := formPanel.AddForm(mod.GetTitle(), cmd.GetTitle())

		// Update relative path to its absolute
		if !strings.HasPrefix(cmd.GetCommandPath(), "/") {
			cmd.SetCommandPath(path.Join(mod.GetModulePath(), cmd.GetCommandPath()))
		}

		taf.modCmdIndex[f.GetId()] = cmd

		for _, a := range cmd.GetArguments() {
			switch a.GetWidgetType() {
			case "dropdown", "list":
				opts := []string{}
				for _, opt := range a.GetOptions() {
					if v, _ := opt.GetValue().(string); v != "" {
						opts = append(opts, v)
					}
				}
				f.AddDropDownSimple(a.GetWidgetLabel(), 0, nil, opts...)
			case "text":
				// Text should have at least one argument, and one optional:
				// <NAME>            to what option to bind its value, e.g. "--name"
				// [DEFAULT_TEXT]    A default text, but can be omitted, e.g. "John Smith"
				if len(a.GetOptions()) > 0 {
					f.AddInputField(a.GetWidgetLabel(), a.GetOptions()[0].GetLabel(), 80, nil, nil) // GetLabel() contains default text, at least for now
				}
			case "toggle":
				f.AddCheckBox(a.GetWidgetLabel(), "", false, nil)
			case "silent":
			}
		}

		// Or next/previous, if not the last form
		f.AddButton("Start", func() {
			formPanel.ShowStdoutWindow()
			go func() {
				cmd := exec.Command(taf.modCmdIndex[f.GetId()].GetCommandPath(), "/opt/bin/blah")
				cmd.Stdout = formPanel.GetStdoutWindow().GetWindow()
				cmd.Stderr = formPanel.GetStdoutWindow().GetWindow()
				if err := cmd.Run(); err != nil {
					teabox.GetTeaboxApp().Stop()
					fmt.Println("Error:", err)
				}
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
