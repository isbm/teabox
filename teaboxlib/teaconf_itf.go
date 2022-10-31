package teaboxlib

import (
	"sort"
)

// Interface
type TeaConfComponent interface {
	SetGroup(id string)
	GetGroup() string
	GetType() string
	GetTitle() string
	SetTitle(title string)
	GetId() string
	GetCommandName() string
	GetCommands() []*TeaConfModCommand
	Len() int
	Add(mod TeaConfComponent) *TeaConfBaseEntity
	GetChildren() []TeaConfComponent
	IsGroupContainer() bool
	SetModulePath(string)
	GetModulePath() string
}

// Entity
type TeaConfBaseEntity struct {
	title      string
	etype      string
	id         string // command or group or module name
	children   []TeaConfComponent
	modulePath string
}

func (tcb *TeaConfBaseEntity) getChildrenContainer() []TeaConfComponent {
	if tcb.children == nil {
		tcb.children = []TeaConfComponent{}
	}
	return tcb.children
}

func (tcb *TeaConfBaseEntity) SetModulePath(p string) {
	tcb.modulePath = p
}

func (tcb *TeaConfBaseEntity) GetModulePath() string {
	return tcb.modulePath
}

func (tcb *TeaConfBaseEntity) Add(mod TeaConfComponent) *TeaConfBaseEntity {
	tcb.children = append(tcb.getChildrenContainer(), mod)
	sort.Slice(tcb.children, func(i, j int) bool {
		return tcb.children[i].GetTitle() < tcb.children[j].GetTitle()
	})
	return tcb
}

func (tcb *TeaConfBaseEntity) GetChildren() []TeaConfComponent {
	return tcb.children
}

func (tcb *TeaConfBaseEntity) IsGroupContainer() bool {
	return tcb.children != nil && len(tcb.children) > 0
}

func (tcb *TeaConfBaseEntity) Len() int {
	return len(tcb.getChildrenContainer())
}

func (tcb *TeaConfBaseEntity) GetTitle() string {
	return tcb.title
}

func (tcb *TeaConfBaseEntity) SetTitle(title string) {
	tcb.title = title
}

func (tcb *TeaConfBaseEntity) GetId() string {
	return tcb.id
}

func (tcb *TeaConfBaseEntity) GetType() string {
	return tcb.etype
}

func (tcb *TeaConfBaseEntity) GetGroup() string {
	return tcb.id
}

func (tcb *TeaConfBaseEntity) SetGroup(id string) {
	tcb.id = id
}

func (tcb *TeaConfBaseEntity) GetCommandName() string {
	// This cannot be a group, so the ID in this case is a command
	return tcb.id
}

// A single module in the menu root, opens immediately a form to run

// Group of the modules by a topic
type TeaConfGroup struct {
	TeaConfBaseEntity
}

func NewTeaConfGroup(id string) *TeaConfGroup {
	tcg := new(TeaConfGroup)
	tcg.id = id
	tcg.etype = "group"
	tcg.SetTitle(id) // group is also a title in this case
	return tcg
}

func (tcg *TeaConfGroup) GetCommands() []*TeaConfModCommand {
	return nil
}

// Command
type TeaConfCmd struct {
	TeaConfBaseEntity
}

func NewTeaConfCmd(id, title string) *TeaConfCmd {
	tc := new(TeaConfCmd)
	tc.SetTitle(title)
	tc.id = id
	tc.etype = "command"
	return tc
}

func (tc *TeaConfCmd) GetCommands() []*TeaConfModCommand {
	return nil
}
