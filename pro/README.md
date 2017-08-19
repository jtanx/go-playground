pro
====================

Super simple wrapper that executes a child command. Any signals that are received are forwarded to the child.

This is useful to work around a bug in MSYS2 winpty 0.4.3 that prevents cmd.exe from being launched directly.
