package host

import "fmt"

func ExampleOriginUrn_Query() {
	urn := OriginUrn("postgresql", "insert", "query-test-resource.prod")

	fmt.Printf("test: OriginUrn() -> %v\n", urn)

	//Output:
	//test: OriginUrn() -> urn:postgresql.region.zone:insert.query-test-resource.prod

}
