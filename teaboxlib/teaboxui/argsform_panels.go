package teaboxui

import "github.com/isbm/crtview"

type TeaboxArgsFormPanels struct {
	itemref map[string]crtview.Primitive
	*crtview.Panels
}

func NewTeaboxArgsFormPanels() *TeaboxArgsFormPanels {
	tafp := &TeaboxArgsFormPanels{
		Panels:  crtview.NewPanels(),
		itemref: map[string]crtview.Primitive{},
	}

	return tafp
}

func (tafp *TeaboxArgsFormPanels) AddPanel(name string, item crtview.Primitive, resize bool, visible bool) {
	tafp.itemref[name] = item
	tafp.Panels.AddPanel(name, item, resize, visible)
}

func (tafp *TeaboxArgsFormPanels) GetPanelByName(name string) crtview.Primitive {
	return tafp.itemref[name]
}
