package lib

import (
	"github.com/iansmith/parigot/api/shared/id"
)

var _hostId = id.NewHostId()

func CurrentHostId() id.HostId {
	return _hostId
}
