package syscall

var _coord *rawLocal = newSyscallDataImpl()

// timeoutInSecs is the number of complete seconds of waiting that
// have to pass before we decide that the startup of a service has
// timed out.
const timeoutInSecs = 5

// startCoordinator controls access to the underlying variable that has
// the singleton of startupCoordinator.
func startCoordinator() *rawLocal {
	return _coord
}
