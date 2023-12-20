package runtimetest

func stringDefault(key string) string {
	return "default"
}

func Example_LookupString() {
	l := NewLookup[string, func(string) string](stringDefault)

	l.Resolve("")

}

func listDefault(key string) []string {
	return []string{"value-0", "value-1"}
}

func Example_LookupList() {
	l := NewLookup[[]string, func(string) []string](listDefault)

	l.Resolve("")

}
