package resource

import (
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// SecurityGroup is ..
type SecurityGroup struct {
	ec2.SecurityGroup
}

// CheckIngress is ...
func (sg *SecurityGroup) CheckIngress(protocol string, sourceIP string, port int) bool {
	return sg.CheckRules(sg.IpPermissions, protocol, sourceIP, port)
}

// CheckEgress is ...
func (sg *SecurityGroup) CheckEgress(protocol string, destinationIP string, port int) bool {
	return sg.CheckRules(sg.IpPermissionsEgress, protocol, destinationIP, port)
}

// CheckRules is ..
func (sg *SecurityGroup) CheckRules(rules []*ec2.IpPermission, protocol string, sourceIP string, port int) bool {
	for _, rule := range rules {
		fmt.Printf("%v", rule)
		if rule.FromPort != nil && rule.ToPort != nil {
			if sg.CheckPortInRange(port, int(*rule.FromPort), int(*rule.ToPort)) &&
				sg.CheckIPRanges(sourceIP, rule.IpRanges) &&
				sg.CheckProtocol() {
				return true
			}
		} else {
			if sg.CheckIPRanges(sourceIP, rule.IpRanges) && sg.CheckProtocol() {
				return true
			}
		}
	}
	return false
}

// CheckPortInRange is ..
func (sg *SecurityGroup) CheckPortInRange(target int, from int, to int) bool {
	if from <= target && to >= target {
		return true
	}
	return false
}

// CheckIPRanges is ..
func (sg *SecurityGroup) CheckIPRanges(targetIP string, networks []*ec2.IpRange) bool {
	ip := net.ParseIP(targetIP)
	for _, network := range networks {
		_, net, err := net.ParseCIDR(*network.CidrIp)
		if err != nil {
			fmt.Print(err)
		} else if net.Contains(ip) {
			return true
		}
	}
	return false
}

// CheckProtocol is ..
func (sg *SecurityGroup) CheckProtocol() bool {
	return true
}
