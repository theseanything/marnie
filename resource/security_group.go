package resource

import (
  "fmt"
  "net"

  "github.com/aws/aws-sdk-go/service/ec2"
)

type SecurityGroup struct{
  ec2.SecurityGroup
}

func (sg *SecurityGroup) CheckIngress(protocol string, source_ip string, port int) bool{
  fmt.Println("Checking ingress")
  if sg.CheckPortInRange(port, sg.IpPermissions.FromPort, sg.IpPermissions.ToPort) &&
     sg.CheckIpRange() &&
     sg.CheckProtocol() {
     return true
   }
  return false
}

func (sg *SecurityGroup) CheckPortInRange(target int, from int, to int) bool {
  if from <= target && to >= target {
    return true
  }
  return false
}

func (sg *SecurityGroup) CheckEgress() {

}

func (sg *SecurityGroup) CheckIpRange() bool {

}
