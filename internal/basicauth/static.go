package basicauth

type Credentials struct {
	Password string
	Payload interface{}
}

type StaticAuthenticator struct {
	ContextKey interface{}
	ACL map[string]Credentials
}

// BasicAuth accepts basic auth credentials
// looks up the username in it's ACL and
// athenticates by comparing passwords. If 
// successful it returns the key and value
// to be added to the connection context
func (a StaticAuthenticator) BasicAuth(
	username string,
	password string,
) (
	key, value interface{},
	ok bool,
	err error,
) {
	credentials, ok := a.ACL[username]
	if !ok || password != credentials.Password {
		return nil, nil, false, nil
	}

	return a.ContextKey, credentials.Payload, true, nil
}
