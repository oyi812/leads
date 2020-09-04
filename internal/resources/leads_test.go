package resources

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"io/ioutil"
	"testing"

	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/resource/testing/mem"
	"github.com/rs/rest-layer/rest"

	// TODO clue that this test suite belongs elsewhere
	"github.com/oyi812/leads/internal/basicauth"
)

type Customer struct{}
const GroupBy = "customer"

type Lead struct {
	Id string `json:"id,omitempty"`
	Customer string `json:"customer,omitempty"`
	FirstName string
	LastName string
	Email string
	Company string `json:",omitempty"`
	PostCode string `json:",omitempty"`
	AcceptTerms bool
	DateCreated string `json:",omitempty"`
}

func TestHandler(t *testing.T) {

	username := "testuser"
	password := "password"
	customer := "testcustomer"

	// store is an in-memory implementation of
	// rest-layer/resource.Storer
	store := mem.NewHandler()
	leads := Leads{GroupBy, Customer{}}
	readWrite := resource.Conf{ AllowedModes: resource.ReadWrite }

	// /leads/[/:lead_id]
	index := resource.NewIndex()
	index.Bind("leads", leads.Schema(), store, readWrite).Use(leads.Hook())

	app, err := rest.NewHandler(index)
	if err != nil {
		t.Fatal("invalid configuration", err)
	}

	// create static ACL for basic auth
	authenticator := basicauth.StaticAuthenticator{
		ContextKey: Customer{},
		ACL: map[string]basicauth.Credentials{
			username: {password, customer},
		},
	}

	// wrap leads app in basic authenticator
	handler := basicauth.Handler(authenticator)(app)

	var a, b Lead

	// create a new lead
	a.FirstName = "bob"
	a.LastName = "cat"
	a.Email = "the@bob.cat"
	a.AcceptTerms = true

	// compose and handle authenticated request
	buf, _ := json.Marshal(&a)
	req := httptest.NewRequest("POST", "/leads", bytes.NewBuffer(buf))

	req.SetBasicAuth(username, password)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &b); err != nil {
		t.Error(err)
	}

	// response has the expected status code
	if sc := resp.StatusCode; sc != 201 {
		t.Errorf("expecting status code %d, have %d", 201, sc)
	}

	// response has a valid header
	ctk := "Content-Type"
	ctv := "application/json"
	if rct := resp.Header.Get(ctk); rct != ctv {
		t.Errorf("expecting %s %s, have %s", ctk, ctv, rct)
	}

	var have, want interface{}

	want = a.FirstName
	have = b.FirstName
	if have.(string) != want.(string) {
		t.Errorf("FirstName: expecting %q, have %q", want, have)
	}
}
