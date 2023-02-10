.. Teabox documentation master file, created by
   sphinx-quickstart on Fri Nov 25 14:06:47 2022.
   You can adapt this file completely to your liking, but it should at least
   contain the root `toctree` directive.

Welcome to Teabox's documentation!
==================================

.. note::
   This documentation covers Teabox |version| — the solution to provide
   modularised CLI programs launcher, allowing to generate terminal UI
   for them. It resembles similar projects like YaST or Anaconda.

   Teabox is in early continuous development version.

.. toctree::
   :maxdepth: 1
   :caption: Contents:

   usage_overview
   configuration
   writing_module_define
   writing_module_preload
   writing_module_lander
   signal_slots
   rebranding
   api_overview
   api_list
   faq


The Concept
-----------

You have a lot of various setup scripts, that are meant to setup or configure
something. One way to keep them all around is to write a lengthy documentation
with the detailed description how to use them. This assumes that a user will
first carefully read the documentation, study the topic, fully learn the purpose
of that, and then use the software.

Another way is to give a user quick possibility to just go "Next/Next/NextFinish",
especially if a particular task is done once or maybe twice per year or even rarely.

The Teabox is doing exactly that: apart of providing a GUI to the variety of those
command line commands, it also sorts them into a meaningful organised system.


Use Case
--------

Create a setup or configuration tooling via variety of different commands, integrating
them into one system, and provide a unified UI.

Contributing
-------------

* `GitHub Repository <https://github.com/isbm/teabox>`__

   Best way to make progress is to open an issue or submit a Pull Request on the GitHub.
