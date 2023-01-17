package dep

import (
	"bytes"
	"fmt"
	"os"

	logmsg "github.com/iansmith/parigot/g/msg/log/v1"
	"github.com/iansmith/parigot/sys/backdoor"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	edges map[string]*EdgeHolder
}

func NewDepGraph() *DepGraph {
	return &DepGraph{
		edges: make(map[string]*EdgeHolder),
	}
}

// GetEdge does no locking so callers need to figure out the exclusive access.
func (d *DepGraph) GetEdge(key DepKey) (*EdgeHolder, bool) {

	eh, ok := d.edges[key.String()]
	return eh, ok
}

// PutEdge does no locking so callers need to figure out the exclusive access.
func (d *DepGraph) PutEdge(key DepKey, eh *EdgeHolder) {

	d.edges[key.String()] = eh
}

// Len does no locking so callers need to figure out the exclusive access.
func (d *DepGraph) Len() int {
	return len(d.edges)
}

// Del does no locking so callers need to figure out the exclusive access.
func (d *DepGraph) Del(key DepKey) {

	delete(d.edges, key.String())
}

// AllEdge() returns a map that is keyed on the String() rep of the key.
// This function makes a copy of the depgraph so it really shouldn't be
// used to manipulate the depgraph it is called on.  If you want to mutate
// the depgraph, try Walk().
// AllEdge does not lock, so the caller must assure exclusive acess.
func (d *DepGraph) AllEdge() map[string]*EdgeHolder {
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
//
// Walk allows mutation of the edgeholders during the walk, but the function
// does not lock.  The caller is responsible for insuring exclusive access.
func (d *DepGraph) Walk(fn func(key string, e *EdgeHolder) bool) {
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

// GetLoopContent returns a string (for humans) that explains the loop it has found
// in the dependency graph.  The function returns "" if there is no loop and panics
// if you try to call it on a depgraph that has all the processes running (so it could
// not possibly have a loop).  This function does not lock because it is intended
// to be called once you have already decided that the application is hosed and
// thus nobody else should be changing the depgraph.
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

// GetDeadNodeContent returns a list of the nodes that cannot possibly be fulfilled.
// A dead node is one that imports "foo.bar" when nothing in the deployment exports
// "foo.bar".  This function does not lock because it assumes that you are calling
// because your deployment is hosed and stuck, so nobody else would be changing
// the DepGraph anyway.
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
// records.  None of the operations on an edge holder lock, so callers are required
// to insure exclusive access.  This is also true of the DepGraph that is made up
// of these, so probably a single lock at the level *above* DepGraph is the easiest
// way to implement this.
type EdgeHolder struct {
	key     DepKey
	export  []string
	require []string
}

// NewEdgeHolder returns a newEdgeHolder with empty require and export lists.
func NewEdgeHolder(key DepKey) *EdgeHolder {
	return &EdgeHolder{key: key}
}

// AddExport adds an element to the export list of this holder.
func (e *EdgeHolder) AddExport(s string) {

	e.export = append(e.export, s)
}

// AddRequire adds an element to the require list of this holder.
func (e *EdgeHolder) AddRequire(s string) {
	e.require = append(e.require, s)
}

// WalkRequire passes the function fn all the values in the require list and
// checks to if fn returns false. If fn returns false, the walk stops.
// This function does not lock.
func (e *EdgeHolder) WalkRequire(fn func(s string) bool) {
	for _, s := range e.require {
		if !fn(s) {
			return
		}
	}
}

// Export returns the list of exports for this holder and this list may be nil.
// The returned list is a COPY of the actual list, which can only be accessed
// behind a lock.
func (e *EdgeHolder) Export() []string {
	result := make([]string, len(e.export))
	copy(result, e.export)
	return result
}

// RequireLen returnns the current number of items in the require list.
func (e *EdgeHolder) RequireLen() int {
	return len(e.require)
}

// ExportLen returns the number of exports currently in the export list.
func (e *EdgeHolder) ExportLen() int {
	return len(e.export)
}

// Key returns the key that this edge holder represents the edges for.  This key
// could represent a local process or a remote service.
func (e *EdgeHolder) Key() DepKey {
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
		depgraphPrint("RemoveRequire ", "considering %s on %s ?? FOUND=%v", req, e.key.String(), found)
		if found {
			changed = true
		} else {
			result = append(result, req)
			depgraphPrint("RemoveRequire ", " %s not found, so what was the content? %#v", req, e.require)
		}
	}
	depgraphPrint("RemoveRequire  ", "did %s change? CHANGE=%v (new result of function is %#v)", e.key.String(), changed, result)
	e.require = result
	depgraphPrint("RemoveRequire ", "exiting---- remaining_requires=%+v and final result is %v", e.require, changed)
	return changed
}

// RemoveRequireSimpleRemoves the given candidate from the list of requirements,
// and does nothing if the candidate is not in the require list.  It returns true
// if the require was actually removed, false if it took no action.
//
// This is the more dangerous version of RemoveRequire() because it should only
// be used when the caller can be sure that the *other* nodes that have the same
// require have already been processed correctly, typically via RemoveRequire().
// This is used primarily so that *new* nodes introduced that have requires that have
// already been successfully exported can be filtered properly.  There is no reason
// to add a new node with a requirement foo if foo is already exported.
func (e *EdgeHolder) RemoveRequireSimple(candidate string) bool {
	result := []string{}
	changed := false
	for i := 0; i < len(e.require); i++ {
		if candidate == e.require[i] {
			changed = true
			continue
		}
		result = append(result, e.require[i])
	}
	e.require = result
	return changed
}
func depgraphPrint(method, spec string, arg ...interface{}) {
	if depgraphVerbose {
		part1 := fmt.Sprintf("depGraph:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		msg := fmt.Sprintf("%s%s", part1, part2)
		req := &logmsg.LogRequest{
			Message: msg,
			Stamp:   timestamppb.Now(),
			Level:   logmsg.LogLevel_LOG_LEVEL_DEBUG,
		}
		backdoor.Log(req, true, false, false, nil)
	}
}
