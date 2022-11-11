package teaboxui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox"
	"gitlab.com/isbm/teabox/teaboxlib"
)

const (
	MAIN_WINDOW   = "main"
	SUBMENU_POPUP = "submenu"
)

type TeaboxWorkspacePanels struct {
	container *crtview.Flex
	*crtview.Panels
}

func NewTeaboxWorkspacePanels(title string) *TeaboxWorkspacePanels {
	tbp := &TeaboxWorkspacePanels{
		Panels: crtview.NewPanels(),
	}
	tbp.SetTitle(title)
	tbp.SetBorder(false)
	tbp.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)

	tbp.container = crtview.NewFlex()
	tbp.container.SetDirection(crtview.FlexColumn)

	wspace := crtview.NewFlex()
	wspace.SetDirection(crtview.FlexRow)
	wspace.AddItem(crtview.NewBox(), 1, 1, false)
	wspace.AddItem(tbp.container, 0, 1, true)
	wspace.AddItem(crtview.NewBox(), 1, 1, false)

	tbp.AddPanel("main", wspace, true, true)

	return tbp
}

func (tbp *TeaboxWorkspacePanels) GetContainer() *crtview.Flex {
	return tbp.container
}

// SetBorder overrides default and always disables it anyways. :)
func (tbp *TeaboxWorkspacePanels) SetBorder(b bool) {
	tbp.Panels.SetBorder(false)
	tbp.Box.SetBorder(false)
}

func (tbp *TeaboxWorkspacePanels) Draw(screen tcell.Screen) {
	if !tbp.IsVisible() {
		return
	}

	tbp.Panels.Draw(screen)

	tbp.Lock()
	defer tbp.Unlock()

	style := tcell.StyleDefault
	hdr := style.Background(teaboxlib.WORKSPACE_HEADER).Foreground(teaboxlib.WORKSPACE_HEADER_TEXT)

	// Header/footer
	w, h := screen.Size()
	for i := 0; i < w; i++ {
		screen.SetContent(i, 0, ' ', nil, hdr)
	}

	for i := 0; i < w; i++ {
		screen.SetContent(i, h-1, ' ', nil, hdr)
	}

	// Title
	hdr.Bold(true)
	for i, c := range tbp.GetTitle() {
		screen.SetContent(i+1, 0, c, nil, hdr)
	}
}

var _teaboxMainWindowRef *TeaboxMainWindow

func InitTeaboxMainWindow() *TeaboxMainWindow {
	if _teaboxMainWindowRef == nil {
		_teaboxMainWindowRef = NewTeaboxMainWindow()
	}
	return _teaboxMainWindowRef
}

func GetTeaboxMainWindow() *TeaboxMainWindow {
	return _teaboxMainWindowRef
}

type TeaboxMainWindow struct {
	title string

	// Windows
	menu *TeaboxMenu

	p          *TeaboxWorkspacePanels
	formWindow *TeaboxArgsForm
}

func NewTeaboxMainWindow() *TeaboxMainWindow {
	tmw := new(TeaboxMainWindow)
	teabox.GetTeaboxApp().EnableMouse(true)
	teabox.GetTeaboxApp().SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			teabox.GetTeaboxApp().SetFocus(tmw.formWindow.GetWidget())
		case tcell.KeyBacktab:
			teabox.GetTeaboxApp().SetFocus(tmw.menu.items)
		default:
			//fmt.Println(event.Key())
		}
		return event
	})
	tmw.title = teabox.GetTeaboxApp().GetGlobalConfig().GetTitle()

	// Whole workspace
	tmw.p = NewTeaboxWorkspacePanels(tmw.title)

	tmw.menu = NewTeaboxMenu()
	tmw.menu.SetOnSelectedFunc(func(i int, li *crtview.ListItem) {
		tmw.formWindow.ShowModuleForm(strings.TrimSpace(li.GetMainText()))
	})
	tmw.p.GetContainer().AddItem(tmw.menu.GetWidget(), teaboxlib.MAIN_MENU_WIDTH, 1, true)

	tmw.formWindow = NewTeaboxArgsForm()
	tmw.formWindow.Init()

	tmw.p.GetContainer().AddItem(tmw.formWindow.GetWidget(), 0, 1, false)

	return tmw
}

func (tmw *TeaboxMainWindow) GetContent() crtview.Primitive {
	return tmw.p
}

// GetMainMenu returns the main menu instance.
func (tmw *TeaboxMainWindow) GetMainMenu() *TeaboxMenu {
	return tmw.menu
}
