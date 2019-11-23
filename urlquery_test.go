package goutils

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type API struct {
	Client  *http.Client
	baseURL string
}

func (api *API) HandleParamsNo(endpoint string) ([]byte, error) {
	resp, err := api.Client.Get(api.baseURL + endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// handling error and doing stuff with body that needs to be unit tested
	return body, err
}

func (api *API) HandleParamsYes() ([]byte, error) {
	resp, err := api.Client.Get(api.baseURL + "/params/yes?str='foobar'&int=5&bool=false")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// handling error and doing stuff with body that needs to be unit tested
	return body, err
}

func TestHandleNoParamsTrue(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, req.URL.String(), "/params/no")
		ok, err := EnsureNoQueryParameters(req)
		if err != nil {
			rw.Write([]byte(`ERROR`))
			return
		}
		if ok {
			// Send response to be tested
			rw.Write([]byte(`OK`))
		}
	}))
	// Close the server when test finishes
	defer server.Close()

	api := API{server.Client(), server.URL}

	// Use Client & URL from our local test server noparams endpoint for no parameters
	body, err := api.HandleParamsNo("/params/no")
	ok(t, err)
	equals(t, []byte("OK"), body)
}

func TestHandleNoParamsFalse(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, req.URL.String(), "/params/no?str='foobar'&int=5&bool=false")
		ok, err := EnsureNoQueryParameters(req)
		if err != nil {
			rw.Write([]byte(`OK`))
			return
		}
		if ok {
			// Send response to be tested
			rw.Write([]byte(`ERROR`))
		}
	}))
	// Close the server when test finishes
	defer server.Close()

	api := API{server.Client(), server.URL}

	// Use Client & URL from our local test server noparams endpoint for no parameters
	body, err := api.HandleParamsNo("/params/no?str='foobar'&int=5&bool=false")
	ok(t, err)
	equals(t, []byte("OK"), body)

}

func TestHandleGetQueryParamsOk(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, req.URL.String(), "/params/yes?str=foobar&int=5&bool=false")
		qp, err := GetQueryParameter(req, "str", true, false, "")
		if err != nil {
			//equals(t, qp, "foobar")
			rw.Write([]byte(`ERROR`))
			return
		}
		if qp == "foobar" {
			// Send response to be tested
			rw.Write([]byte(`OK`))
		}
	}))
	// Close the server when test finishes
	defer server.Close()

	api := API{server.Client(), server.URL}

	// Use Client & URL from our local test server noparams endpoint for no parameters
	body, err := api.HandleParamsNo("/params/yes?str=foobar&int=5&bool=false")
	ok(t, err)
	equals(t, []byte("OK"), body)
}

func TestHandleGetQueryParamsDefault(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, req.URL.String(), "/params/yes?int=5&bool=false")
		qp, err := GetQueryParameter(req, "str", true, true, "foobar")
		if err != nil {
			//equals(t, qp, "foobar")
			rw.Write([]byte(`ERROR`))
			return
		}
		if qp == "foobar" {
			// Send response to be tested
			rw.Write([]byte(`OK`))
		}
	}))
	// Close the server when test finishes
	defer server.Close()

	api := API{server.Client(), server.URL}

	// Use Client & URL from our local test server noparams endpoint for no parameters
	body, err := api.HandleParamsNo("/params/yes?int=5&bool=false")
	ok(t, err)
	equals(t, []byte("OK"), body)
}

func TestHandleGetQueryParamsRequireNoDefault(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, req.URL.String(), "/params/yes?int=5&bool=false")
		_, err := GetQueryParameter(req, "str", true, false, "")
		if err != nil {
			//equals(t, qp, "foobar")
			rw.Write([]byte(`OK`))
			return
		}

		rw.Write([]byte(`ERROR`))

	}))
	// Close the server when test finishes
	defer server.Close()

	api := API{server.Client(), server.URL}

	// Use Client & URL from our local test server noparams endpoint for no parameters
	body, err := api.HandleParamsNo("/params/yes?int=5&bool=false")
	ok(t, err)
	equals(t, []byte("OK"), body)
}

func TestHandleGetQueryParamsNotRequireDefault(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, req.URL.String(), "/params/yes?int=5&bool=false")
		qp, err := GetQueryParameter(req, "str", false, true, "foobar")
		if err != nil {
			rw.Write([]byte(`ERROR`))
			return
		}
		if qp == "foobar" {
			// Send response to be tested
			rw.Write([]byte(`OK`))
		}

	}))
	// Close the server when test finishes
	defer server.Close()

	api := API{server.Client(), server.URL}

	// Use Client & URL from our local test server noparams endpoint for no parameters
	body, err := api.HandleParamsNo("/params/yes?int=5&bool=false")
	ok(t, err)
	equals(t, []byte("OK"), body)
}

func TestHandleGetQueryParamsNotRequireNoDefault(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, req.URL.String(), "/params/yes?int=5&bool=false")
		qp, err := GetQueryParameter(req, "str", false, false, "")
		if err != nil {
			rw.Write([]byte(`ERROR`))
			return
		}
		if qp == "" {
			// Send response to be tested
			rw.Write([]byte(`OK`))
		}

	}))
	// Close the server when test finishes
	defer server.Close()

	api := API{server.Client(), server.URL}

	// Use Client & URL from our local test server noparams endpoint for no parameters
	body, err := api.HandleParamsNo("/params/yes?int=5&bool=false")
	ok(t, err)
	equals(t, []byte("OK"), body)
}
