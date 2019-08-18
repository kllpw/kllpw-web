package client

import (
	"testing"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
)
func TestRegisterClient(t *testing.T) {
	var cMan = NewManager("Test1")
	req := populatedRequest()
	isClientReg := cMan.RegisterClient(httptest.NewRecorder(), req) 
	if !isClientReg {
		t.Error("No clients registered should not be valid")	
	}
	isClientReg = cMan.RegisterClient(httptest.NewRecorder(), req)
	if isClientReg {
		t.Error("Client has been registered should not be valid")	
	}
}

func TestGetClient(t *testing.T) {
	var cMan = NewManager("Test2")
	req := populatedRequest()
	client := cMan.GetClient(httptest.NewRecorder(), req) 
	if client != nil {
		t.Error("Client should not be registered yet")
	}
	isClientReg := cMan.RegisterClient(httptest.NewRecorder(), req) 
	if !isClientReg {
		t.Error("No clients registered should not be valid")	
	}
	login := cMan.LoginClient(httptest.NewRecorder(), req)
	if !login {
		t.Error("Clients registered should login")
		t.Fail()
	}
	client = cMan.GetClient(httptest.NewRecorder(), req) 
	if client == nil {
		t.Error("Client should be registered and returned")
	}
	

}

func TestIsValidClient(t *testing.T) {
	var cMan = NewManager("Test3")
	req := populatedRequest()
	isClientReg := cMan.RegisterClient(httptest.NewRecorder(), req)
	if !isClientReg {
		t.Error("No clients registered")	
	}
	valid := cMan.IsValidClient(httptest.NewRecorder(), req)
	if valid {
		t.Error("No clients registered should not be valid")
		t.Fail()
	}
	login := cMan.LoginClient(httptest.NewRecorder(), req)
	if !login {
		t.Error("Clients registered should login")
		t.Fail()
	}
	valid = cMan.IsValidClient(httptest.NewRecorder(), req)
	if !valid {
		t.Error("Clients registered should validate")
		t.Fail()
	}
	cMan.LogoutClient(httptest.NewRecorder(), req)
	valid = cMan.IsValidClient(httptest.NewRecorder(), req)
	if valid {
		t.Error("Client logged out so should not validate")
		t.Fail()
	}
}

func TestLoginClient(t *testing.T) {
	var cMan = NewManager("Test4")
	req := populatedRequest()
	isClientReg := cMan.RegisterClient(httptest.NewRecorder(), req)
	if !isClientReg {
		t.Error("No clients registered should not valid")	
	}
	login := cMan.LoginClient(httptest.NewRecorder(), req)
	if !login {
		t.Error("Clients registered should login")
		t.Fail()
	}
	badLogin := cMan.LoginClient(httptest.NewRecorder(), badRequest())
	if badLogin {
		t.Error("Invalid password should not login")
		t.Fail()
	}
}

func populatedRequest() *http.Request {
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	req.Header.Add("Authorization","Basic " + basicAuth("username1","password123"))
	return req
}

func badRequest() *http.Request {
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	req.Header.Add("Authorization","Basic " + basicAuth("username1","badpassword"))
	return req
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	 return base64.StdEncoding.EncodeToString([]byte(auth))
}
