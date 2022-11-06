package sys

import (
	"bytes"
	"fmt"
	"log"
)

type depPair struct {
	proc       *Process
	service    string
	linkExport string
}

func (d *depPair) String() string {
	if d.linkExport == "" {
		return fmt.Sprintf("require %s in %s;", d.service, d.proc.String())
	}
	return fmt.Sprintf("require %s (via export of %s) in %s;", d.service, d.linkExport, d.proc.String())
}

// DependencyLoop does a DFS looking for a cycle.
func (n *NameServer) dependencyLoop(cand *depPair, seen []*depPair) string {
	log.Printf("---------------- dependency loop ---------\n")
	log.Printf("candidate: %s\n", cand.String())
	for _, s := range seen {
		log.Printf("\tseen %s\n", s.String())
	}
	log.Printf("---------------- END ---------\n")
	found := false
	//check to see if we found it
	for _, pair := range seen {
		if cand.service == pair.service {
			found = true
			break
		}
	}
	if found {
		log.Printf("got it! %s\n", cand.String())
		var buf bytes.Buffer
		for _, pair := range seen {
			buf.WriteString(pair.String())
		}
		// add the loop end
		buf.WriteString(cand.String())
		log.Printf("got it part 2: %s", buf.String())
		return buf.String()
	}
	//we didn't find it so we are now in seen
	seen = append(seen, cand)
	child := n.findCandidateProcessesByExport(cand)
	for _, c := range child {
		loop := n.dependencyLoop(c, seen)
		if loop != "" {
			return loop
		}
	}
	log.Printf("RET EMPTY %s", cand.String())
	return ""
}

// findCandidateProcessesByExport returns the set of processes that export a particular
// service.
func (n *NameServer) findCandidateProcessesByExport(pair *depPair) []*depPair {
	candidateList := []*depPair{}
	for _, node := range n.dependencyGraph {
		for _, exp := range node.export {
			if exp == pair.service {
				for _, r := range node.require {
					candidateList = append(candidateList,
						&depPair{proc: node.proc,
							linkExport: exp,
							service:    r})
				}
				break
			}
		}
	}
	return candidateList
}

func (n *NameServer) GetLoopContent() string {
	if len(n.dependencyGraph) == 0 {
		panic("should not be sending scanning for loop when every process is running!")
	}
	candidateList := []*depPair{}
	// we want to try all combos
	for _, v := range n.dependencyGraph {
		for _, req := range v.require {
			candidateList = append(candidateList, &depPair{v.proc, req, ""})
		}
	}
	for _, candidate := range candidateList {
		loop := n.dependencyLoop(candidate, []*depPair{})
		log.Printf("deploop returned '%s'", loop)
		if loop != "" {
			return loop
		}
	}
	return ""
}

// getDeadNodeContent returns a list of the nodes that cannot possibly be fulfilled.
func (n *NameServer) GetDeadNodeContent() string {
	candidateList := []*depPair{}
	// we want to try all combos
	for _, v := range n.dependencyGraph {
		for _, req := range v.require {
			candidateList = append(candidateList, &depPair{v.proc, req, ""})
		}
	}
	//now build a list of the possible exports of the whole graph
	possibleExport := []string{}
	for _, v := range n.dependencyGraph {
		for _, export := range v.export {
			possibleExport = append(possibleExport, export)
		}
	}
	log.Printf("possible exports: %+v", possibleExport)
	// strip out all the candidates who could be unlocked by the export
	result := []*depPair{}
outer:
	for _, candidate := range candidateList {
		for _, export := range possibleExport {
			if candidate.service == export {
				continue outer
			}
		}
		log.Printf("candidate passed all exports: %+v", candidate)
		result = append(result, candidate)
	}
	var buf bytes.Buffer
	for _, r := range result {
		buf.WriteString(r.String())
	}
	return buf.String()
}

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
	nameserverPrint("REMOVEREQUIRE ", " dead list? %d %+v", len(deadList), deadList)
	for _, req := range e.require {
		found := false
		for _, dead := range deadList {
			if dead == req {
				found = true
				break
			}
		}
		nameserverPrint("REMOVEREQUIRE ", " req %s not found on dead list", req)
		if !found {
			result = append(result, req)
			changed = true
		}
	}
	e.require = result
	return changed
}
