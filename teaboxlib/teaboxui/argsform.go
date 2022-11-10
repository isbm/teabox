package teaboxui

import (
	"fmt"
	"path"
	"strings"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox"
	"gitlab.com/isbm/teabox/teaboxlib"
	"gitlab.com/isbm/teabox/teaboxlib/teaboxui/teawidgets"
)

// TeaFormsPanel is a layer of windows, and it contains many TeaForm instances to switch between them.
type TeaFormsPanel struct {
	landingPage  teawidgets.TeaboxLandingWindow
	moduleConfig *teaboxlib.TeaConfModule
	objref       map[string]interface{}
	*crtview.Panels
}

func NewTeaFormsPanel(conf *teaboxlib.TeaConfModule) *TeaFormsPanel {
	tfp := &TeaFormsPanel{
		Panels:       crtview.NewPanels(),
		objref:       map[string]interface{}{},
		moduleConfig: conf,
	}

	// Init landing
	switch tfp.moduleConfig.GetLandingPageType() {
	case "logger":
		tfp.landingPage = teawidgets.NewTeaSTDOUTWindow()
		tfp.AddPanel(teawidgets.LANDING_WINDOW_LOGGER, tfp.landingPage.(crtview.Primitive), true, false)
	default:
		teabox.GetTeaboxApp().Stop()
		fmt.Printf("Unfortauntely, type \"%s\" of landing page is not implemented yet\n", tfp.moduleConfig.GetLandingPageType())
	}

	return tfp
}

// GetModuleConfig returns module configuration, as was specified in its init.conf
func (tfp *TeaFormsPanel) GetModuleConfig() *teaboxlib.TeaConfModule {
	return tfp.moduleConfig
}

func (tfp *TeaFormsPanel) GetLandingPage() teawidgets.TeaboxLandingWindow {
	return tfp.landingPage
}

func (tfp *TeaFormsPanel) AddPanel(name string, item crtview.Primitive, resize bool, visible bool) {
	tfp.objref[name] = item
	tfp.Panels.AddPanel(name, item, resize, visible)
}

func (tfp *TeaFormsPanel) StartListener() error {
	// TODO: Add widget update handler action. Currently a noop dummy
	teabox.GetTeaboxApp().GetCallbackServer().AddLocalAction(func(call *teaboxlib.TeaboxAPICall) {})

	// Run the Unix server instance
	if err := teabox.GetTeaboxApp().GetCallbackServer().Start(tfp.moduleConfig.GetCallbackPath()); err != nil {
		teabox.GetTeaboxApp().Stop()
		fmt.Println(err) // That would be a general system problem
	}

	return nil
}

// ShowLandingWindow and start Unix socket server listener with the current callback pack
func (tfp *TeaFormsPanel) ShowLandingWindow() error {
	// Setup local action for the future instance
	teabox.GetTeaboxApp().GetCallbackServer().AddLocalAction(tfp.landingPage.GetWindowAction())
	tfp.SetCurrentPanel(teawidgets.LANDING_WINDOW_LOGGER)

	return nil
}

// StopLandingWindow switches back to the caller form and stops the Unix socket listener
func (tfp *TeaFormsPanel) StopLandingWindow(fid string) error {
	tfp.SetCurrentPanel(fid)
	return tfp.landingPage.StopListener()
}

func (tfp *TeaFormsPanel) GetFormItem(title, subtitle string) interface{} {
	return tfp.objref[fmt.Sprintf("%s - %s", title, subtitle)]
}

func (tfp *TeaFormsPanel) AddForm(title, subtitle string) *teawidgets.TeaboxArgsMainWindow {
	f := teawidgets.NewTeaboxArgsMainWindow(title, subtitle)
	tfp.AddPanel(f.GetId(), f, true, tfp.GetPanelCount() == 1)

	return tfp.GetFormItem(title, subtitle).(*teawidgets.TeaboxArgsMainWindow)
}

// TeaboxArgsForm contains a layers with TeaForms on it, also their output, intro screen, callback screens etc.
type TeaboxArgsForm struct {
	TeaboxBaseWindow

	/*
		Multi-pages forms for all modules of the suite.
		Menu switches between the modules, displaying a default first form, when selected.
	*/
	allModulesForms *TeaboxArgsFormPanels

	/*
		This is not super-nice, but still, per each form (via GetId() method) it holds the following:

		  1. A map of all named arguments with their values, e.g. "--path=/dev/null" etc.
		  2. An array of named arguments without their values in order to put an order to the map above.
		  3. An array of flags, like "-d -v -x -y -z" etc. Currently flags are always passed first.

		Note, these are all arguments to the *module*, so if Teabox is not flexible enough to call a random
		userland command, one can always wrap it with a shell script. :-)

		NOTE: If you still think this is dumb, feel free to make it better and send your PR!
	*/
	modCmdIndex map[string]*teaboxlib.TeaConfModCommand

	wzlib_logger.WzLogger
}

func NewTeaboxArgsForm() *TeaboxArgsForm {
	taf := new(TeaboxArgsForm)
	taf.modCmdIndex = map[string]*teaboxlib.TeaConfModCommand{}

	return taf
}

func (taf *TeaboxArgsForm) GetWidget() crtview.Primitive {
	return taf.allModulesForms
}

