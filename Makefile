
.PHONY: example

generate: 
	protoc -I=gensql/proto --go_out=gensql/proto --go_opt=paths=source_relative ./gensql/proto/*.proto

build: generate
	cd protoc-gen-sqlproto && make build

example: build 
	protoc --plugin ./protoc-gen-sqlproto/protoc-gen-sqlproto -I=. --sqlproto_out . ./example/todo/*.proto
