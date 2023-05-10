package apiimpl

import "hash/crc32"

// if your message doesn't start with this, you have lost sync or had some other encoding
// problem.  if you can, it might be a good time to try reconnect to the sender.
var MagicStringOfBytes = uint64(0x1789071417760704)
var FrontMatterSize = 8 + 4
var TrailerSize = 4

var KoopmanTable = crc32.MakeTable(crc32.Koopman)
var ReadBufferSize = 4096

// unix domain socket for talking to the logviewer... note that the SocketEnvVar
// should be "" when you are running an app inside the dev container.  You only
// need SocketEnvVar for things running on the local machine, like the log viewer app.
const SocketEnvVar = "PARIGOT_SOCKET_DIR"
const SocketName = "logviewer.sock"
