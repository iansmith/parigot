package interface_

type Pctx interface {
	Log() Log
	Entry(group string, name string) (string, bool)
	SetEntry(group string, name string, value string) bool
	ToBytes() ([]byte, error)
}
