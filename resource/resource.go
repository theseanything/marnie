package resource

import "net"

// Resource mimicks aws class
type Resource struct {
	ip net.IPAddr
}

// CheckConnectionTo checks id is in rules
func (r *Resource) CheckConnectionTo(protocol string, ipAddress string, port int) bool {
	return false
}
