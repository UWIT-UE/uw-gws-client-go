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
	clientConfig.CAFile = *caFile
	clientConfig.ClientCert = *certFile
	clientConfig.ClientKey = *keyFile

	gwsClient, err := gws.NewClient(clientConfig)
	if err != nil {
		log.Fatal(err)
	}

	grp1, _ := gwsClient.GetGroup("u_devtools_admin")
	fmt.Println("egid", grp1.Regid, "name", grp1.DisplayName)

	fmt.Println("Terminating the application...")
}
