package teaboxlib

// # Loader
// --------
//
// Set progress percentage on the progress bar.
//
//  init-set-progress:int:42

var INIT_SET_PROGRESS string = "init-set-progress"

// Increment progress percentage by a previously allocated step.
// For example, if there was allocated 3 steps (each 33%), then
// incrementing by one will move progress bar by 33%. Example:
//
//   init-inc-progress::

var INIT_INC_PROGRESS string = "init-inc-progress"

// Allocate progress increment steps. Basically, 100 divided by
// the number of these steps. Example allocating a step by 33%:
//
//   init-alloc-progress:int:3

var INIT_ALLOC_PROGRESS string = "init-alloc-progress"

// Set status in the init form as a string value. Example:
//
//   init-set-status::Darth Wader

var INIT_SET_STATUS string = "init-set-status"

// Reset the entire init form, flushing all the values to the initial.
// Example usage:
//
//   init-reset::

var INIT_RESET string = "init-reset"

// # Logger (lander)
// -----------------
//
// Set status of the STDOUT dumper window. Example usage:
//
//	logger-status::"text to set"

var LOGGER_STATUS string = "logger-status"

// Set title of the STDOUT dumper window. Example usage:
//
//	logger-title::"Output of the Apt Package Manager"

var LOGGER_TITLE string = "logger-title"

// # Progress (lander)
// -------------------
//
// Set event status of an event. Example usage:
//
//   progress-event::"Godzilla happens"

var PROGRESS_EVENT string = "progress-event"

// Allocate progress steps to increment till full state. Example usage:
//
//   progress-allocate:int:5

var PROGRESS_ALLOCATE string = "progress-allocate"

// Increment progress by one step. Example usage:
//
//   progress-next::

var PROGRESS_NEXT string = "progress-next"

// Set progress value directly. Example usage:
//
//   progress-set:int:42

var PROGRESS_SET string = "progress-set"

// Set lookup prefix. This is a static string, which will be a trigger to pick it up
// and set as a event message. Example usage:
//
//   progress-lookup-prefix::STATUS>
//
// When a command is called, STDOUT is piped through. Each line which starts with "STATUS>"
// will be taken as a event message and will update eventbar.

var PROGRESS_LOOKUP_PREFIX string = "progress-lookup-prefix"

// Set lookup regex. This is a POSIX (!) regular expression, which acts as same as
// "progress-lookup-prefix", except it is a regex. Example usage:
//
//   progress-lookup-regex::"^STATUS*"

var PROGRESS_LOOKUP_REGEX string = "progress-lookup-regex"

// # Form (any)
// ------------
//
// Set a value of a field, which is selected by its label. The label
// has to be exactly the same as in the configuration. Example usage:
//
//	field-set-by-label::{Shadow Location}/etc/shadow
var FORM_SET_BY_LABEL string = "field-set-by-label"

// Set a value of a field, selected by its order in the YAML description,
// started from 0. If you added few fields, simply find out its index order
// and access it that way. Example usage:
//
//   field-set-by-ord::{0}/etc/shadow

var FORM_SET_BY_ORD string = "field-set-by-ord"

// Add a value to an existing one for a field, finding the field by its label
// or its order in the YAML description of the module. Adding value on a text
// field will merge the data. Adding value on a list widget will create another
// item. Adding value on a checkbox will override the existing value.
//
// Example usage is identical as setting the value.
//

var FORM_ADD_BY_LABEL string = "field-add-by-label"
var FORM_ADD_BY_ORD string = "field-add-by-ord"

// Clear value of a field by label. Example usage:
//
//   field-reset-by-label::{Label Of The Field}

var FORM_CLR_BY_LABEL string = "field-reset-by-label"

// Clear value of a field by order index. Example usage:
//
//	field-reset-by-label::{0}
var FORM_CLR_BY_ORD string = "field-reset-by-ord"
