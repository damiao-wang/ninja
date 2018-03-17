package main

import (
	"ninja/base/misc/cmdt"
	"ninja/blog/service"
)

func main() {
	cmdt.SetName("博客")
	service.Init()
	cmdt.Execute()
}
