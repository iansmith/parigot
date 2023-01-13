package dep

import (
	"bytes"
	"fmt"
	"os"
	"sync"
)

var depgraphVerbose = false || os.Getenv("PARIGOT_VERBOSE") != ""

// DepKey is the interface that lets local and remote handles to a service
// be used in the same way.  The String() method returns a nicely printed
// version of what it is for humans, Name() is just the name that was configure
// in the deployment file (for sorting).
type DepKey interface {
	String() string
	Name() string
}

// DepGraph keeps the edges needed by each process.  This graph can and does
// change as we resolve dependencies of a given process.  We use a sync.Mutex
// here because the sync.Map data structure is not optimized for this case of
// many changes to existing members of the map.
type DepGraph struct {
	lock  *sync.Mutex
	edges map[string]*EdgeHolder
}

func NewDepGraph() *DepGraph {
	return &DepGraph{
		lock:  &sync.Mutex{},
		edges: make(map[string]*EdgeHolder),
	}
}
func (d *DepGraph) GetEdge(key DepKey) (*EdgeHolder, bool) {
	d.lock.Lock()
	defer d.lock.Unlock()

	eh, ok := d.edges[key.String()]
	return eh, ok
}
func (d *DepGraph) PutEdge(key DepKey, eh *EdgeHolder) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.edges[key.String()] = eh
}
func (d *DepGraph) Len() int {
	d.lock.Lock()
	defer d.lock.Unlock()

	return len(d.edges)
}
func (d *DepGraph) Del(key DepKey) {
	d.lock.Lock()
	defer d.lock.Unlock()

	delete(d.edges, key.String())
}

// AllEdge() returns a map that is keyed on the String() rep of the key.
// This function makes a copy of the depgraph so it really shouldn't be
// used to manipulate the depgraph it is called on.  If you want to mutate
// the depgraph, try Walk().
func (d *DepGraph) AllEdge() map[string]*EdgeHolder {
	d.lock.Lock()
	defer d.lock.Unlock()
	result := make(map[string]*EdgeHolder)
	for key, value := range d.edges {
		result[key] = value
	}
	return result
}

// Walk iterates over all the edge holders and gives you a chance to modify
// the edges present via fn.  The function should not add or delete entire
// edge holders during the walk.  If the function returns false, the Walk will
// stop at that point.
func (d *DepGraph) Walk(fn func(key string, e *EdgeHolder) bool) {
	d.lock.Lock()
	defer d.lock.Unlock()

	for k, e := range d.edges {
		if !fn(k, e) {
			break
		}
	}
}

// depInfo is the implementation of DepKey that knows about the difference between
// a process (local) and a service (remote).
type depInfo struct {
	key        DepKey
	service    string
	linkExport string
}

func (d *depInfo) String() string {
	if d.linkExport == "" {
		return fmt.Sprintf("require %s in %s;", d.service, d.key.String())
	}
	return fmt.Sprintf("require %s (via export of %s) in %s;", d.service, d.linkExport, d.key.String())
}

// DependencyLoop does a DFS looking for a cycle.
func (g *DepGraph) dependencyLoop(cand *depInfo, seen []*depInfo) string {
	found := false
	//check to see if we found it
	for _, pair := range seen {
		if cand.service == pair.service {
			//log.Printf("\t\t\tfound candidate in seen: %s", cand.service)
			found = true
			break
		}
	}
	if found {
		var buf bytes.Buffer
		buf.WriteString("-----dependencyLoop found ---------\n")
		for _, pair := range seen {
			buf.WriteString(pair.String() + "\n")
		}
		// add the loop end
		buf.WriteString(cand.String() + "\n")
		buf.WriteString("-----dependencyLoop ed ---------\n")
		//log.Printf("\t\t\tprinted out the loop found")
		return buf.String()
	}
	//we didn't find it so we are now in seen
	seen = append(seen, cand)
	child := g.findCandidateProcessesByExport(cand)
	for _, c := range child {
		loop := g.dependencyLoop(c, seen)
		if loop != "" {
			return loop
		}
	}
	//log.Printf("\t\t\tRET EMPTY %s", cand.String())
	return ""
}

// findCandidateProcessesByExport returns the set of processes that export a particular
// service.
func (g *DepGraph) findCandidateProcessesByExport(pair *depInfo) []*depInfo {
	candidateList := []*depInfo{}
	for _, node := range g.edges {
		for _, exp := range node.export {
			if exp == pair.service {
				for _, r := range node.require {
					candidateList = append(candidateList,
						&depInfo{
							key:        node.key,
							linkExport: exp,
							service:    r})
				}
				break
			}
		}
	}
	return candidateList
}

func (g *DepGraph) GetLoopContent() string {
	if len(g.edges) == 0 {
		panic("should not be sending scanning for loop when every process is running!")
	}
	candidateList := []*depInfo{}
	// we want to try all combos
	for _, v := range g.edges {
		for _, req := range v.require {
			candidateList = append(candidateList,
				&depInfo{
					key:        v.key,
					service:    req,
					linkExport: ""})
		}
	}
	for _, candidate := range candidateList {
		loop := g.dependencyLoop(candidate, []*depInfo{})
		//		log.Printf("\t\t\tdeploop returned emtpty? %v", loop == "")
		if loop != "" {
			return loop
		}
	}
	return ""
}

