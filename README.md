# Teabox
A simple way of making modular scripted terminal GUIs

# What it is?
You know `whiptail` or `dialog` or `xmessage` or `zenity` etc? So this is like these, except it like [YaST](https://yast.opensuse.org).
Essentially, you can make your own YaST, just with Black Jack and... well, you've got the idea.

# In a Nutshell

Essentially, it builds GUI for your module. Your module can be written literally in [_whatever any language_ you like](https://github.com/fabianishere/brainfuck).
Then you describe UI in YAML config, say what conditions needs to be met in order that module to be available, say what parameters you
want to see and voila, you can launch your command.

# It Is More Than That

You can communicate via Unix socket to the UI and pre-load it with values, control widgets from your module etc.


# Limitations (always will be)
No X11. This is terminal-only UI.
