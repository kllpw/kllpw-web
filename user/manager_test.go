package user

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	var cMan = NewManager("Test1")
	req := populatedRequest()
	isUserReg := cMan.RegisterUser(httptest.NewRecorder(), req)
	if !isUserReg {
		t.Error("No users registered should not be valid")
	}
	isUserReg = cMan.RegisterUser(httptest.NewRecorder(), req)
	if isUserReg {
		t.Error("User has been registered should not be valid")
	}
}

func TestGetUser(t *testing.T) {
	var cMan = NewManager("Test2")
	req := populatedRequest()
	client := cMan.GetUser(httptest.NewRecorder(), req)
	if client != nil {
		t.Error("User should not be registered yet")
	}
	isUserReg := cMan.RegisterUser(httptest.NewRecorder(), req)
	if !isUserReg {
		t.Error("No users registered should not be valid")
	}
	login := cMan.LoginUser(httptest.NewRecorder(), req)
	if !login {
		t.Error("Users registered should login")
		t.Fail()
	}
	client = cMan.GetUser(httptest.NewRecorder(), req)
	if client == nil {
		t.Error("User should be registered and returned")
	}
}

func TestIsValidUser(t *testing.T) {
	var cMan = NewManager("Test3")
	req := populatedRequest()
	isUserReg := cMan.RegisterUser(httptest.NewRecorder(), req)
	if !isUserReg {
		t.Error("No users registered")
	}
	valid := cMan.IsUserAuthenticated(httptest.NewRecorder(), req)
	if valid {
		t.Error("No users registered should not be valid")
		t.Fail()
	}
	login := cMan.LoginUser(httptest.NewRecorder(), req)
	if !login {
		t.Error("Users registered should login")
		t.Fail()
	}
	valid = cMan.IsUserAuthenticated(httptest.NewRecorder(), req)
	if !valid {
		t.Error("Users registered should validate")
		t.Fail()
	}
	cMan.LogoutUser(httptest.NewRecorder(), req)
	valid = cMan.IsUserAuthenticated(httptest.NewRecorder(), req)
	if valid {
		t.Error("User logged out so should not validate")
		t.Fail()
	}
}

func TestLoginUser(t *testing.T) {
	var cMan = NewManager("Test4")
	req := populatedRequest()
	isUserReg := cMan.RegisterUser(httptest.NewRecorder(), req)
	if !isUserReg {
		t.Error("No users registered should not valid")
	}
	login := cMan.LoginUser(httptest.NewRecorder(), req)
	if !login {
		t.Error("Users registered should login")
		t.Fail()
	}
	badLogin := cMan.LoginUser(httptest.NewRecorder(), badRequest())
	if badLogin {
		t.Error("Invalid password should not login")
		t.Fail()
	}
}

func populatedRequest() *http.Request {
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	req.Header.Add("Authorization", "Basic "+basicAuth("username1", "password123"))
	return req
}

func badRequest() *http.Request {
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	req.Header.Add("Authorization", "Basic "+basicAuth("username1", "badpassword"))
	return req
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
