package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/uwit-ue/uw-gws-client-go/gws"
)

var (
	certFile = flag.String("cert", "someCertFile", "A PEM eoncoded certificate file.")
	keyFile  = flag.String("key", "someKeyFile", "A PEM encoded private key file.")
	caFile   = flag.String("CA", "someCertCAFile", "A PEM eoncoded CA's certificate file.")
)

func main() {
	var grp gws.GroupResponse

	flag.Parse()

	fmt.Println("Starting the application...")

	// Vault API calling
	clientConfig := &gws.Config{}
	tlsConfig := &gws.TLSConfig{
		CACert:     *caFile,
		ClientCert: *certFile,
		ClientKey:  *keyFile,
	}
	if err := clientConfig.ConfigureTLS(tlsConfig); err != nil {
		log.Fatal(err)
	}

	gwsClient, err := gws.NewClient(clientConfig)
	if err != nil {
		log.Fatal(err)
	}

	resp1, err := gwsClient.Config.HTTPClient.GET("https://iam-ws.u.washington.edu/group_sws/v3/group/u_devtools_admin")
	if err != nil {
		log.Fatal(err)
	}

	// Dump response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(data))

	fmt.Println("Part Two...")

	// Load client cert
	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(*caFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig2 := &tls.Config{
		Certificates:  []tls.Certificate{cert},
		RootCAs:       caCertPool,
		Renegotiation: tls.RenegotiateOnceAsClient,
	}
	tlsConfig2.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig2}
	client := &http.Client{Transport: transport}

	// Do GET something
	resp, err := client.Get("https://iam-ws.u.washington.edu/group_sws/v3/group/u_devtools_admin")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Dump response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(data))

	err = json.Unmarshal(data, &grp)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("here are my properties")
	log.Println(grp.Meta.ID)
	log.Println(grp.Data.Gid)
	//log.Println(grp.Created)

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
