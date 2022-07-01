package main

var commands = map[string][]string{
	"all": {
		"go build ./...",
	},
	"install": {
		"go generate ./version/.",
		"go install  .",
	},
	"release": {
		"go generate ./version/.",
		"go install  -ldflags '-w -s' ./.",
	},
	"test": {
		"go test -v ./...",
	},
	"bench": {
		"go bench ./...",
	},
	"buidl": {
		"go generate ./version/.",
		"go install  ./buidl/.",
	},
	"gen": {
		"go generate ./...",
	},
	"tag": {
		"go run ./version/update/. patch",
	},
}
