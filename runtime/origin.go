package runtime

import "fmt"

// Origin - struct for origin information
type Origin struct {
	Region     string
	Zone       string
	SubZone    string
	Service    string
	InstanceId string
}

func (o *Origin) String() string {
	return fmt.Sprintf("%s:%s:%s:%s:%s", o.Region, o.Zone, o.SubZone, o.Service, o.InstanceId)
}

func NewOrigin(s string) *Origin {
	o := new(Origin)
	o.Region = s
	return o
}
