# What  "Teabox" Is?

## TL;DR

This is a "`whiptail` on steroids" that allows you automatically
generate terminal UI for your CLI scripts, commands and launch them.

But wait! This only a very short story...

## A bit longer story

Prior to "Teabox", there was a `whiptail` and a Bash script. Then this
script started enormously grow up and became unmaintainable. Classics,
you know...

This brought to the need breaking this script into modules. However,
hackers do not like to be intimidated to write only in one specific
language. Therefore it raised the need of supporting many of other
languages, or build a system which literally _does not care_ which
language you are using in your module.

If you know a bit about [SUSE YaST2 (Yet Another Setup Tool)](https://yast.opensuse.org),
then you may also know that essentially this is a container of modules,
that provides you a standard way of organising them as being grouped
into a certain order by topic etc, so then these modules can do a specific
job, they are intended to.

Project "Teabox" resembles exactly that.

