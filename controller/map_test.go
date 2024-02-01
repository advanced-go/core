package controller

import (
	"embed"
	"fmt"
	"io/fs"
)

var (
	//go:embed controllertest/*
	f embed.FS
)

const (
	controllersPath = "controllertest/controllers.json"
	//newConfigUri = "file://[cwd]/controllertest/controllers-new.json"
)

func ExampleNewMap() {
	buf, err := fs.ReadFile(f, controllersPath)
	if err != nil {
		fmt.Printf("test: fs.ReadFile() -> [err:%v]\n", err)
		return
	}
	m, status := NewMap(buf)
	fmt.Printf("test: NewMap() -> [ctrls:%v] [status:%v]\n", m != nil, status)

	k := "query"
	c, status0 := m.Get(k)
	fmt.Printf("test: Get(\"%v\") -> [route:%v] [duration:%v] [status:%v]\n", k, c.Name, c.Duration, status0)

	k = "exec"
	c, status0 = m.Get(k)
	fmt.Printf("test: Get(\"%v\") -> [route:%v] [duration:%v] [status:%v]\n", k, c.Name, c.Duration, status0)

	//Output:
	//test: NewMap() -> [ctrls:true] [status:<nil>]
	//test: Get("query") -> [route:query] [duration:2s] [status:<nil>]
	//test: Get("exec") -> [route:exec] [duration:800ms] [status:<nil>]

}
