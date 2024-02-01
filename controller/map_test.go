package controller

import (
	"encoding/json"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"os"
)

const (
	configUri    = "file://[cwd]/controllertest/controllers.json"
	newConfigUri = "file://[cwd]/controllertest/controllers-new.json"
)

func ExampleNewMap() {
	buf, err := os.ReadFile(runtime.FileName(configUri))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
		return
	}
	m, status := NewMap(buf)
	fmt.Printf("test: NewMap() -> [ctrls:%v] [status:%v]\n", m != nil, status)

	k := "query"
	c, status0 := m.Get(k)
	fmt.Printf("test: Get(\"%v\") -> [route:%v] [method:%v] [uri:%v] [duration:%v] [status:%v]\n", k, c.Route, c.Method, c.Uri, c.Duration, status0)

	k = "exec"
	c, status0 = m.Get(k)
	fmt.Printf("test: Get(\"%v\") -> [route:%v] [method:%v] [uri:%v] [duration:%v] [status:%v]\n", k, c.Route, c.Method, c.Uri, c.Duration, status0)

	//Output:
	//test: NewMap() -> [ctrls:true] [status:OK]
	//test: Get("query") -> [route:query-route] [method:query] [uri:github/advanced-go/postgresql/pgxsql:query-test-resource.prod] [duration:2s] [status:OK]
	//test: Get("exec") -> [route:exec-route] [method:insert] [uri:github/advanced-go/postgresql/pgxsql:insert.exec-test-resource.dev] [duration:800ms] [status:OK]

}

func _ExampleNewMap_WriteJson() {
	var ctrls []config

	var cfg config
	cfg.Name = "query"
	cfg.Route = "query-route"
	cfg.Method = "query"
	cfg.Duration = "2s"
	cfg.Uri = "github/advanced-go/postgresql/pgxsql:query-test-resource.prod.exec-test-resource.dev"
	ctrls = append(ctrls, cfg)

	cfg.Name = "exec"
	cfg.Route = "exec-route"
	cfg.Method = "exec"
	cfg.Duration = "800ms"
	cfg.Uri = "github/advanced-go/postgresql/pgxsql:insert.exec-test-resource.dev"
	ctrls = append(ctrls, cfg)

	buf, err := json.Marshal(ctrls)
	if err != nil {
		fmt.Printf("test: json.Marshal() -> [err:%v]\n", err)
		return
	}
	err = os.WriteFile(runtime.FileName(newConfigUri), buf, 667)
	if err != nil {
		fmt.Printf("test: os.WriteFile() -> [err:%v]\n", err)
	}

	//Output:

}
