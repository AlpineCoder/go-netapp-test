package main

import (
	"fmt"
	"net/http"
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
		Debug:             true,
	},
)
var initiator = "iqn.1994-05.com.redhat:a048115b3c7e"

var nfsClient = "baz.bar.com"
var ruleIndex = 10

var lunContainerVol = "lcvol"
var lunName = "mylun"

var lunPath = "/vol/" + lunContainerVol + "/" + lunName

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

	respa, http_respa, aerr := client.VServer.AddInitiator(os.Args[5], "trident_tst", initiator)

	fmt.Println(respa)
	fmt.Println(http_respa)
	fmt.Println(aerr)

	// Lun creation.
	// 1. A vol to hold our LUN(s)
	fmt.Println("Volume creation")
	volCreateOpts := netapp.VolumeCreateOptions{
		ContainingAggregateName: "cdot01_01_FC_1",
		Encrypt:                 false,
		ExportPolicy:            "default",
		Size:                    "1G",
		SnapshotPolicy:          "none",
		UnixPermissions:         "777",
		Volume:                  lunContainerVol,
		VolumeSecurityStyle:     "unix",
	}
	_, _, _ = client.VolumeOperations.Create(os.Args[5], &volCreateOpts)
	fmt.Println("Lun creation")

	lunCreateOpts := netapp.LunCreateOptions{
		OsType: "linux",
		Size:   104857600,
		Path:   lunPath,
	}

	lcresp, lchresp, lcerr := client.LunOperations.Create(os.Args[5], &lunCreateOpts)
	fmt.Println(lcresp)
	fmt.Println(lchresp)
	fmt.Println(lcerr)

	_, _, _ = client.LunOperations.Map(os.Args[5], lunPath, "trident_tst")

}

func removeItems() {

	respa, http_respa, aerr := client.VServer.RemoveInitiator(os.Args[5], "trident_prd", initiator, true)
	fmt.Println(respa)
	fmt.Println(http_respa)
	fmt.Println(aerr)

	r, h, e := client.VServer.ListExportRules(os.Args[5])
	fmt.Println(r)
	fmt.Println(h)
	fmt.Println(e)

	for _, rule := range r.Results.AttributesList.VServerExportRuleInfo {
		if rule.ClientMatch == nfsClient {
			client.VServer.DeleteExportRule(os.Args[5], rule.PolicyName, rule.RuleIndex)
		} else {
			fmt.Println(rule.ClientMatch + " " + strconv.Itoa(rule.RuleIndex))
		}

	}

	// offline a lun that you created!
	fmt.Println("=================== Offline LUN ===================================")
	var srr *netapp.SingleResultResponse
	var httpResp *http.Response
	var err error
	_, _, _ = client.LunOperations.Unmap(os.Args[5], lunPath, "trident_tst")

	srr, httpResp, err = client.LunOperations.Operation(os.Args[5], lunPath, netapp.LunOfflineOperation)
	if err != nil {
		fmt.Println(srr)
		fmt.Println(httpResp)
	}
	// delete a lun that you created!
	fmt.Println(httpResp)
	srr, httpResp, err = client.LunOperations.Operation(os.Args[5], lunPath, netapp.LunDestroyOperation)
	if err != nil {
		fmt.Println(srr)
		fmt.Println(httpResp)
	}

	_, _, _ = client.VolumeOperations.Operation(os.Args[5], lunContainerVol, netapp.VolumeOfflineOperation)
	_, _, _ = client.VolumeOperations.Operation(os.Args[5], lunContainerVol, netapp.VolumeDestroyOperation)
}

func main() {
	if os.Args[1] == "create" {
		createItems()
	} else if os.Args[1] == "remove" {
		removeItems()
	}
}
