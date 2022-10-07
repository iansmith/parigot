all: proto/gen/greet/greet.pb.go

proto/gen/greet/greet.pb.go: proto/greet/greet.proto
	pushd proto >& /dev/null ; buf generate; popd >& /dev/null

curltest:
	curl \
        --header "Content-Type: application/json" \
        --data '{"name": "Jane"}' \
        http://localhost:8080/greet.Greet/Greet

grpcurltest:
	/bin/bash -c "grpcurl -protoset <(buf build -o -) \
	-plaintext -d '{\"name\": \"Jane\"}' \
	localhost:8080 greet.Greet/Greet"

# runs the server
run:
	go run main.go

# if you open another shell with the server still running you can also do this
runclient:
	go run cmd/client/main.go

