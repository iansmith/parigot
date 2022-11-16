package main

import (
	"log"
	"time"

	"github.com/iansmith/parigot/g/pb/net"
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/dep"
	"google.golang.org/protobuf/types/known/anypb"
)

// func waitForReady(info *waitInfo) {
// 	log.Printf("blocking client called waitForReady (id is %d)", info.waitId)
// 	t := <-info.ch
// 	log.Printf("blocking client finished waiting in waitForReady (id is %d)", info.waitId)

// 	resp := &net.RunBlockResponse{}
// 	waitLock.Lock()
// 	defer waitLock.Unlock()

// 	if !t {
// 		log.Printf("we've been told that we were timed out (id %d)", info.waitId)
// 		resp.ErrId = lib.MarshalKernelErrId(lib.NewKernelError(lib.KernelNotFound))
// 	} else {
// 		resp.ErrId = lib.MarshalKernelErrId(lib.NoKernelErr())
// 	}
// 	log.Printf("sending req resp to original RunBlock requestor")

// 	a := anypb.Any{}
// 	err := a.MarshalFrom(resp)
// 	if err != nil {
// 		info.respChan <- nil
// 	} else {
// 		info.respChan <- &a
// 	}

// 	// remove from the map before we return and release the lock
// 	_, ok := runBlockWaitingList[info.waitKey.String()]
// 	if !ok {
// 		log.Printf("unable to find the client info waiting for export %s (id %d) ", info.waitKey, info.waitId)
// 	} else {
// 		delete(runBlockWaitingList, info.waitKey.String())
// 	}
// 	// we've sent the client the response and removed them from the map, so we are done
// }

// func notifyWaiter(exportedTypeName string) {
// 	graph := core.DependencyGraph().AllEdge()
// 	candidateList := []dep.DepKey{}
// 	for _, eh := range graph {
// 		//look through the edges
// 		for _, req := range eh.Require() {
// 			if req == exportedTypeName {
// 				candidateList = append(candidateList, eh.Key())
// 				eh.RemoveRequire([]string{req})
// 				break
// 			}
// 		}
// 	}
// 	for _, candidate := range candidateList {
// 		core.RunIfReady(candidate, func(key dep.DepKey) {
// 			wait := runBlockWaitingList[key.String()]
// 			if wait == nil {
// 				log.Printf("NOTIFYWAITER unable to find %s on the waiting list", key)
// 				return // can't do anything here
// 			}
// 			// we need to tell him to hit it
// 			wait.ch <- true
// 		})
// 	}
// }

// timeoutHandler is run periodically by the main goroutine looking for folks on the waiting
// list that have been there too long.
func timeoutHandler() {
	// note that we assert the lock up here because it is not safe to read the map requireWaitingList
	// with others modifying it
	waitLock.Lock()
	defer waitLock.Unlock()
	for s, info := range runBlockWaitingList {
		if time.Now().Sub(info.waitStart) > time.Duration(timeoutClient*time.Second) {
			log.Printf("found a timeout (id=%d), removing from the waiting list and sending response", info.waitId)
			delete(runBlockWaitingList, s)
			sendRunBlockResponse(info, false)
		}
	}
}

// checkForReadyWaiter running a goroutine that is notified by the main goroutine
// to do a check for folks that can run.  Because the state can only change when
// the set of waiters changes, the main goroutine calls this when a new program
// hits the runblockWaitingList.
func checkForReadyWaiter(newCh chan dep.DepKey, notifyCh chan *waitInfo) {
	for {
		key := <-newCh
		x := func() {
			waitLock.Lock()
			defer waitLock.Unlock()

			// run if ready may call the func many times if multiple folks are ready
			core.RunIfReady(key, func(readyKey dep.DepKey) {
				info, ok := runBlockWaitingList[readyKey.String()]
				if !ok {
					log.Printf("WARNING: didn't find %s on the runBlockWaitingList even though NSCore thinks its ready", readyKey.String())
					return
				}
				log.Printf("found that id=%d is ready, so putting him in the ready chan and removing from wait list", info.waitId)
				notifyCh <- info
				delete(runBlockWaitingList, readyKey.String())
			})
		}

		// just so we get the defer and are sure we get an unlock
		x()
	}
}

func sendRunBlockResponse(info *waitInfo, success bool) {
	resp := &net.RunBlockResponse{
		ErrId:    lib.MarshalKernelErrId(lib.NoKernelErr()),
		TimedOut: !success,
	}
	var a anypb.Any
	err := a.MarshalFrom(resp)
	if err != nil {
		log.Printf("WARN: unable to marshal response to %s, internal error", info.waitKey)
		return
	}
	info.respChan <- &a
}
