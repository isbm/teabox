package teabox

import (
	"github.com/isbm/crtview"
)

var _appRef *TeaboxApplication

type TeaboxApplication struct {
	*crtview.Application
}

func GetTeaboxApp() *TeaboxApplication {
	if _appRef == nil {
		_appRef = &TeaboxApplication{
			Application: crtview.NewApplication(),
		}
	}
	return _appRef
}
