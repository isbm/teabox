package teaboxui

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox/teaboxlib"
)

type TeaboxArgsForm struct {
	conf *teaboxlib.TeaConf
	TeaboxBaseWindow
	layers *crtview.Panels
}

func NewTeaboxArgsForm(conf *teaboxlib.TeaConf) *TeaboxArgsForm {
	taf := new(TeaboxArgsForm)
	taf.conf = conf

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

	taf.layers.AddPanel("__root__", intro, true, true)

	for _, mod := range taf.conf.GetModuleStructure() {
		if len(mod.GetCommands()) < 1 {
			continue
		}

		f := crtview.NewForm()
		f.SetBorder(true)
		f.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
		f.SetFieldTextColor(tcell.ColorWhite)
		f.SetFieldBackgroundColor(teaboxlib.WORKSPACE_HEADER)
		f.SetFieldBackgroundColorFocused(tcell.ColorGreenYellow)
		f.SetFieldTextColorFocused(teaboxlib.WORKSPACE_BACKGROUND)

		f.SetButtonBackgroundColor(teaboxlib.FORM_BUTTON_BACKGROUND)
		f.SetButtonBackgroundColorFocused(teaboxlib.FORM_BUTTON_BACKGROUND_SELECTED)
		f.SetButtonTextColor(teaboxlib.FORM_BUTTON_TEXT)
		f.SetButtonTextColorFocused(teaboxlib.FORM_BUTTON_TEXT_SELECTED)

		for _, cmd := range mod.GetCommands() { // One form can have many tabs!
			f.SetTitle(fmt.Sprintf("%s - %s", mod.GetTitle(), cmd.GetTitle()))
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
					var v string
					if len(a.GetOptions()) > 0 {
						if data, _ := a.GetOptions()[0].GetValue().(string); data != "" {
							v = data
						}
					}
					f.AddInputField(a.GetWidgetLabel(), v, 80, nil, nil)
				case "toggle":
					f.AddCheckBox(a.GetWidgetLabel(), "", false, nil)
				case "silent":
				}
			}
			break // currently we take only a first command
		}

		f.AddButton("Start", nil)
		f.AddButton("Cancel", func() {
			os.Exit(1)
		})

		taf.layers.AddPanel(mod.GetTitle(), f, true, false)
	}

	return taf
}
