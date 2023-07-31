package runtime

import "fmt"

const (
	OriginRegionKey     = "ORIGIN-REGION"
	OriginZoneKey       = "ORIGIN-ZONE"
	OriginSubZoneKey    = "ORIGIN-SUB-ZONE"
	OriginServiceKey    = "ORIGIN-SERVICE"
	OriginInstanceIdKey = "ORIGIN-INSTANCE-ID"
)

var origin Origin

// Origin - struct for origin information
type Origin struct {
	Region     string
	Zone       string
	SubZone    string
	Service    string
	InstanceId string
}

func init() {
	origin.Region = "region"
	origin.Zone = "zone"
}

func SetOrigin(region, zone, subZone, service, instanceId string) {
	origin.Region = region
	origin.Zone = zone
	origin.SubZone = subZone
	origin.Service = service
	origin.InstanceId = instanceId
}

func OriginString() string {
	return fmt.Sprintf("%s:%s:%s:%s:%s", origin.Region, origin.Zone, origin.SubZone, origin.Service, origin.InstanceId)
}

func OriginUrn(nid string, nss, resource string) string {
	return fmt.Sprintf("urn:%v.%v.%v:%v.%v", nid, origin.Region, origin.Zone, nss, resource)
}
