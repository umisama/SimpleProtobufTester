#bin/sh
go build -o=../main ../main.go
../main -r ExampleRoot -u http://localhost:8080 -j ./types.json -p ./types.proto
