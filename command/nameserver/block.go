package main

import (
	"log"
	"time"

	"github.com/iansmith/parigot/api/proto/g/pb/net"
	"github.com/iansmith/parigot/lib"
	"github.com/iansmith/parigot/sys/dep"
	"google.golang.org/protobuf/types/known/anypb"
)

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
