package teawidgets

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/isbm/crtview"
	"gitlab.com/isbm/teabox"
	"gitlab.com/isbm/teabox/teaboxlib"
)

/*
landerChecklistItem class is meant to be internal to TeaProgressWindowLander and never used outside of it.
It is implementing a visuals for todo/done list.
*/
type landerChecklist struct {
	stateDone, stateTodo string
	itemsIdx             []string
	itemsState           map[string]bool
	items                map[string][]*crtview.TextView
	drawables            []*crtview.Flex

	*crtview.Flex
}

func meowLanderChecklist() *landerChecklist { // :-P
	lcl := &landerChecklist{
		Flex:      crtview.NewFlex(),
		stateDone: "✔",
		stateTodo: "┄",
	}
	lcl.SetDirection(crtview.FlexRow)
	return lcl.ClearItems()
}

// ClearItems, all of them. Empty the list entirely.
func (lcl *landerChecklist) ClearItems() *landerChecklist {
	lcl.itemsIdx = []string{}
	lcl.itemsState = map[string]bool{}
	lcl.items = map[string][]*crtview.TextView{} // There are always two of them: status (0) and a label (1)
	lcl.drawables = []*crtview.Flex{}
	lcl.Flex.ClearItems()

	return lcl
}

// ResetItems only sets all done items in todo mode.
func (lcl *landerChecklist) ResetItems() *landerChecklist {
	for k := range lcl.items {
		lcl.itemsState[k] = false
	}

	return lcl
}

// AddItem adds a label, accessible by its id.
func (lcl *landerChecklist) AddItem(id, label string) *landerChecklist {
	lcl.itemsIdx = append(lcl.itemsIdx, id)
	lcl.itemsState[id] = false

	// State
	stt := crtview.NewTextView()
	stt.SetText(lcl.stateTodo)
	stt.SetBackgroundColor(lcl.GetBackgroundColor())

	// Label
	lbl := crtview.NewTextView()
	lbl.SetText(label)
	lbl.SetBackgroundColor(lcl.GetBackgroundColor())

	// Register widgets
	lcl.items[id] = []*crtview.TextView{stt, lbl}

	// Wrap into a holder
	lst := crtview.NewTextView()
	lst.SetBackgroundColor(lcl.GetBackgroundColor())
	lst.SetText("[")
	lst.SetTextColor(teaboxlib.WORKSPACE_HEADER)

	rst := crtview.NewTextView()
	rst.SetBackgroundColor(lcl.GetBackgroundColor())
	rst.SetText("]")
	rst.SetTextColor(teaboxlib.WORKSPACE_HEADER)

	i := crtview.NewFlex()
	i.SetDirection(crtview.FlexColumn)
	i.AddItem(lst, 1, 1, false)
	i.AddItem(stt, 1, 1, false)
	i.AddItem(rst, 3, 1, false)
	i.AddItem(lbl, 0, 1, false)

	lcl.drawables = append(lcl.drawables, i)

	// Every time the item is added, the whole layout is re-created, because
	// a filler needs to be added at the end.
	// Remove all items from the container
	lcl.Flex.ClearItems()
	for _, clItem := range lcl.drawables {
		// Add an item
		lcl.Flex.AddItem(clItem, 1, 0, false)
	}

	spacer := crtview.NewTextView()
	spacer.SetBackgroundColor(lcl.GetBackgroundColor())
	lcl.Flex.AddItem(spacer, 0, 1, false)

	return lcl
}

func (lcl *landerChecklist) CompleteItem(id string) *landerChecklist {
	if _, exists := lcl.itemsState[id]; exists {
		lcl.itemsState[id] = true
		lcl.items[id][0].SetText(lcl.stateDone)
	}

	return lcl
}

/*
TeaProgressWindowLander class, implementing a lander page, which shows the overall progress,
and additionally allows to display a list of steps done.
*/

type TeaProgressWindowLander struct {
	steps        int
	stepsOffset  int
	lookupPrefix string
	lookupRegex  string

	checklist   *landerChecklist
	action      func(call *teaboxlib.TeaboxAPICall)
	title       *crtview.TextView    // Title of the lander page
	eventBar    *crtview.TextView    // Like a status bar, but shows a chunk of the progress
	generalInfo *crtview.TextView    // General text information (static per module)
	progressBar *crtview.ProgressBar // Progressbar itself

	*teaCommonBaseWindowLander
	*crtview.Flex
}

func NewTeaProgressWindowLander() *TeaProgressWindowLander {
	return (&TeaProgressWindowLander{
		Flex:        crtview.NewFlex(),
		checklist:   meowLanderChecklist(),
		title:       crtview.NewTextView(),
		eventBar:    crtview.NewTextView(),
		generalInfo: crtview.NewTextView(),
		progressBar: crtview.NewProgressBar(),
	}).init()
}

