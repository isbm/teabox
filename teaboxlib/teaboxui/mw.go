package teaboxui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/isbm/crtview"
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

type TeaboxMainWindow struct {
	appRef *crtview.Application
	title  string

	// Windows
	menu *TeaboxMenu

	p          *TeaboxWorkspacePanels
	conf       *teaboxlib.TeaConf
	formWindow *TeaboxArgsForm
}

func NewTeaboxMainWindow(app *crtview.Application, conf *teaboxlib.TeaConf) *TeaboxMainWindow {
	tmw := new(TeaboxMainWindow)
	tmw.appRef = app
	tmw.appRef.EnableMouse(true)
	tmw.appRef.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			tmw.appRef.SetFocus(tmw.formWindow.GetWidget())
		case tcell.KeyBacktab:
			tmw.appRef.SetFocus(tmw.menu.items)
		default:
			//fmt.Println(event.Key())
		}
		return event
	})
	tmw.SetConfig(conf)
	tmw.title = tmw.conf.GetTitle()

	// Whole workspace
	tmw.p = NewTeaboxWorkspacePanels(tmw.title)

	tmw.menu = NewTeaboxMenu(tmw.appRef, tmw.conf)
	tmw.menu.SetOnSelectedFunc(func(i int, li *crtview.ListItem) {
		tmw.formWindow.ShowForm(strings.TrimSpace(li.GetMainText()))
	})
	tmw.p.GetContainer().AddItem(tmw.menu.GetWidget(), teaboxlib.MAIN_MENU_WIDTH, 1, true)

	tmw.formWindow = NewTeaboxArgsForm(tmw.conf)
	tmw.formWindow.Init()

	tmw.p.GetContainer().AddItem(tmw.formWindow.GetWidget(), 0, 1, false)

	return tmw
}

// SetConfig of the teabox content
func (tmw *TeaboxMainWindow) SetConfig(conf *teaboxlib.TeaConf) *TeaboxMainWindow {
	tmw.conf = conf
	return tmw
}

func (tmw *TeaboxMainWindow) GetContent() crtview.Primitive {
	return tmw.p
}

/*
.SetBackgroundColor(tcell.ColorBlue)
		box.SetBorderColor(tcell.ColorWhite)
		box.SetBorder(true).SetTitle("Hello, world!")
		box.SetText("Output:\n").SetChangedFunc(func() { app.Draw() })
		app.SetRoot(box, true)

*/
