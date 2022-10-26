package teaboxui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox/teaboxlib"
)

type TeaboxMenu struct {
	conf   *teaboxlib.TeaConf
	items  *crtview.List
	layers *crtview.Panels

	TeaboxBaseWindow
}

func NewTeaboxMenu(app *crtview.Application, conf *teaboxlib.TeaConf) *TeaboxMenu {
	tm := new(TeaboxMenu)
	tm.conf = conf
	tm.appref = app
	tm.Init()

	return tm
}

func (tm *TeaboxMenu) SetConfig(conf *teaboxlib.TeaConf) *TeaboxMenu {
	tm.conf = conf
	return tm
}

func (tm *TeaboxMenu) ShowSubmenu(id string) {
	tm.layers.SetCurrentPanel(id)
}

func (tm *TeaboxMenu) GetWidget() crtview.Primitive {
	return tm.layers
}

func (tm *TeaboxMenu) makeSubmenu(mod teaboxlib.TeaConfComponent) {
	list := crtview.NewList()
	list.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	list.ShowSecondaryText(false)
	list.SetBorder(true)
	list.SetTitle(mod.GetTitle())
	list.SetTitleAlign(crtview.AlignRight)

	for _, c := range mod.GetChildren() {
		list.AddItem(crtview.NewListItem(fmt.Sprintf("%-"+strconv.Itoa(teaboxlib.MAIN_MENU_WIDTH-2)+"s", c.GetTitle())))
	}

	// Spacer
	list.AddItem(crtview.NewListItem(strings.Repeat(teaboxlib.LABEL_SEP, teaboxlib.MAIN_MENU_WIDTH-2)))
	list.SetItemEnabled(list.GetItemCount()-1, false)

	// Return item
	list.AddItem(crtview.NewListItem(fmt.Sprintf("%-"+strconv.Itoa(teaboxlib.MAIN_MENU_WIDTH-2)+"s", teaboxlib.LABEL_BACK)))

	// Return hook
	list.SetSelectedFunc(func(i int, li *crtview.ListItem) {
		if strings.TrimSpace(li.GetMainText()) == teaboxlib.LABEL_BACK {
			tm.ShowSubmenu("mainmenu")
		}
	})

	tm.layers.AddPanel(mod.GetTitle(), list, true, false)
}

// Init the ui
func (tm *TeaboxMenu) Init() TeaboxWindow {
	tm.layers = crtview.NewPanels()

	tm.items = crtview.NewList()
	tm.items.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	tm.items.ShowSecondaryText(false)
	tm.items.SetBorder(true)
	tm.items.SetTitle("Modules")
	tm.items.SetTitleAlign(crtview.AlignRight)

	// Put a hook on exit
	tm.items.SetSelectedFunc(func(i int, li *crtview.ListItem) {
		ref := li.GetReference().(teaboxlib.TeaConfComponent)
		if ref.GetTitle() == teaboxlib.LABEL_EXIT {
			tm.appref.Stop()
		} else if ref.GetType() == "group" {
			tm.ShowSubmenu(ref.GetTitle())
		} else if ref.GetType() == "module" {
			tm.appref.Stop()
			spew.Dump(ref.(*teaboxlib.TeaConfModule))
		}
	})

	for idx, mod := range tm.conf.GetModuleStructure() {
		if mod.GetType() == "group" {
			tm.makeSubmenu(mod)
		}
		if mod.GetTitle() == teaboxlib.LABEL_EXIT {
			tm.items.AddItem(crtview.NewListItem(strings.Repeat(teaboxlib.LABEL_SEP, teaboxlib.MAIN_MENU_WIDTH-2)))
			tm.items.SetItemEnabled(idx, false)
		}

		item := crtview.NewListItem(fmt.Sprintf("%-"+strconv.Itoa(teaboxlib.MAIN_MENU_WIDTH-2)+"s", mod.GetTitle()))
		item.SetReference(mod)
		tm.items.AddItem(item)
	}

	tm.layers.AddPanel("mainmenu", tm.items, true, true)

	return tm
}
