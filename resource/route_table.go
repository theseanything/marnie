package resource

import (
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// RouteTable is ..
type RouteTable struct {
	ec2.RouteTable
}

// CheckRoutes is ...
func (t *RouteTable) CheckRoutes(targetIP string) bool {
	for _, r := range t.Routes {
		_, dest, err := net.ParseCIDR(*r.DestinationCidrBlock)
		if err != nil {
			fmt.Println(err)
		}
		if dest.Contains(net.ParseIP(targetIP)) && r.GatewayId != nil {
			return true
		}
	}
	return false
}
