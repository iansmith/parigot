module github.com/iansmith/parigot/example/httpsimple

go 1.21

replace github.com/iansmith/parigot => /workspaces/parigot

require (
	github.com/iansmith/parigot v0.0.0-20230724152722-d9d85542a70c
	github.com/iansmith/parigot/example/helloworld v0.0.0-20230824141946-b86f0489829d
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/fatih/color v1.15.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/tetratelabs/wazero v1.5.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
)

//#	github.com/iansmith/parigot v0.0.0-20230702130819-d67ad66557ef // indirect
