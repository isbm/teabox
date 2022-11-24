package teaboxlib

// # Loader
// --------
//
// Set progress percentage on the progress bar.
//
//  init.progress:int:42

var INIT_SET_PROGRESS string = "init.progress"

// Increment progress percentage by a previously allocated step.
// For example, if there was allocated 3 steps (each 33%), then
// incrementing by one will move progress bar by 33%. Example:
//
//   init.progress.inc::

var INIT_INC_PROGRESS string = "init.progress.next"

// Allocate progress increment steps. Basically, 100 divided by
// the number of these steps. Example allocating a step by 33%:
//
//   init.progress.alloc:int:3

var INIT_ALLOC_PROGRESS string = "init.progress.alloc"

// Set status in the init form as a string value. Example:
//
//   init.status::Darth Wader

var INIT_SET_STATUS string = "init.status"

// Reset the entire init form, flushing all the values to the initial.
// Example usage:
//
//   init.reset::

var INIT_RESET string = "init.reset"

// # Logger (lander)
// -----------------
//
// Set status of the STDOUT dumper window. Example usage:
//
//	logger.status::"text to set"

var LOGGER_STATUS string = "logger.status"

// Set title of the STDOUT dumper window. Example usage:
//
//	logger.title::"Output of the Apt Package Manager"

var LOGGER_TITLE string = "logger.title"

/*
Common lander is used for a "common feedback" and has the following features:
- A cheklist of things that are going to happen
- An area of a text message which accepts dynamic colors (tags)
- A progress bar with a short status

All these features can be shown or hidden and controlled by the API during
the module performance runtime.
*/
// # Common (lander)
// -------------------
//
// Set event status of an event (the label above the progress bar). Example usage:
//
//   common.progress.event::"Godzilla happens"

var COMMON_PROGRESS_EVENT string = "common.progress.event"

// Allocate progress steps to increment till full state. Example usage:
//
//   common.progress.allocate:int:5

var COMMON_PROGRESS_ALLOCATE string = "common.progress.allocate"

// Increment progress by one step. Example usage:
//
//   common.progress.next::

var COMMON_PROGRESS_NEXT string = "common.progress.next"

// Set progress value directly. Example usage:
//
//   common.progress.set:int:42

var COMMON_PROGRESS_SET string = "common.progress.set"

// Set lookup prefix. This is a static string, which will be a trigger to pick it up
// and set as a event message. Example usage:
//
//   common.progress.lookup.prefix::STATUS>
//
// When a command is called, STDOUT is piped through. Each line which starts with "STATUS>"
// will be taken as a event message and will update eventbar.

var COMMON_LOOKUP_PREFIX string = "common.lookup.prefix"

// Set lookup globbing regex (Unix). This is a Unix regular expression, which acts as same as
// shell filename matching (the rules applies exactly the same). Example usage:
//
//   common.progress.lookup.glob::"?ord*"

var COMMON_LOOKUP_GLOB string = "common.lookup.glob"

// Set lookup regex. This is a POSIX (!) regular expression, which acts as same as
// "progress-lookup-prefix", except it is a regex. Example usage:
//
//   common.progress.lookup.regex::"^STATUS*"

var COMMON_LOOKUP_REGEX string = "common.lookup.regex"

// Add a todo list item in the list. Each list item needs to have an ID, therefore key/value
// format should be used. Ordering is stil preserved in the same order as API were called. Example usage:
//
//	 common.list.add::{foo}"Foo should happen"

var COMMON_LIST_ADD_ITEM string = "common.list.add"

// Complete an item by passing its ID. Example usage:
//
//	 common.list.complete::foo

var COMMON_LIST_COMPLETE_ITEM string = "common.list.complete"

// Reset all the list to the "todo" state. Example usage:
//
//	 common.list.reset::

var COMMON_LIST_RESET string = "common.list.reset"

// Add a text to the existing text in common info area. Example usage:
//
//	 common.info.add::"This is an undocumented feature in Windows."

var COMMON_INFO_ADD string = "common.info.add"

// Set a text to the common info area. Example usage:
//
//	 common.info.set::"ATM cell has no roaming feature turned on, notebooks can't connect. Yes, this is called a design limitation!"

var COMMON_INFO_SET string = "common.info.set"

// Set a title. Example usage:
//
//	 common.title::"Big To Little Endian Conversion Error"

var COMMON_TITLE string = "common.title"

// Reset everything on the progress lander page. Example usage:
//
//   common.reset::

var COMMON_RESET string = "common.reset"

/*
Form is the main UI to collect all the parameters in order to launch the module command.

It is launched after the Init screen (above), during which all the form values can be
pre-loaded into form fields from a setup script for their further modification and then
passing back to the action script.

Form cycle essentially works like a simple server-side HTTP form:
- When module starts, the form is covered by the Init screen, which shows onlt a progress
  bar, and during this time the main form is pre-loaded.

- Form is displayed according to the description of the module, and user is changing the
  values on its widgets.

- After form launches the module script, a lander window is displayed, which shows the
  progress of the module process (either logger or common lander)

- A pop-up alert says either the result is successful or failed.
*/
// # Form (any)
// ------------
//
// Set a value of a field, which is selected by its label. The label
// has to be exactly the same as in the configuration. Example usage:
//
//	field.set.by-label::{Shadow Location}/etc/shadow
var FORM_SET_BY_LABEL string = "field.set.by-label"

// Set a value of a field, selected by its order in the YAML description,
// started from 0. If you added few fields, simply find out its index order
// and access it that way. Example usage:
//
//   field.set.by-ord::{0}/etc/shadow

var FORM_SET_BY_ORD string = "field.set.by-ord"

// Add a value to an existing one for a field, finding the field by its label
// or its order in the YAML description of the module. Adding value on a text
// field will merge the data. Adding value on a list widget will create another
// item. Adding value on a checkbox will override the existing value.
//
// Example usage is identical as setting the value.
//

var FORM_ADD_BY_LABEL string = "field.add.by-label"
var FORM_ADD_BY_ORD string = "field.add.by-ord"

// Clear value of a field by label. Example usage:
//
//   field.reset.by-label::{Label Of The Field}

var FORM_CLR_BY_LABEL string = "field.reset.by-label"

// Clear value of a field by order index. Example usage:
//
//	field.reset.by-label::{0}
var FORM_CLR_BY_ORD string = "field.reset.by-ord"
