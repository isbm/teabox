# Configuration

## Overview

A basic concept is that Teabox is calling modules, organised as a tree.
There is a three types of the configuration:

1. General configuration for everything
2. Configuration of the whole module suite (there can be many)
3. Specific module configuration

## General configuration

Teabox is generally configured in three locations:

- `/etc/<appname>.conf`
- `~/.<appname>rc`
- `./<appname>.conf`

Teabox if called `teabox` it will work only in a development mode, i.e. will only validate
modules, their syntax and integrity. In order to call it as a configuration suite,
a symlink needs to be created or the whole binary renamed. In packaged system it is usually
an original `teabox` is shipped and a symlink to it is done, e.g:

    /usr/bin/acme -> /usr/bin/teabox

In this case, one of three configurations should exist: `/etc/acme.conf`, `~/.acmerc` or
 `acme.conf` in the current directory. Typically it is recommended to use `/etc` location.

### Available Options

Below is the explanation of a general configuration.

```yaml
# Path where all modules reside
content: /path/to/the/tree/of/modules

# Communication Unix socket. There should be permissions to write it.
# Usually /tmp is the safest place. The socket is created only when module
# is actually called, and is removed after module finished run.
callback: /tmp/teabox.sock

# Global environment, which will be re-exported with each module call.
env:
  PYTHONPATH: /opt/scary/dungeons
```

This config also contains branding theme (colors) for the Teabox instance. But it is described
in a separate chapter, called "Branding/Theming the Teabox".

## Suite Configuration

The whole tree of the modules is called "suite". It has its "entry" configuration file,
and then each sub-directory with an actual module has its own configuration file. Each of these
files are called `init.conf`.

### Available Options

Currently `init.conf` of the entire suite has only one option available:

```yaml
title: My Setup Of Something
```

## Module Configuration

Module configuration is the most complex one, as it defines its startup UI (paremeters),
which executable to call etc. It is also called `init.conf` but it is described in the section
when creating an actual module.

But yes, it also has `title` option. 😊
