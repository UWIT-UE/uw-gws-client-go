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
	clientConfig.APIUrl = "https://eval.groups.uw.edu/group_sws/v3"
	clientConfig.CAFile = *caFile
	clientConfig.ClientCert = *certFile
	clientConfig.ClientKey = *keyFile

	gwsClient, err := gws.NewClient(clientConfig)
	if err != nil {
		log.Fatal(err)
	}

	grp1, err := gwsClient.GetGroup("u_devtools_admin")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("egid", grp1.Regid, "name", grp1.DisplayName)

	members1, err := gwsClient.GetEffectiveMembership("u_devtools_admin")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("membership", members1)

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
	// fmt.Println("egid", grp2.Regid, "name", grp2.DisplayName)

	err = gwsClient.DeleteGroup("u_unixgrp_testgroup3")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Terminating the application...")
}
