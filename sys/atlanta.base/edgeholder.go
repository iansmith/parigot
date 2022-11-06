package sys

// edgeHolder holds the in and out edges of the dependency graph for a single process.
type edgeHolder struct {
	proc    *Process
	export  []string
	require []string
}

// isReady returns true if all the services named in the require() requests are
// are now running. In other words, the folks that _exported_ those servires can accept
// the requests from those that want to use them.
func (e *edgeHolder) isReady() bool {
	return len(e.require) == 0
}

// removeRequired takes in a list of newly "dead" services and removes any of them
// that it finds in this edgeHolder's list of requirements.  This call is used when
// we discover that some service is ready to run, then we take all of *its* exports
// run them through all the processes edgeHolders, to see if any new processes become
// ready.
func (e *edgeHolder) removeRequire(deadList []string) bool {
	result := []string{}
	changed := false

	for _, req := range e.require {
		found := false
		for _, dead := range deadList {
			if dead == req {
				found = true
				break
			}
		}
		if !found {
			result = append(result, req)
			changed = true
		}
	}
	e.require = result
	return changed
}
