// Copyright © 2017 Sean Rankine <srdeveloper@icloud.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"net"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/theseanything/marnie/resource"
)

type IdentifierType uint8

const (
	IpAddress = iota
	InstanceId
	NameTag
	None
)

type Target struct {
	value  string
	idType IdentifierType
}

var conn resource.Connection

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println("No arguments given.") // return instead
		}

		var instance *resource.Instance
		var err error
		identifier := args[0]
		identifierType := parseIdentifier(identifier)

		fmt.Println("Searching for EC2 instance...")
		switch identifierType {
		case IpAddress:
			instance, err = resource.NewInstanceFromIp(identifier)
		case InstanceId:
			instance, err = resource.NewInstanceFromId(identifier)
		case NameTag:
			instance, err = resource.NewInstanceFromNameTag(identifier)
		}

		if err != nil {
			fmt.Println("Error: ", err)
		}

		if instance != nil {
			fmt.Println(*instance.InstanceId)
			instance.CheckSecurityGroups(conn)
		}

	},
}

func init() {
	checkCmd.Flags().StringVarP(&conn.protocol, "protocol", "p", "tcp", "protocol used")
	checkCmd.Flags().StringVarP(&conn.destIp, "dest-ip", "d", "22", "ip of the destination")
	checkCmd.Flags().StringVarP(&conn.destPort, "dest-port", "D", "22", "port of the destination")
	checkCmd.Flags().StringVarP(&conn.srcIp, "src-ip", "s", "", "ip of the source")
	checkCmd.Flags().StringVarP(&conn.srcPort, "src-port", "S", "", "port of the source")

	checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.AddCommand(checkCmd)
}

func parseIdentifier(value string) IdentifierType {
	if ip := net.ParseIP(value); ip != nil {
		return IpAddress
	} else if regexp.MustCompile("^i-[0-9a-f]{17}$").MatchString(value) {
		return InstanceId
	} else {
		return NameTag
	}
}
