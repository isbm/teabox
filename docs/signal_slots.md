# Writing Dynamic UI

Often is the case that the UI is not just a form to collect data and that's all about it.
In fact, we want sometimes hide/show some widgets, reset their values, validate their contents
and these sort of things.

Teabox offers a simple signalling mechanism. It more works like old-school Web CGI script.

## Module State

Your Bash *(or Python or [whatever](https://esolangs.org/wiki/Brainfuck))* script can have
its state during the session. If you activated a form in Teabox, it will create a session
for you, which essentially is a very simple in-memory key/value storage when Teabox is running.
Typical use-cases for the session are:

- Maintain module state
- Exchange data between separate scripts
- Exchange data between modules

You can refer to the {doc}`api_list` description for details how to access it.

### Data Visibility

Each key is actually named with the prefix of the module name. If the name is "`my_module`",
then a key e.g. "`foo`" in reality will be stored as "`my_module.foo`". So there is no way to
access from the other module, because then request key will be prefixed with that module and so on.

Therefore all keys are always private to a module.

### Public Visibility

However, if there is a need to share data across the modules *(not recommended due to "ad-hoc"
nature of this approach)*, this is still possible to declare a key as public.

To declare a key as public, it needs to be prefixed with a colon "`:`", like so:

    session::set{:name}"John Smith"

In this case key "`:name`" will not be prefixed and in this way can be accessed from anywhere,
e.g. from another module within the same runtime session of Teabox.

## Signal Slots

Signal slots allows to perform an action when an event occurs. There are following events
supported:
- Widget was selected
- Widget was de-selected
- Widget was changed (its value)

Each of these widget states can trigger an action that calls any script within the module 
directory, where the parameters are defined in "`init.conf`" file of the module itself.

Example:

```yaml
args:
  - type: toggle
    name: --
    attributes:
      - view-only
    label: "Show other widget"
    signals:
      selected: target.sh --show-optional-widgets
      deselected: target.sh --hide-optional-widgets
```

The script "`target.sh`" then should implement all the logic and call any of Teabox's {doc}`api_list` to do something, for example hide or show a widget, or fill it with some value etc.
