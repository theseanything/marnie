package resource

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Instance mimicks aws class
type Instance struct {
	ec2.Instance
}

// CheckSecurityGroups checks id is in rules
func (i *Instance) CheckSecurityGroups(protocol string, ipAddress string, port int) bool {

	var groupIds []*string
	for _, sg := range i.SecurityGroups {
		groupIds = append(groupIds, sg.GroupId)
	}
	sess := session.Must(session.NewSession())
	svc := ec2.New(sess)

	params := &ec2.DescribeSecurityGroupsInput{
		GroupIds: groupIds,
	}
	response, _ := svc.DescribeSecurityGroups(params)

	var securityGroups []*SecurityGroup
	for _, sg := range response.SecurityGroups {
		s := SecurityGroup{*sg}
		if s.CheckIngress(protocol, ipAddress, port) && s.CheckEgress(protocol, ipAddress, port) {
			return true
		}
		securityGroups = append(securityGroups, &s)
	}

	return false
}

// CheckNACLs checks id is in rules
func (i *Instance) CheckNACLs(protocol string, ipAddress string, port int) bool {
	sess := session.Must(session.NewSession())
	svc := ec2.New(sess)

	resp, err := svc.DescribeNetworkAcls(&ec2.DescribeNetworkAclsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("association.subnet-id"),
				Values: []*string{aws.String(*i.SubnetId)},
			},
		},
	})

	if err != nil {
		fmt.Println(err)
	}

	for _, n := range resp.NetworkAcls {
		for _, e := range n.Entries {
			fmt.Println(e)
		}
	}

	return false
}

// CheckRouteTables checks id is in rules
func (i *Instance) CheckRouteTables(protocol string, ipAddress string, port int) bool {
	var publicIPs []*string
	var privateIPs []*string

	for _, n := range i.NetworkInterfaces {
		for _, a := range n.PrivateIpAddresses {
			privateIPs = append(privateIPs, a.PrivateIpAddress)
			if a.Association != nil {
				publicIPs = append(publicIPs, a.Association.PublicIp)
			}
		}
	}

	sess := session.Must(session.NewSession())
	svc := ec2.New(sess)

	resp, err := svc.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("association.subnet-id"),
				Values: []*string{aws.String(*i.SubnetId)},
			},
		},
	})

	if len(resp.RouteTables) == 0 {
		resp, err = svc.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{
					Name:   aws.String("association.main"),
					Values: []*string{aws.String("true")},
				},
			},
		})
	}

	if err != nil {
		fmt.Println(err)
	}

	for _, t := range resp.RouteTables {
		fmt.Println(t)
		routeTable := RouteTable{*t}
		routeTable.CheckRoutes(*i.PublicIpAddress)
	}

	return false
}

// NewInstance ityGroups checks id is in rules
func NewInstance(params *ec2.DescribeInstancesInput) *Instance {
	sess := session.Must(session.NewSession())
	svc := ec2.New(sess)
	response, err := svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances", err.Error())
		log.Fatal(err.Error())
	}
	var ids []*ec2.Instance
	for _, reservation := range response.Reservations {
		for _, instance := range reservation.Instances {
			ids = append(ids, instance)
		}
	}
	numberOfInstances := len(ids)
	if numberOfInstances == 1 {
		return &Instance{*ids[0]}
	} else if numberOfInstances > 1 {
		fmt.Println("Too many instances found.")
		return nil
	}
	fmt.Println("No instance found.")
	return nil
}

// NewInstanceFromNameTag ityGroups checks id is in rules
func NewInstanceFromNameTag(value string) *Instance {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(value),
				},
			},
		},
	}
	return NewInstance(params)
}

// NewInstanceFromId checks id is in rules
func NewInstanceFromId(value string) *Instance {
	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{&value},
	}
	return NewInstance(params)
}

// NewInstanceFromIp checks id is in rules
func NewInstanceFromIp(value string) *Instance {
	filters := []string{
		"network-interface.addresses.private-ip-address",
		"network-interface.ipv6-addresses.ipv6-address",
		"network-interface.addresses.association.public-ip",
	}
	for _, filter := range filters {
		params := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String(filter),
					Values: []*string{
						aws.String(value),
					},
				},
			},
		}
		if instance := NewInstance(params); instance != nil {
			return instance
		}
	}
	return nil
}
