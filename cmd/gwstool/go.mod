module github.com/uwit-ue/uw-gws-client-go/cmd/gwstool

go 1.23.1

require (
	github.com/spf13/cobra v1.9.1
	github.com/uwit-ue/uw-gws-client-go v0.0.0
)

require (
	github.com/go-resty/resty/v2 v2.16.5 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	golang.org/x/net v0.33.0 // indirect
)

// Use local version of the library
replace github.com/uwit-ue/uw-gws-client-go => ../..
