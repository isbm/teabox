package teaboxlib

// Interface
type TeaConfComponent interface {
	SetGroup(id string)
	GetGroup() string
	GetType() string
	GetTitle() string
	GetCommandName() string
	GetArgs() []string
	Len() int
	Add(mod TeaConfComponent) *TeaConfBaseEntity
}

// Entity
type TeaConfBaseEntity struct {
	title    string
	etype    string
	id       string // command or group or module name
	cargs    []string
	children []TeaConfComponent
}

func (tcb *TeaConfBaseEntity) getChildrenContainer() []TeaConfComponent {
	if tcb.children == nil {
		tcb.children = []TeaConfComponent{}
	}
	return tcb.children
}

func (tcb *TeaConfBaseEntity) Add(mod TeaConfComponent) *TeaConfBaseEntity {
	tcb.children = append(tcb.getChildrenContainer(), mod)
	return tcb
}

func (tcb *TeaConfBaseEntity) Len() int {
	return len(tcb.getChildrenContainer())
}

func (tcb *TeaConfBaseEntity) GetTitle() string {
	return tcb.title
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

func (tcb *TeaConfBaseEntity) GetArgs() []string {
	if tcb.cargs == nil {
		tcb.cargs = []string{}
	}

	return tcb.cargs
}

// A single module in the menu root, opens immediately a form to run
type TeaConfModule struct {
	TeaConfBaseEntity
}

func NewTeaConfModule(title string) *TeaConfModule {
	tcm := new(TeaConfModule)
	tcm.title = title
	tcm.etype = "module"
	return tcm
}

// Group of the modules by a topic
type TeaConfGroup struct {
	TeaConfBaseEntity
}

func NewTeaConfGroup(id string) *TeaConfGroup {
	tcg := new(TeaConfGroup)
	tcg.id = id
	tcg.etype = "group"
	return tcg
}

// Command
type TeaConfCmd struct {
	TeaConfBaseEntity
}

func NewTeaConfCmd(id, title string) *TeaConfCmd {
	tc := new(TeaConfCmd)
	tc.title = title
	tc.id = id
	tc.etype = "command"
	return tc
}
