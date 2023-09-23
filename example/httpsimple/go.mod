module github.com/iansmith/parigot/example/httpsimple

go 1.21

//replace github.com/iansmith/parigot => /workspaces/parigot

require (
	github.com/iansmith/parigot v0.0.0-20230724152722-d9d85542a70c
	google.golang.org/protobuf v1.31.0
)

require github.com/tetratelabs/wazero v1.5.0 // indirect

//#	github.com/iansmith/parigot v0.0.0-20230702130819-d67ad66557ef // indirect
