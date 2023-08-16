package syscall

import (
	"log"
	"os"
	"strconv"

	"github.com/iansmith/parigot/api/shared/id"
)

var _hostId id.HostId

func CurrentHostId() id.HostId {
	if _hostId.IsEmptyValue() {
		high, err := strconv.ParseUint(os.Getenv("HOSTID_HIGH"), 16, 64)
		if err != nil {
			log.Printf("failed trying to parse high of host id: %v", err)
			return id.HostIdZeroValue()
		}
		low, err := strconv.ParseUint(os.Getenv("HOSTID_LOW"), 16, 64)
		if err != nil {
			log.Printf("failed trying to parse low of host id: %v", err)
			return id.HostIdZeroValue()
		}
		_hostId = id.HostIdFromPair(uint64(high), uint64(low))
	}
	return _hostId
}
