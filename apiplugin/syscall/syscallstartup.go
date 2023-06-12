package main

var _coord *startupCoordinator

// timeoutInSecs is the number of complete seconds of waiting that
// have to pass before we decide that the startup of a service has
// timed out.
const timeoutInSecs = 15

// coordinator controls access to the underlying variable that has
// the singleton of startupCoordinator.
func coordinator() *startupCoordinator {
	if _coord == nil {
		_coord = newSyscallDataImpl()
	}
	return _coord
}
