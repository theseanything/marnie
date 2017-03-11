package resource

import "github.com/aws/aws-sdk-go/service/ec2"

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
func (sg *SecurityGroup) CheckRules(rules []*ec2.IpPermission, protocol string, targetIP string, port int) bool {
	for _, rule := range rules {
		r := Rule{*rule}
		if r.CheckPort(port) && r.CheckProtocol(protocol) && r.CheckTargetIPv4(targetIP) {
			return true
		}
	}
	return false
}
