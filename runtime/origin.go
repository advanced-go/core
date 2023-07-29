package runtime

import "fmt"

const (
	OriginRegionKey     = "ORIGIN-REGION"
	OriginZoneKey       = "ORIGIN-ZONE"
	OriginSubZoneKey    = "ORIGIN-SUB-ZONE"
	OriginServiceKey    = "ORIGIN-SERVICE"
	OriginInstanceIdKey = "ORIGIN-INSTANCE-ID"
)

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

func (o *Origin) String2() string {
	return fmt.Sprintf(OriginRegionKey+"%s:%s:%s:%s:%s", o.Region, o.Zone, o.SubZone, o.Service, o.InstanceId)
}

func NewOrigin(s string) *Origin {
	o := new(Origin)
	o.Region = s
	return o
}
