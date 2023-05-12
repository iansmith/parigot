This package exists only because the functions in this package
need to use BOTH the syscall code (apiwasm/syscall) *and* the lib
code.  If we put this code in either of those, we create an import loop.