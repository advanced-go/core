package controller

import (
	"fmt"
)

func ExampleTable_Add_Exists_LookupByName() {
	name := "test-route"
	t := newTable(true, false)
	fmt.Printf("test: empty() -> [%v]\n", t.isEmpty())

	err := t.AddController(newRoute("", nil, nil, nil, nil))
	fmt.Printf("test: Add(nil) -> [err:%v] [count:%v] [exists:%v] [lookup:%v]\n", err, t.count(), t.exists(name), t.LookupByName(name))

	err = t.AddController(newRoute(name, nil, nil, nil, nil))
	fmt.Printf("test: Add(controller) -> [err:%v] [count:%v] [exists:%v] [lookup:%v]\n", err, t.count(), t.exists(name), t.LookupByName(name) != nil)

	t.remove("")
	fmt.Printf("test: remove(\"\") -> [count:%v] [exists:%v] [lookup:%v]\n", t.count(), t.exists(name), t.LookupByName(name) != nil)

	t.remove(name)
	fmt.Printf("test: remove(name) -> [count:%v] [exists:%v] [lookup:%v]\n", t.count(), t.exists(name), t.LookupByName(name))

	//Output:
	//test: empty() -> [true]
	//test: Add(nil) -> [err:[invalid argument: route name is empty]] [count:0] [exists:false] [lookup:<nil>]
	//test: Add(controller) -> [err:[]] [count:1] [exists:true] [lookup:true]
	//test: remove("") -> [count:1] [exists:true] [lookup:true]
	//test: remove(name) -> [count:0] [exists:false] [lookup:<nil>]

}

func ExampleTable_LookupUri() {
	name := "test-route"
	t := newTable(true, false)
	fmt.Printf("test: empty() -> [%v]\n", t.isEmpty())

	r := t.LookupUri("", "")
	fmt.Printf("test: LookupUri(nil,nil) -> [controller:%v]\n", r.Name())

	uri := "urn:postgres:query.access-log"
	r = t.LookupUri(uri, "")
	fmt.Printf("test: LookupUri(%v) -> [controller:%v]\n", uri, r.Name())

	t.SetUriMatcher(func(uri, method string) (string, bool) {
		return name, true
	},
	)
	ok := t.AddController(newRoute(name, NewTimeoutConfig(true, 503, 100), nil, nil, nil))
	fmt.Printf("test: Add(controller) -> [controller:%v] [count:%v] [exists:%v]\n", ok, t.count(), t.exists(name))

	r = t.LookupUri(uri, "")
	fmt.Printf("test: LookupUri(%v) ->  [controller:%v]\n", uri, r.Name())

	//Output:
	//test: empty() -> [true]
	//test: LookupUri(nil,nil) -> [controller:*]
	//test: LookupUri(urn:postgres:query.access-log) -> [controller:*]
	//test: Add(controller) -> [controller:[]] [count:1] [exists:true]
	//test: LookupUri(urn:postgres:query.access-log) ->  [controller:test-route]

}

func ExampleTable_Lookup_Name_Default() {
	//name := "test-route"
	t := newTable(true, true)
	fmt.Printf("test: empty() -> [%v]\n", t.isEmpty())

	act := t.LookupByName("")
	fmt.Printf("test: Lookup(nil) -> [controller:%v]\n", act != nil)

	act = t.LookupByName("test")
	fmt.Printf("test: Lookup(\"test\") -> [controller:%v]\n", act != nil)

	//Output:
	//test: empty() -> [true]
	//test: Lookup(nil) -> [controller:false]
	//test: Lookup("test") -> [controller:true]

}
