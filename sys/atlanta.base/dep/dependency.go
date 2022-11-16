package dep

import (
	"bytes"
	"fmt"
)

var depgraphVerbose = true

type DepKey interface {
	String() string
	IsKey() bool
}

type DepGraph struct {
	edges map[string]*EdgeHolder
}

func NewDepGraph() *DepGraph {
	return &DepGraph{
		edges: make(map[string]*EdgeHolder),
	}
}
func (d *DepGraph) GetEdge(key DepKey) (*EdgeHolder, bool) {
	eh, ok := d.edges[key.String()]
	return eh, ok
}
func (d *DepGraph) PutEdge(key DepKey, eh *EdgeHolder) {
	d.edges[key.String()] = eh
}
func (d *DepGraph) Len() int {
	return len(d.edges)
}
func (d *DepGraph) Del(key DepKey) {
	delete(d.edges, key.String())
}

// Note that the map you get back here from AllEdge is keyed on the STRING()
// rep of the key.
func (d *DepGraph) AllEdge() map[string]*EdgeHolder {
	return d.edges
}

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
			//			log.Printf("\t\t\tfound candidate in seen: %s", cand.service)
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
		//		log.Printf("\t\t\tprinted out the loop found")
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
	//	log.Printf("\t\t\tRET EMPTY %s", cand.String())
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
		for _, export := range v.export {
			possibleExport = append(possibleExport, export)
		}
	}
	//	log.Printf("\t\t\tpossible exports: %+v", possibleExport)
	// strip out all the candidates who could be unlocked by the export
	result := []*depInfo{}
outer:
	for _, candidate := range candidateList {
		for _, export := range possibleExport {
			if candidate.service == export {
				continue outer
			}
		}
		//log.Printf("\t\t\tcandidate passed all exports: %+v, so is dead?", candidate)
		result = append(result, candidate)
	}
	//log.Printf("\t\t\tresult is length %d", len(result))
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
	key     DepKey
	export  []string
	require []string
}

func NewEdgeHolder(key DepKey) *EdgeHolder {
	return &EdgeHolder{key: key}
}

func (e *EdgeHolder) AddExport(s string) {
	e.export = append(e.export, s)
}

func (e *EdgeHolder) AddRequire(s string) {
	e.require = append(e.require, s)
}
func (e *EdgeHolder) Require() []string {
	return e.require
}
func (e *EdgeHolder) Export() []string {
	return e.export
}
func (e *EdgeHolder) RequireLen() int {
	return len(e.require)
}
func (e *EdgeHolder) ExportLen() int {
	return len(e.export)
}
func (e *EdgeHolder) Key() DepKey {
	return e.key
}

// IsReady returns true if all the services named in the require() requests are
// are now running. In other words, the folks that _exported_ those servires can accept
// the requests from those that want to use them.
func (e *EdgeHolder) IsReady() bool {
	return len(e.require) == 0
}

// RemoveRequired takes in a list of newly "dead" services and removes any of them
// that it finds in this edgeHolder's list of requirements.  This call is used when
// we discover that some service is ready to run, then we take all of *its* exports
// run them through all the processes edgeHolders, to see if any new processes become
// ready.
func (e *EdgeHolder) RemoveRequire(deadList []string) bool {
	result := []string{}
	changed := false
	depgraphPrint("REMOVEREQUIRE considering if node %s is now enabled to run", e.key.String())
	depgraphPrint("REMOVEREQUIRE ", " exports to remove list size? %d values? %+v", len(deadList), deadList)
	for _, req := range e.require {
		found := false
		for _, dead := range deadList {
			if dead == req {
				found = true
				break
			}
		}
		depgraphPrint("REMOVEREQUIRE ", " req %s not found on dead list of %s", req, e.key.String())
		if !found {
			result = append(result, req)
			changed = true
		}
	}
	depgraphPrint("REMOVEREQUIRE", "did %s change? %v", e.key.String(), changed)
	e.require = result
	return changed
}

func depgraphPrint(method, spec string, arg ...interface{}) {
	if depgraphVerbose {
		part1 := fmt.Sprintf("depGraph:%s", method)
		part2 := fmt.Sprintf(spec, arg...)
		print(part1, part2, "\n")
	}
}
