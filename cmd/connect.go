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

	"github.com/spf13/cobra"
	"github.com/theseanything/marnie/resource"
)

// ResourceType is
type ResourceType uint8

const (
	External = iota
	Instance
)

var sourceProtocol string
var sourceIP string
var sourcePort int
var sourceType ResourceType

var nameTag string
var instanceID string
var targetIP string
var destinationType ResourceType

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		var sourceResource Resource
		var destinationResource Resource

		sourceResource.CheckConnectionTo(protocol, destinationResource, port)
		nameTag := args[0]

		fmt.Print("Searching for instance with name tag \"", nameTag, "\": ")
		instance := *resource.NewInstanceFromNameTag(nameTag)
		fmt.Println(*instance.InstanceId)

		response := "NOPE"
		if instance.CheckRouteTables(sourceProtocol, sourceIP, sourcePort) {
			response = "OK"
		}
		fmt.Println("Route Tables: ", response)

		response = "ERROR"
		if instance.CheckSecurityGroups(sourceProtocol, sourceIP, sourcePort) {
			response = "OK"
		}
		fmt.Println("Security Groups: ", response)

		response = "ERROR"
		if instance.CheckNACLs(sourceProtocol, sourceIP, sourcePort) {
			response = "OK"
		}
		fmt.Println("NACLs: ", response)

	},
}

func init() {
	connectCmd.Flags().StringVarP(&sourceProtocol, "protocol", "p", "tcp", "The protcol being used to connect.")
	connectCmd.Flags().StringVarP(&sourceIP, "ip", "i", "10.0.0.0", "The source IP being used to connect from.")
	connectCmd.Flags().IntVarP(&sourcePort, "port", "P", 80, "The resource port attempt to connect to.")
	RootCmd.AddCommand(connectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// connectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// connectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