func (taf *TeaboxArgsForm) ShowModuleForm(id string) {
	f, ok := taf.allModulesForms.GetPanelByName(id).(*TeaFormsPanel)
	if ok {
		if err := f.StartListener(); err != nil {
			teabox.GetTeaboxApp().GetScreen().Clear()
			taf.GetLogger().Panic(err)
		}

		// TODO: Add a Unix socket hook for preloader and values to the form
		/*
			1. Insert an action that will receive RPC calls to update progress bar and status bar (what's going on during preload)
			   This action will update currently visible loader window.
			   There must be a way to reset it after it is hidden

			2. Insert an action that will receive RPC calls to update currently active but hidden form.
			3. Remove all the actions after the cycle is finished.
		*/

		if f.GetModuleConfig().GetSetupCommand() != "" {
			// Show preload form first
			taf.allModulesForms.SetCurrentPanel(teawidgets.LOAD_WINDOW_COMMON)
			loader := taf.allModulesForms.GetPanelByName(teawidgets.LOAD_WINDOW_COMMON).(*teawidgets.TeaboxArgsLoadingWindow)
			loader.SetAfterLoadAction(func() {
				// TODO: Hook-up Unix receiver with the pre-loader

				// Load finished, so show the main form
				taf.allModulesForms.SetCurrentPanel(id)
				teabox.GetTeaboxApp().Draw()
			})

			// Set receiver hook
			teabox.GetTeaboxApp().GetCallbackServer().AddLocalAction(loader.GetSocketAcceptAction())

			// Loader command could have relative path or absolute.
			// Current directory ("./") is not supported
			var cmd string
			if !strings.HasPrefix(f.GetModuleConfig().GetSetupCommand(), "/") {
				cmd = path.Join(f.GetModuleConfig().GetModulePath(), f.GetModuleConfig().GetSetupCommand())
			} else {
				cmd = f.GetModuleConfig().GetSetupCommand()
			}

			// Call the loader to pre-populate everything
			if err := loader.Load(cmd, f.GetModuleConfig().GetSetupCommandArgs()...); err != nil {
				teabox.GetTeaboxApp().GetScreen().Clear()
				taf.GetLogger().Panic(err)
			}
		} else {
			// No loader specified, show the form "as is" directly
			taf.allModulesForms.SetCurrentPanel(id)
		}
	} else {
		panic(fmt.Sprintf("Panel %s was not found", id))
	}
}

func (taf *TeaboxArgsForm) Init() TeaboxWindow {
	taf.allModulesForms = NewTeaboxArgsFormPanels()

	// Add intro window (common for all)
	taf.allModulesForms.AddPanel(teawidgets.INTRO_WINDOW_COMMON, teawidgets.NewTeaboxArgsIntroWindow(), true, true)

	// Add argloading window (common for all)
	taf.allModulesForms.AddPanel(teawidgets.LOAD_WINDOW_COMMON, teawidgets.NewTeaboxArgsLoadingWindow(), true, false)

	// Build all individual forms
	for _, mod := range teabox.GetTeaboxApp().GetGlobalConfig().GetModuleStructure() {
		taf.generateForms(mod)
	}

	return taf
}

// ShowIntroScreen hides current form and shows the statrup one
func (taf *TeaboxArgsForm) ShowIntroScreen() {
	taf.allModulesForms.SetCurrentPanel(teawidgets.INTRO_WINDOW_COMMON)
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
	// TODO: Define action for updating widgets inside the form
	formPanel := NewTeaFormsPanel(mod)

	for _, cmd := range mod.GetCommands() { // One form can have many tabs!
		f := formPanel.AddForm(mod.GetTitle(), cmd.GetTitle()).
			SetStaticFlags(cmd)

		// Update relative path to its absolute
		if !strings.HasPrefix(cmd.GetCommandPath(), "/") {
			cmd.SetCommandPath(path.Join(mod.GetModulePath(), cmd.GetCommandPath()))
		}
		taf.modCmdIndex[f.GetId()] = cmd

		// Add arguments
		f.AddArgWidgets(cmd)

		// Or next/previous, if not the last form
		f.AddButton("Start", func() {
			// Show resulting end-widget. Those are:
			// - STDOUT "dumb" writer, shows just an output, like a terminal
			// - Checklist done/todo progress screen that has various features, such as progress-bar, status etc (TODO)
			//
			// NOTE: landing window also starts the listener, to which Action() below connects via resulting command Action() calls.
			formPanel.ShowLandingWindow()
			go func() {
				// Run command on the landing window
				if err := formPanel.GetLandingPage().Action(taf.modCmdIndex[f.GetId()].GetCommandPath(), f.GetCommandArguments(f.GetId())...); err != nil {
					teabox.GetTeaboxApp().Stop()
					fmt.Println("Error:", err)
				}

				// Reset landing window to the caller form as done and stop the listener
				formPanel.StopLandingWindow(f.GetId())
			}()
		})

		f.AddButton("Cancel", func() {
			teabox.GetTeaboxApp().SetFocus(GetTeaboxMainWindow().GetMainMenu().GetWidget())
			taf.ShowIntroScreen()
		})

		break // currently we take only a first command
	}

	taf.allModulesForms.AddPanel(mod.GetTitle(), formPanel, true, false)
}
