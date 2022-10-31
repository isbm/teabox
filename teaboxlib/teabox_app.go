package teaboxlib

import "github.com/isbm/crtview"

var _appRef *crtview.Application

func GetTeaboxApp() *crtview.Application {
	if _appRef == nil {
		_appRef = crtview.NewApplication()
	}
	return _appRef
}
