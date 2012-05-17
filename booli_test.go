package booli

import ( 	
			"testing"
			"net/http"
		)

// Simulate other then 200 response, for example 403
type HttpGetFail struct {
}

func (h *HttpGetFail) Get(url string) (r *http.Response, err error) {
	// Todo: make broken response
	resp := http.Response{}
	return &resp, nil
}

// Simulate broken json from booli
type HttpGetBrokenJson struct {
}

func (h *HttpGetBrokenJson) Get(url string) (r *http.Response, err error) {
	// Todo: make broken json response
	resp := http.Response{}
	return &resp, nil
}

// Simulate a ok response from booli
type HttpGetOk struct {
}

func (h *HttpGetOk) Get(url string) (r *http.Response, err error) {
	// Todo: make Ok booli response
	resp := http.Response{}
	return &resp, nil
}

// Check a full search condition string
type HttpGetCheckSearchCondition struct {
}

func (h *HttpGetCheckSearchCondition) Get(url string) (r *http.Response, err error) {
	// Todo: make Ok booli response
	resp := http.Response{}
	return &resp, nil
}
		
func TestGetResultImpl (t *testing.T) {
	// Test caller id empty
	_, err := GetResultImpl(SearchCondition{}, "", "xxx", nil)
	if err == nil {
		t.Errorf("Should be missing caller id error!")
	}
	
	// Test key empty
	_, err = GetResultImpl(SearchCondition{}, "xxx", "", nil)
	if err == nil {
		t.Errorf("Should be key empty error!")
	}
	
	// Test missing vital searchconditions
	_, err = GetResultImpl(SearchCondition{}, "xxx", "xxx", nil)
	if err == nil {
		t.Errorf("Should be missing search condition error")
	}
}