// getDeadNodeContent returns a list of the nodes that cannot possibly be fulfilled.
func (g *DepGraph) GetDeadNodeContent() string {
	candidateList := []*depInfo{}
	// we want to try all combos
	for _, v := range g.edges {
		for _, req := range v.require {
			candidateList = append(candidateList,
				&depInfo{
					key:        v.key,
					service:    req,
					linkExport: ""})
		}
	}
	//now build a list of the possible exports of the whole graph
	possibleExport := []string{}
	for _, v := range g.edges {
		possibleExport = append(possibleExport, v.export...)
	}
	// pull out all the candidates who could be unlocked by the export
	result := []*depInfo{}
outer:
	for _, candidate := range candidateList {
		for _, export := range possibleExport {
			if candidate.service == export {
				continue outer
			}
		}
		result = append(result, candidate)
	}
	var buf bytes.Buffer
	for _, r := range result {
		buf.WriteString(r.String())
	}
	return buf.String()
}

// EdgeHolder holds the in and out edges of the dependency graph for a single node.
// The key field is because when you have a value in the graph (this object)
// you may want the REAL key, not the string rep of the key, which is what the map
// records.
type EdgeHolder struct {
	lock    *sync.Mutex
	key     DepKey
	export  []string
	require []string
}

// NewEdgeHolder returns a newEdgeHolder with empty require and export lists.
func NewEdgeHolder(key DepKey) *EdgeHolder {
	return &EdgeHolder{key: key, lock: &sync.Mutex{}}
}

// AddExport adds an element to the export list of this holder.
func (e *EdgeHolder) AddExport(s string) {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.export = append(e.export, s)
}

// AddRequire adds an element to the require list of this holder.
func (e *EdgeHolder) AddRequire(s string) {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.require = append(e.require, s)
}

// Require returns the list of requirements and this list may be nil.
func (e *EdgeHolder) Require() []string {
	e.lock.Lock()
	defer e.lock.Unlock()

	return e.require
}

// Export returns the list of exports for this holder and this list may be nil.
func (e *EdgeHolder) Export() []string {
	e.lock.Lock()
	defer e.lock.Unlock()

	return e.export
}

// RequireLen returnns the current number of items in the require list.
func (e *EdgeHolder) RequireLen() int {
	e.lock.Lock()
	defer e.lock.Unlock()

	return len(e.require)
}

// ExportLen returns the number of exports currently in the export list.
func (e *EdgeHolder) ExportLen() int {
	e.lock.Lock()
	defer e.lock.Unlock()

	return len(e.export)
}

// Key returns the key that this edge holder represents the edges for.  This key
// could represent a local process or a remote service.
func (e *EdgeHolder) Key() DepKey {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.key
}

// IsReady returns true if all the services named in the require() requests are
// are now running. In other words, the folks that _exported_ those servires can accept
// the requests from those that want to use them.
func (e *EdgeHolder) IsReady() bool {
	return len(e.require) == 0
}

// RemoveRequire takes in a list of newly exported services and removes any of them
// that it finds in this edgeHolder's list of requirements.
//
// This call is used when we discover that some service is ready to run, then we
// take all of *its* exports and run them through all the processes edgeHolders,
// to see if any new processes become ready because of this change.
//
// This function will lock its own edges that it is changing, but it does lock the graph.
// That is the responsibility of Walk() on the dependency graph which is how this should
// be called.
func (e *EdgeHolder) RemoveRequire(exportedList []string) bool {
	e.lock.Lock()
	defer e.lock.Unlock()

	result := []string{}
	changed := false
	depgraphPrint("RemoveRequire ", "start--------considering if node %s is now enabled to run", e.key.String())
	depgraphPrint("RemoveRequire ", "exports to remove list size? %d values? %+v", len(exportedList), exportedList)
	for _, req := range e.require {
		found := false
		for _, exported := range exportedList {
			if exported == req {
				found = true
				break
			}
		}
		depgraphPrint("RemoveRequire ", " exportedList is %+v", exportedList)
		depgraphPrint("RemoveRequire ", " req %s FOUND on %s ?? FOUND=%v", req, e.key.String(), found)
		if found {
			changed = true
		} else {
			result = append(result, req)
			depgraphPrint("RemoveRequire ", " %s not found, so what was the content? %#v", req, e.require)
		}
	}
	depgraphPrint("RemoveRequire  ", "did %s change? CHANGE=%v (new result is %#v)", e.key.String(), changed, result)
	e.require = result
	depgraphPrint("RemoveRequire ", "exiting---- e.require=%+v and final result is %v", e.require, changed)
	return changed
}

func depgraphPrint(method, spec string, arg ...interface{}) {
	if depgraphVerbose {
		part1 := fmt.Sprintf("depGraph:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		print(part1, part2, "\n")
	}
}
