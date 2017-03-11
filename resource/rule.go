package resource

import (
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// Rule is ..
type Rule struct {
	ec2.IpPermission
}

// CheckPort is ..
func (r *Rule) CheckPort(target int) bool {
	from := r.FromPort
	to := r.ToPort
	return (from == nil && to == nil) || (*from <= int64(target) && *to >= int64(target))
}

// CheckTargetIPv4 is ..
func (r *Rule) CheckTargetIPv4(targetIP string) bool {
	ip := net.ParseIP(targetIP)
	for _, network := range r.IpRanges {
		_, net, err := net.ParseCIDR(*network.CidrIp)
		if err != nil {
			fmt.Print(err)
		} else if net.Contains(ip) {
			return true
		}
	}
	return false
}

// CheckTargetIPv6 is ..
func (r *Rule) CheckTargetIPv6(targetIP string) bool {
	ip := net.ParseIP(targetIP)
	for _, network := range r.Ipv6Ranges {
		_, net, err := net.ParseCIDR(*network.CidrIpv6)
		if err != nil {
			fmt.Print(err)
		} else if net.Contains(ip) {
			return true
		}
	}
	return false
}

// CheckProtocol is ..
func (r *Rule) CheckProtocol(protocol string) bool {
	return *r.IpProtocol == "-1" || *r.IpProtocol == protocol
}
