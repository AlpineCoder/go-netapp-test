package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AlpineCoder/go-netapp/netapp"
)

var client = netapp.NewClient(
	os.Args[2],
	"1.160",
	&netapp.ClientOptions{
		BasicAuthUser:     os.Args[3],
		BasicAuthPassword: os.Args[4],
		SSLVerify:         false,
	},
)
var initiator = "iqn.1994-05.com.redhat:a048115b3c7e"

var nfsClient = "baz.bar.com"
var ruleIndex = 10

func createItems() {

	protocols := []string{"NFS"}
	anyrule := []string{"Any"}

	resp, http_resp, err := client.VServer.CreateExportRule(os.Args[5], &netapp.VServerExportRuleInfo{
		ClientMatch:       nfsClient,
		PolicyName:        "default",
		RuleIndex:         ruleIndex,
		Protocol:          &protocols,
		ReadOnlyRule:      &anyrule,
		ReadWriteRule:     &anyrule,
		SuperUserSecurity: &anyrule,
	})

	fmt.Println(resp)
	fmt.Println(http_resp)
	fmt.Println(err)

	respa, http_respa, aerr := client.VServer.AddInitiator(os.Args[5], "trident", initiator)

	fmt.Println(respa)
	fmt.Println(http_respa)
	fmt.Println(aerr)

}

func removeItems() {

	respa, http_respa, aerr := client.VServer.RemoveInitiator(os.Args[5], "trident", initiator, true)
	fmt.Println(respa)
	fmt.Println(http_respa)
	fmt.Println(aerr)

	r, h, e := client.VServer.ListExportRules(os.Args[5])
	fmt.Println(r)
	fmt.Println(h)
	fmt.Println(e)

	for _, rule := range r.Results.AttributesList.VServerExportRuleInfo {
		if rule.ClientMatch == nfsClient && rule.RuleIndex == ruleIndex {
			client.VServer.DeleteExportRule(os.Args[5], rule.PolicyName, rule.RuleIndex)
		} else {
			fmt.Println(rule.ClientMatch + " " + strconv.Itoa(rule.RuleIndex))
		}

	}
}

func main() {
	if os.Args[1] == "create" {
		createItems()
	} else if os.Args[1] == "remove" {
		removeItems()
	}
}
