package main

import (
	"log"
	"flag"
	"net/http"

	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/resource/testing/mem"
	"github.com/rs/rest-layer/rest"

	"github.com/oyi812/leads/internal/basicauth"
	"github.com/oyi812/leads/internal/resources"
)

var hostFlag = flag.String("host", "localhost", "restful service host")
var portFlag = flag.String("port", "8080", "restful service port")

// Customer is used for identifying
// the grouping/authentication field value
// in the context of the connection
type Customer struct{}
const GroupBy = "customer"

func main() {

	// config
	flag.Parse()
	addr := *hostFlag + ":" + *portFlag

	// store is an in-memory implementation of
	// rest-layer/resource.Storer
	store := mem.NewHandler()
	leads := resources.Leads{GroupBy, Customer{}}
	readWrite := resource.Conf{ AllowedModes: resource.ReadWrite }

	// /leads/[/:lead_id]
	index := resource.NewIndex()
	index.Bind("leads", leads.Schema(), store, readWrite).Use(leads.Hook())

	app, err := rest.NewHandler(index)
	if err != nil {
		log.Fatalln("invalid configuration", err)
	}

	// create static ACL for basic auth
	authenticator := basicauth.StaticAuthenticator{
		ContextKey: Customer{},
		ACL: map[string]basicauth.Credentials{
			"testuser": {"password", "testcustomer"},
		},
	}

	// wrap leads app in basic authenticator
	handler := basicauth.Handler(authenticator)(app)

	log.Println("listening on", addr)
	err = http.ListenAndServe(addr, handler)
	log.Fatalln("stopped listening on", addr, "error:", err)
}
