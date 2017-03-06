package resource

import (
  "fmt"
  "log"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ec2"

)

type Instance struct {
  ec2.Instance
}

func (i *Instance) CheckSecurityGroups(protocol string, ip_address string, port string) bool {
  fmt.Println(protocol)
  fmt.Println(ip_address)
  fmt.Println(port)
  return true
}

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
    for _, instance := range reservation.Instances{
      ids = append(ids, instance)
    }
  }

  return &Instance{*ids[0]}
}
