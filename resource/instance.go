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
		s.CheckIngress(protocol, ipAddress, port)
		s.CheckEgress(protocol, ipAddress, port)
		securityGroups = append(securityGroups, &s)
	}

	fmt.Println(response)

	return true
}

// NewInstanceFromNameTag ityGroups checks id is in rules
func NewInstanceFromNameTag(value string) *Instance {
	sess := session.Must(session.NewSession())
	svc := ec2.New(sess)
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

	return &Instance{*ids[0]}
}
