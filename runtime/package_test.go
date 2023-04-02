package runtime

import (
	"fmt"
	"net/url"
)

func ExamplePackagePath() {
	u, err := url.Parse(PkgUrl)

	fmt.Printf("test: PackageUrl() -> %v %v %v", PkgUrl, err, u.String())

	//Output:
	//fail
}
