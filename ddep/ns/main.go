package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/iansmith/parigot/api/plugin/syscall"
	"github.com/iansmith/parigot/api/shared/id"

	"github.com/quic-go/quic-go/http3"
)

var teamName = flag.String("t", "", "sets the team name")

func main() {
	flag.Parse()

	team := strings.TrimSpace(*teamName)
	if team == "" {
		log.Fatalf("unable to start docker containers without a team name (for use as a dns zone)")
	}
	log.Printf("cert file: %s, %s", os.Getenv("HTTP3_CERT_FILE"), os.Getenv("HTTP3_KEY_FILE"))

	mux := http.NewServeMux()
	mux.HandleFunc("/parigot/host", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("request info: %+v, %+v", req.URL, req.Header)
		w.WriteHeader(200)
	})
	err := http3.ListenAndServe(":8000", os.Getenv("HTTP3_CERT_FILE"), os.Getenv("HTTP3_KEY_FILE"), mux)
	log.Fatalf("%T,%v", err, err)

}

func newRawRemote() syscall.SyscallDataRaw {
	return &rawRemote{}
}

type rawRemote struct {
}

// SetService puts a service into SyscallData.  This should only be
// called once for each package_ and name pair. It returns the
// ServiceId for the service named, creating a new one if necessary.
// The client flag should be set to true only when the requesting
// party is a client.  All services should pass false here.  This
// flag effectively means that the requester (package_,name) does not
// to export their service to be ready to run.
// If the bool result is false, then the pair already existed and
// we made no changes to it.
func (n *rawRemote) SetService(ctx context.Context, package_, name string, client bool) (syscall.Service, bool) {

}

// Export finds a service by the given sid and then marks that
// service as being exported. This function returns nil if
// there is no such service.
func (n *rawRemote) Export(ctx context.Context, svc id.ServiceId) syscall.Service {

}

// Import introduces a dendency between the sourge and dest
// services. Thus,  dest must be running before source can run.
// This function returns a kernel error as an int32 in two primary cases.
// 1. one of the src or destination could not be found.  2. The
// newly introduced edge would create a cycle.
func (n *rawRemote) Import(ctx context.Context, src, dest id.ServiceId) int32 {

}

// Launch blocks the caller until all the prerequistes have been
// launched.  It returns a kernel error as an int32 if there is a
// problem.  The most common problem is that the Launch() timed out.
func (n *rawRemote) Launch(context.Context, id.ServiceId) int32 {

}

// PathExists returns true if there is a sequence of dependency
// graph vertices that eventually leads from source to target.
func (n *rawRemote) PathExists(ctx context.Context, source, target string) bool {

}
