package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

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
	clientConfig.APIUrl = "https://eval.groups.uw.edu/group_sws/v3"
	clientConfig.CAFile = *caFile
	clientConfig.ClientCert = *certFile
	clientConfig.ClientKey = *keyFile

	gwsClient, err := gws.NewClient(clientConfig)
	if err != nil {
		log.Fatal(err)
	}

	// TEST getgroup
	// grp1, err := gwsClient.GetGroup("u_devtools_admin")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("egid", grp1.Regid, "name", grp1.DisplayName)

	// TEST getmembership
	// members1, err := gwsClient.GetMembership("u_devtools_admin")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("membership", members1)

	// TEST geteffectivemembership
	// members2, err := gwsClient.GetEffectiveMembership("u_devtools_admin")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("eff membership", members2)
	// fmt.Println("slice size, cap", len(members2), cap(members2))

	// TEST getmembercount
	// memberC, err := gwsClient.GetMemberCount("u_devtools_admin")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("membership count", memberC)

	// TEST geteffectivemembercount
	// memberC2, err := gwsClient.GetEffectiveMemberCount("u_devtools_admin")
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
	// newg := &gws.Group{
	// 	ID:          "u_unixgrp_testgroup3",
	// 	DisplayName: "A test group u_unixgrp_testgroup3",
	// 	Description: "lalala",
	// 	Admins:      gws.ToEntityList(&gws.Entity{Type: "uwnetid", ID: "erich"}),
	// }

	// grp2, err := gwsClient.CreateGroup(*newg)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("grp2 %+v\n", grp2)
	// fmt.Println("regid", grp2.Regid, "name", grp2.DisplayName)
	// fmt.Println("sleep")
	// time.Sleep(30 * time.Second)

	// Example for updating a group
	// origGrp, err := gwsClient.GetGroup("u_unixgrp_testgroup3")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// origGrp.DisplayName = "Updated " + origGrp.DisplayName
	// updatedGrp, err := gwsClient.UpdateGroup(origGrp)
	// fmt.Println("updated display name:", updatedGrp.DisplayName)

	// TEST deletegroup
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

	nf1, err := gwsClient.AddMembers("u_clos_test", "clos", "erich0", "erich1", "erich2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("nf", len(nf1), strings.Join(nf1, ", "))

	time.Sleep(10 * time.Second)

	members1, err := gwsClient.GetMembership("u_clos_test")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("membership", members1)

	time.Sleep(2 * time.Second)

	// err = gwsClient.RemoveMembers("u_clos_test", "erich1", "erich2", "notfoundmember333")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Removing")

	err = gwsClient.RemoveAllMembers("u_clos_test")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Removing All")

	time.Sleep(10 * time.Second)

	members2, err := gwsClient.GetMembership("u_clos_test")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("membership", members2)

	fmt.Println("Terminating the application...")
}
