//go:build tools
// +build tools

package main

//go:generate protoc -I./data --go_out=./data --go_opt=paths=source_relative ./data/payload.proto
