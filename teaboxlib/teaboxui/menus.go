package teaboxui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox"
	"gitlab.com/isbm/teabox/teaboxlib"
)

type TeaboxMenu struct {
	lastSubmenu  string
	items        *crtview.List
	submenuItems map[string]*crtview.List
	layers       *crtview.Panels

	onSelectecFunction func(i int, li *crtview.ListItem)
	TeaboxBaseWindow
}

func NewTeaboxMenu() *TeaboxMenu {
	tm := new(TeaboxMenu)
	tm.submenuItems = map[string]*crtview.List{}
	tm.Init()

	return tm
}

func (tm *TeaboxMenu) SetOnSelectedFunc(f func(i int, li *crtview.ListItem)) *TeaboxMenu {
	tm.onSelectecFunction = f
	return tm
}

func (tm *TeaboxMenu) ShowSubmenu(id string) {
	tm.layers.SetCurrentPanel(id)
	if id == "mainmenu" {
		id = ""
	}
	tm.lastSubmenu = id
}

func (tm *TeaboxMenu) FocusCurrentMenu() {
	state := false
	if tm.lastSubmenu != "" {
		if submenu, found := tm.submenuItems[tm.lastSubmenu]; found {
			teabox.GetTeaboxApp().SetFocus(submenu)
			state = true
		}
	}

	if !state {
		teabox.GetTeaboxApp().SetFocus(tm.items)
	}
}

func (tm *TeaboxMenu) GetWidget() crtview.Primitive {
	return tm.layers
}

// Creates a list for the menu
func (tm *TeaboxMenu) createMenu(title string) *crtview.List {
	menuStub := crtview.NewList()
	menuStub.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	menuStub.ShowSecondaryText(false)
	menuStub.SetBorder(true)
	menuStub.SetFocusedBorderStyle(crtview.BorderSingle)
	menuStub.SetBorderColorFocused(teaboxlib.MENU_BORDER_SELECTED)
	menuStub.SetBorderColor(teaboxlib.MENU_BORDER)
	menuStub.SetSelectedBackgroundColor(teaboxlib.MENU_ITEM_SELECTED)
	menuStub.SetSelectedTextColor(teaboxlib.MENU_ITEM)
	menuStub.SetDisabledItemTextColor(teaboxlib.MENU_BORDER)
	menuStub.SetTitle(title)
	menuStub.SetTitleAlign(crtview.AlignRight)

	return menuStub
}

func (tm *TeaboxMenu) makeSubmenu(mod teaboxlib.TeaConfComponent) {
	menu := tm.createMenu(mod.GetTitle())
	for _, c := range mod.GetChildren() {
		menu.AddItem(crtview.NewListItem(fmt.Sprintf("%-"+strconv.Itoa(teaboxlib.MAIN_MENU_WIDTH-2)+"s", c.GetTitle())))
	}

	// Spacer
	menu.AddItem(crtview.NewListItem(strings.Repeat(teaboxlib.LABEL_SEP, teaboxlib.MAIN_MENU_WIDTH-2)))
	menu.SetItemEnabled(menu.GetItemCount()-1, false)

	// Return item
	menu.AddItem(crtview.NewListItem(fmt.Sprintf("%-"+strconv.Itoa(teaboxlib.MAIN_MENU_WIDTH-2)+"s", teaboxlib.LABEL_BACK)))

	// Return hook
	menu.SetSelectedFunc(func(i int, li *crtview.ListItem) {
		if strings.TrimSpace(li.GetMainText()) == teaboxlib.LABEL_BACK {
			tm.ShowSubmenu("mainmenu")
		} else {
			tm.onSelectecFunction(i, li)
		}
	})

	tm.submenuItems[mod.GetTitle()] = menu
	tm.layers.AddPanel(mod.GetTitle(), menu, true, false)
}

// Init the ui
func (tm *TeaboxMenu) Init() TeaboxWindow {
	tm.layers = crtview.NewPanels()
	tm.items = tm.createMenu("Modules")

	// Put a hook on exit
	tm.items.SetSelectedFunc(func(i int, li *crtview.ListItem) {
		ref := li.GetReference().(teaboxlib.TeaConfComponent)
		if ref.GetTitle() == teaboxlib.LABEL_EXIT {
			teabox.GetTeaboxApp().Stop("And remember: have a lot of fun!")
		} else if ref.IsGroupContainer() {
			tm.ShowSubmenu(ref.GetTitle())
		} else if ref.GetType() == "module" {
			tm.onSelectecFunction(i, li)
		}
	})

	for idx, mod := range teabox.GetTeaboxApp().GetGlobalConfig().GetModuleStructure() {
		suff := ""
		if mod.IsGroupContainer() || mod.GetGroup() != "" {
			tm.makeSubmenu(mod)
			suff = teaboxlib.LABEL_MORE
		}
		if mod.GetTitle() == teaboxlib.LABEL_EXIT {
			tm.items.AddItem(crtview.NewListItem(strings.Repeat(teaboxlib.LABEL_SEP, teaboxlib.MAIN_MENU_WIDTH-2)))
			tm.items.SetItemEnabled(idx, false)
			suff = "" // reset suffix for exit
		}

		item := crtview.NewListItem(fmt.Sprintf("%-"+strconv.Itoa(teaboxlib.MAIN_MENU_WIDTH-2)+"s", mod.GetTitle()+suff))
		item.SetReference(mod)
		tm.items.AddItem(item)
	}

	tm.layers.AddPanel("mainmenu", tm.items, true, true)

	return tm
}
