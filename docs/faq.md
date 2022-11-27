# FAQ

## Why making anything like that if there is YaST?

Yes, there is. As YaST developers write:

> YaST consists of several pieces that work together despite being distributed
> among several packages and repositories.
>
> [...]
> 
> That **complex system** is translated into files and folders in a consistent way
> all along the full set of YaST repositories and therefore in the installed system,
> as detailed in the YaST code organization description.

YaST is very complex, big, not always portable and has a lot of features one
might not always need. Also YaST requires Ruby runtime and all related efforts with
packaging and maintenance of all the Ruby subsystem. The "Teabox", in contrast, strives
to be as small as possible, one small compiled binary, which require no runtime.

## What language I have to use to write my module?

[Anything](https://esolangs.org/wiki/Brainfuck) that can produce an executable that can
be specified as a target command and can accept parameters, do output and call Unix socket
on your localhost. To ease your life, you can start from a regular shell scripting.

## What Teabox is not covering?

Project "Teabox" is not:

- An attempt to (re)implement YaST
- A setup tool
- A replacement for a similar tools (`whiptail` or `dialog` etc)
- A terminal GUI _library_ for further development