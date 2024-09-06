package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/uwit-ue/uw-gws-client-go/gws"
)

var (
	certFile = flag.String("cert", "someCertFile", "A PEM eoncoded certificate file.")
	keyFile  = flag.String("key", "someKeyFile", "A PEM encoded private key file.")
	caFile   = flag.String("CA", "someCertCAFile", "A PEM eoncoded CA's certificate file.")
)

func main() {
	flag.Parse()
	fmt.Println("Starting the application...")

	clientConfig := gws.DefaultConfig()
	//clientConfig.APIUrl = "https://eval.groups.uw.edu/group_sws/v3"
	clientConfig.APIUrl = "https://groups.uw.edu/group_sws/v3"
	clientConfig.CAFile = *caFile
	clientConfig.ClientCert = *certFile
	clientConfig.ClientKey = *keyFile

	gwsClient, err := gws.NewClient(clientConfig)
	if err != nil {
		log.Fatal(err)
	}

	// TEST getgroup
	grp1, err := gwsClient.GetGroup("u_devtools_admin")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("egid", grp1.Regid, "name", grp1.DisplayName)

	// TEST getmembership
	// members1, err := gwsClient.GetMembership("u_devtools_admin")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("membership", members1)

	// TEST geteffectivemembership
	// members2, err := gwsClient.GetEffectiveMembership("u_erich_membertypes")
	// //members2, err := gwsClient.GetMembership("u_erich_membertypes")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("eff membership", members2)

	// fmt.Println("eff membership comma", members2.Match(gws.MemberTypeUWNetID).ToCommaString())

	// TEST MemberCount
	// memberC, err := gwsClient.MemberCount("u_devtools_admin")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("membership count", memberC)

	// TEST EffectiveMemberCount
	// memberC2, err := gwsClient.EffectiveMemberCount("u_devtools_admin")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("eff membership count", memberC2)

	// TEST getmember
	// member1, err := gwsClient.GetEffectiveMember("u_devtools_admin", "erich5")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("id", member1.ID, "type", member1.Type)

	// TEST ismember
	// ismem, err := gwsClient.IsEffectiveMember("u_devtools_admin2", "erich2")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("erich membership:", ismem)

	// TEST creategroup
	// fmt.Println("create group")
	// newg := &gws.Group{
	// 	ID:          "u_unixgrp_testgroup3",
	// 	DisplayName: "A test group u_unixgrp_testgroup3",
	// 	Description: "lalala",
	// }

	// grp2, err := gwsClient.CreateGroup(newg)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("grp2 %+v\n", grp2)
	// fmt.Println("regid", grp2.Regid, "name", grp2.DisplayName)
	// // fmt.Println("sleep")
	// // time.Sleep(30 * time.Second)

	// // Example for updating a group
	// fmt.Println("update group")
	// origGrp, err := gwsClient.GetGroup("u_unixgrp_testgroup3")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// origGrp.DisplayName = "Updated " + origGrp.DisplayName
	// updatedGrp, err := gwsClient.UpdateGroup(origGrp)
	// fmt.Println("updated display name:", updatedGrp.DisplayName)

	// // TEST deletegroup
	// err = gwsClient.DeleteGroup("u_unixgrp_testgroup3")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// ss := gws.NewSearch()
	// ss = ss.WithName("u_devtools_admin")
	// //ss = ss.WithMember("erich1")
	// //ss = ss.OnlyDirectMembers()
	// ss = ss.InEffectiveMembers()
	// i, err := gwsClient.DoSearch(ss)
	// fmt.Println("returned", i, err)
	// fmt.Println("first group", i[0].DisplayName)

	// ss2 := gws.NewSearch().WithMember("erich1").InEffectiveMembers()
	// gresult, err := gwsClient.DoSearch(ss2)

	// if err == nil {
	// 	for _, g := range gresult {
	// 		fmt.Println(g.ID)
	// 	}
	// }

	// nf1, err := gwsClient.AddMembers("u_clos_test", "clos", "erich0", "erich1", "erich2")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("nf", len(nf1), strings.Join(nf1, ", "))

	// time.Sleep(10 * time.Second)

	// members1, err := gwsClient.GetMembership("u_clos_test")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("membership", members1)

	// time.Sleep(2 * time.Second)

	// err = gwsClient.DeleteMembers("u_clos_test", "erich1", "erich2", "notfoundmember333")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Removing")

	// err = gwsClient.DeleteAllMembers("u_clos_test")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Removing All")

	// time.Sleep(10 * time.Second)

	// members2, err := gwsClient.GetMembership("u_clos_test")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("membership", members2)

	// members5 := gws.NewMemberList()
	// members5 = members5.AddUWNetIDMembers("erich1", "erich2", "erich0")
	// fmt.Println("manipulated membership", members5)
	// // members5.AddDNSMembers("ref.s.uw.edu", "clos.s.uw.edu")
	// // fmt.Println("manipulated membership", members5)
	// // members5.AddEPPNMembers("erich@quavy.com")
	// // fmt.Println("manipulated membership", members5)
	// // members5.AddDNSMembers("erich")

	// gwsClient.SetMembership("u_erich_wasempty", members5)
	// fmt.Println("Terminating the application...")

	// Manipulating group admins/owners

	// gwsClient, err := gws.NewClient(clientConfig)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// ss := gws.NewSearch()
	// ss = ss.WithOwner("somename.cac.washington.edu")
	// grps, err := gwsClient.DoSearch(ss)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // iterate over grps and update the admin list
	// for _, grp := range grps {
	// 	fmt.Println("Change group:", grp.ID)

	// 	fullGrp, err := gwsClient.GetGroup(grp.ID)
	// 	if err != nil {
	// 		fmt.Println("Error fetching group ", grp.ID, ":", err)
	// 		continue
	// 	}

	// 	fmt.Println("   old admins:", fullGrp.Admins.ToCommaString())
	// 	if fullGrp.IsAdmin("somename.cac.washington.edu") && !fullGrp.IsAdmin("newname.s.uw.edu") {
	// 		fullGrp.AddAdmin("newname.s.uw.edu")
	// 		fmt.Println("   new admins:", fullGrp.Admins.ToCommaString())
	// 	} else {
	// 		continue
	// 	}
	// 	updatedGrp, err := gwsClient.UpdateGroup(fullGrp)
	// 	if err != nil {
	// 		fmt.Println("Error updating group ", fullGrp.ID, ":", err)
	// 		continue
	// 	}
	// 	fmt.Println("   updated admins:", updatedGrp.Admins.ToCommaString())
	// }

}