func (pl *TeaProgressWindowLander) init() *TeaProgressWindowLander {
	// Construct UI here
	pl.SetDirection(crtview.FlexRow)

	pl.title.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	pl.title.SetTitleColor(tcell.ColorWhite.TrueColor())
	pl.title.SetBorderPadding(1, 1, 10, 10)
	pl.AddItem(pl.title, 3, 1, false)

	// Mid spacer
	spacer := crtview.NewTextView()
	spacer.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)

	chlPadder := crtview.NewFlex()
	chlPadder.SetDirection(crtview.FlexColumn)
	chlPadder.AddItem(spacer, 10, 1, false)
	pl.checklist.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	chlPadder.AddItem(pl.checklist, 0, 1, false)
	chlPadder.AddItem(spacer, 10, 1, false)
	pl.AddItem(chlPadder, 0, 1, false)

	pl.generalInfo.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	pl.generalInfo.SetTextColor(teaboxlib.WORKSPACE_HEADER_TEXT)
	pl.generalInfo.SetWordWrap(true)
	pl.generalInfo.SetDynamicColors(true)
	pl.generalInfo.SetBorderPadding(1, 1, 10, 10)
	pl.AddItem(pl.generalInfo, 0, 1, false)

	// Add spacer instead of checklist
	//pl.AddItem(spacer, 0, 1, false)

	// Event bar
	pl.eventBar.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	pl.eventBar.SetTextColor(tcell.ColorWhite.TrueColor())
	pl.eventBar.SetTextAlign(crtview.AlignCenter)
	pl.AddItem(pl.eventBar, 1, 1, false)

	// Progress bar
	pl.progressBar.SetBackgroundColor(teaboxlib.WORKSPACE_BACKGROUND)
	pl.progressBar.SetEmptyColor(teaboxlib.WORKSPACE_HEADER)
	pl.progressBar.SetFilledColor(tcell.ColorYellow.TrueColor())
	pl.progressBar.SetProgress(42)
	pl.progressBar.SetBorder(true)
	pl.progressBar.SetBorderColor(teaboxlib.WORKSPACE_HEADER)
	pl.progressBar.SetBorderPadding(1, 1, 3, 3)

	pgbBorder := crtview.NewFlex()
	pgbBorder.SetBackgroundColor(tcell.ColorRed)
	pgbBorder.AddItem(pl.progressBar, 0, 1, false)

	pgbPadder := crtview.NewFlex()
	pgbPadder.SetDirection(crtview.FlexColumn)
	pgbPadder.AddItem(spacer, 10, 1, false)
	pgbPadder.AddItem(pgbBorder, 0, 1, false)
	pgbPadder.AddItem(spacer, 10, 1, false)

	pl.AddItem(pgbPadder, 6, 1, false)

	// Bottom spacer
	//pl.AddItem(spacer, 4, 1, false)

	// Define API receiver
	pl.action = func(call *teaboxlib.TeaboxAPICall) {
		switch call.GetClass() {
		case teaboxlib.COMMON_PROGRESS_EVENT:
			pl.eventBar.SetText(call.GetString())
		case teaboxlib.COMMON_PROGRESS_NEXT:
			pl.stepsOffset++
			if pl.stepsOffset >= pl.steps {
				pl.stepsOffset = pl.steps
			}
			pl.progressBar.SetProgress((100 / pl.steps) * pl.stepsOffset)
		case teaboxlib.COMMON_PROGRESS_ALLOCATE:
			pl.steps = call.GetInt()
		case teaboxlib.COMMON_PROGRESS_SET:
			pl.progressBar.SetProgress(call.GetInt())

		case teaboxlib.COMMON_LOOKUP_PREFIX:
			pl.lookupPrefix = call.GetString()
		case teaboxlib.COMMON_LOOKUP_REGEX:
			pl.lookupRegex = call.GetString()

		case teaboxlib.COMMON_LIST_ADD_ITEM:
			pl.checklist.AddItem(call.GetKey(), call.GetString())
		case teaboxlib.COMMON_LIST_COMPLETE_ITEM:
			pl.checklist.CompleteItem(call.GetString())
		case teaboxlib.COMMON_LIST_RESET:
			pl.checklist.ResetItems()
		case teaboxlib.COMMON_INFO_ADD:
			pl.generalInfo.SetText(pl.generalInfo.GetText(true) + call.GetString())

		case teaboxlib.COMMON_INFO_SET:
			pl.generalInfo.SetText(call.GetString())

		case teaboxlib.COMMON_TITLE:
			pl.title.SetText(call.GetString())
		}

		teabox.GetTeaboxApp().Draw()
	}

	return pl
}

func (pl *TeaProgressWindowLander) Action(cmdpath string, cmdargs ...string) error {
	// TODO: Watch stdout and update the UI on prefix or regex
	cmd := exec.Command(cmdpath, cmdargs...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(fmt.Sprintf("Error: command \"%s %s\" quit as %s", cmdpath, strings.Join(cmdargs, " "), err.Error()))
	}

	teabox.GetTeaboxApp().Draw()
	return nil
}

func (pl *TeaProgressWindowLander) AsWidgetPrimitive() crtview.Primitive {
	var w TeaboxLandingWindow = pl
	return w.(crtview.Primitive)
}

// Return window receiver action on Unix socket calls, specific per this widget.
// This action is called by Unix socket among others, and it picks up stuff that are needed.
func (pl *TeaProgressWindowLander) GetWindowAction() func(call *teaboxlib.TeaboxAPICall) {
	return pl.action
}

// Reset all the values to the initial state
func (pl *TeaProgressWindowLander) Reset() {
	pl.checklist.ClearItems()
	pl.generalInfo.SetText("")
	pl.progressBar.SetProgress(0)
	pl.steps = 0
	pl.stepsOffset = 0
	pl.lookupPrefix = ""
	pl.lookupRegex = ""

	pl.title.SetText("Welcome!")
	pl.eventBar.SetText("")
}
