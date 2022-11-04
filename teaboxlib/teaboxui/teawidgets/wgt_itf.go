package teawidgets

// TeaboxWindowWidget interface.
// Window widgets are the entire panels, that can use the whole window space.
type TeaboxWindowWidget interface {
	// Action implementation, which takes:
	//   - "callback" path of the Unix socket to start a listener on it
	//   - "cmdpath" is a command absolute path (preferrably) or env returned
	//   - "cmdargs" are the arguments to the "cmdpath"
	Action(callback, cmdpath string, cmdargs ...string) error
}
