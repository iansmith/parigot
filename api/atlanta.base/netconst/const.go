package netconst

import "hash/crc32"

// if your message doesn't start with this, you have lost sync and should close the connection
// so we can try to reconnect
var MagicStringOfBytes = uint64(0x1789071417760704)
var FrontMatterSize = 8 + 4
var TrailerSize = 4

var KoopmanTable = crc32.MakeTable(crc32.Koopman)
var ReadBufferSize = 4096
