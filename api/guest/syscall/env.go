package syscall

import (
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/iansmith/parigot/api/shared/id"
)

var _hostId id.HostId

// CurrentHostId provides the interface to the runner's chosen host id
// for our WASM machine.  That value is communicated through environment variables.
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

var localTimeZone *time.Location

// CurrentTimezone provides the interface to the configuration file's timezone string
// for our WASM machine.  That value is communicated through environment variables.
func CurrentTimezone() *time.Location {
	if localTimeZone != nil {
		return localTimeZone
	}
	raw := strings.TrimSpace(os.Getenv("TZ"))
	if raw == "local" {
		slog.Error("Warning: local timezone is not supported (probably because running in a WASM sandbox), using UTC", "timezone", raw)
		raw = "UTC"
	}
	var err error
	localTimeZone, err = time.LoadLocation(raw)
	if err == nil {
		return localTimeZone
	}
	if raw == "" || raw == "UTC" {
		//should never happen
		panic("unable to get UTC timezone:" + err.Error())
	}
	log.Printf("%s failed to find timezone %s, using UTC", CurrentHostId().Short(), raw)
	raw = ""
	localTimeZone, err = time.LoadLocation(raw)
	if err != nil {
		//should never happen
		panic("unable to get UTC timezone:" + err.Error())
	}
	return localTimeZone
}
