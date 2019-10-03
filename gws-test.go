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
	//var grp gws.GroupResponse

	flag.Parse()

	fmt.Println("Starting the application...")

	// Vault API calling
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

// jsonData := map[string]string{"firstname": "Nic", "lastname": "Raboy"}
// jsonValue, _ := json.Marshal(jsonData)
// request, _ := http.NewRequest("POST", "https://httpbin.org/post", bytes.NewBuffer(jsonValue))
// request.Header.Set("Content-Type", "application/json")
// client := &http.Client{}
// response, err := client.Do(request)
// if err != nil {
// 	fmt.Printf("The HTTP request failed with error %s\n", err)
// } else {
// 	data, _ := ioutil.ReadAll(response.Body)
// 	fmt.Println(string(data))
// }
