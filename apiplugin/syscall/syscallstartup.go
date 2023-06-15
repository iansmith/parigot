package syscall

var _coord *startupCoordinator = newSyscallDataImpl()

// timeoutInSecs is the number of complete seconds of waiting that
// have to pass before we decide that the startup of a service has
// timed out.
const timeoutInSecs = 1

// coordinator controls access to the underlying variable that has
// the singleton of startupCoordinator.
func coordinator() *startupCoordinator {
	return _coord
}
