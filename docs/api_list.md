# API List

List of API calls from your script to the Teabox during the runtime of your module.

## Loader (init) window

### `init.progress`
Set progress percentage on the progress bar.

    init.progress:int:42

### `init.progress.alloc`

Allocate progress increment steps. Basically, 100 divided by
the number of these steps. Example allocating a step by 33%:

    init.progress.alloc:int:3

### `init.progress.next`

Increment progress percentage by a previously allocated step.
For example, if there was allocated 3 steps (each 33%), then
incrementing by one will move progress bar by 33%. Example:

    init.progress.inc::

### `init.status`

Set status in the init form as a string value. Example:

    init.status::"Darth Wader is happening. Run."

### `init.reset`

Reset the entire init form, flushing all the values to the initial. Example usage:

    init.reset::


## "Logger" landing window

### `logger.status`

Set status of the STDOUT dumper window. Example usage:

    logger.status::"text to set"

### `logger.title`

Set title of the STDOUT dumper window. Example usage:

    logger.title::"Output of the Apt Package Manager"


## "Common" landing window

### `common.progress.event`
Set event status of an event (the label above the progress bar). Example usage:

    common.progress.event::"Godzilla came to your backyard. Smile."

### `common.progress.allocate`

Allocate progress steps to increment till full state. Example usage:

    common.progress.allocate:int:5

### `common.progress.next`

Increment progress by one step. Example usage:

    common.progress.next::

### `common.progress.set`

Set progress value directly. Example usage:

    common.progress.set:int:42

### `common.lookup.prefix`

Set lookup prefix. This is a static string, which will be a trigger to pick it up
and set as a event message. Example usage:

    common.progress.lookup.prefix::STATUS>

When a command is called, STDOUT is piped through. In this example, each line which
starts with "`STATUS>`" will be taken as a event message and will update eventbar.


### `common.lookup.glob`

Set lookup globbing regex (Unix). This is a Unix regular expression, which acts as same as
shell filename matching (the rules applies exactly the same). Example usage:

    common.progress.lookup.glob::"??it*"

The expression above requires to have any character instead of `?` and take everything else
after it by asterisk `*`.

### `common.lookup.regex`

Set lookup regular expression, which acts as same as "`common.lookup.glob`", except it is a regex.
Example usage:

    common.progress.lookup.regex::"^STATUS*"

### `common.list.add`

Add a todo list item in the list. Each list item needs to have an ID, therefore key/value
format should be used. Ordering is stil preserved in the same order as API were called. Example usage:

	common.list.add::{virus}"Transmit virus from computers to sysadmins"

### `common.list.complete`

Complete an item by passing its ID. The item is a checkbox on the common lander page, which will be
marked as checked. In this case ID is required. Example usage:

	common.list.complete::sandwitch

### `common.list.reset`

Reset all the items in the checklist to the "todo" state (i.e. uncheck). Example usage:

	common.list.reset::

### `common.info.add`

Add a text to the existing text in common info area. Previous text will be preserved. Example usage:

	common.info.add::"This is an undocumented feature in Windows."

> 👉🏻 Note, that this will **not** add a whitespace between the previous text.

### `common.info.set`

Set a text to the common info area, removing previous text. Example usage:

	common.info.set::"ATM cell has no roaming feature turned on, notebooks can't connect."


### `common.title`

Set a title. Example usage:

	common.title::"Big To Little Endian Conversion Error"

### `common.reset`

Reset everything on the progress lander page. Example usage:

    common.reset::


## Form

### `field.set.by-label`

Set a value of a field, which is selected by its label. The label
has to be exactly the same as in the configuration. Example usage:

	field.set.by-label::{Shadow Location}/etc/shadow

### `field.set.by-ord`

Set a value of a field, selected by its order in the YAML description,
started from 0. If you added few fields, simply find out its index order
and access it that way. Example usage:

    field.set.by-ord::{0}/etc/shadow

### `field.add.by-label`

Add a value to an existing one for a field, finding the field by its label
or its order in the YAML description of the module. Adding value on a text
field will merge the data. Adding value on a list widget will create another
item. Adding value on a checkbox will override the existing value.

Example usage is identical as setting the value.


### `field.add.by-ord`

This API call works the same as `field.set.by-ord`, except it adds the value
to the field.

### `field.reset.by-label`

Clear value of a field by label. Example usage:

    field.reset.by-label::{Label Of The Field}

### `field.reset.by-ord`

Clear value of a field by order index. Example usage:

    field.reset.by-label::{0}

## Session (State Storage)

Session is a very simple key/value in-memory storage to maintain module state across scripts
and share the information between them, if it is needed. Syntax is the same as key/value
accessing fields.
### `sesstion.set`

Set a value to the session using a key. Example usage:

    session.set::{name}"John Smith"
    session.set:int:{age}42

### `session.get`

Get a value from the session, using a key. Example usage (keys are _always_ strings):

    session.get::name
    session.get::age

### `session.keys`

Get a list of available keys in the session. It returns a list of them. Example:

    session.keys::

### `session.delete`

Delete a particular value by key from the session. Example (keys are _always_ strings):

    session.delete::name

### `session.flush`

Flush the entire session, emptying it. Example:

    session.flush::
