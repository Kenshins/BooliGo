package booli

import ( 	
			"testing"
			"net/http"
			"errors"
			"strings"
			"io"
		)
		
const (
		GetNonAuth = 1
		GetOk = 2
		GetCheckSearchCond = 3
	  )

// A mock io.ReadCloser
type ReadCloser struct {
	Reader
	Closer
}

type Reader struct {
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n = copy(p, []byte("{}"))
	return n, io.EOF
}

type Closer struct {
}

func (c *Closer) Close() error {
	return nil
}

// Mock get struct to test different behaviours of different responses and search conditions
type MockHttpGet struct {
	UrlMatch string
	TestType int
}

func (h *MockHttpGet) Get(url string) (r *http.Response, err error) {
	switch h.TestType {
	case GetCheckSearchCond:
		resp := http.Response{Status: "200 OK", StatusCode: 200, Body: &ReadCloser{}}
		urlParts := strings.Split(url,"&callerId")
		if h.UrlMatch != urlParts[0] {
			return &resp, errors.New("Missmatch in constructed url and expected url in unit test!, expected \n" + h.UrlMatch + " was \n" +  urlParts[0])
		}
		return &resp, nil
	case GetOk:
		resp := http.Response{Status: "200 OK", StatusCode: 200, Body: &ReadCloser{}}
		return &resp, nil
	case GetNonAuth:
		// Todo: make broken response
		resp := http.Response{Status: "403 OK", StatusCode: 403, Body: &ReadCloser{}}
		return &resp, nil
	}
	return r, errors.New("Missing Test type to use mock get in unit test!")	
}

type searchMatchPos struct {
	SearchCond SearchCondition
	UrlMatch string
}

type searchMatchNeg struct {
	SearchCond SearchCondition
	ExpectedError string
}

var searchConditionPositiveTests = []searchMatchPos { {SearchCondition{Q: "nacka"}, "http://api.booli.se/listings?offset=0&limit=3&q=nacka"},
												 {SearchCondition{Q: "svapasjarvi"}, "http://api.booli.se/listings?offset=0&limit=3&q=svapasjarvi" }, 
												 {SearchCondition{Q: "nacka", Center: "20,20", Dim: "300,300", Bbox: "1,1,1,1"}, "http://api.booli.se/listings?offset=0&limit=3&bbox=1,1,1,1&dim=300,300&center=20,20&q=nacka" },
												 {SearchCondition{Q: "nacka", Center: "1,1", Dim: "1,1", Bbox: "1,1,1,1", AreaId: "1,2,3", MinPrice: 200000, MaxPrice: 2000000, MinRooms: 2, MaxRooms: 4, MaxRent: 500, MinLivingArea: 10, MaxLivingArea: 500, MinPlotArea: 200, MaxPlotArea: 6000, ObjectType: "villa, radhus", MinCreated: "20100101", MaxCreated: "20100115", Limit:0, Offset:0}, "http://api.booli.se/listings?offset=0&limit=3&maxCreated=20100115&minCreated=20100101&objectType=villa, radhus&maxPlotArea=6000&minPlotArea=200&maxLivingArea=500&minLivingArea=10&maxRent=500&maxRooms=4&minRooms=2&maxPrice=2000000&minPrice=200000&areaId=1,2,3&bbox=1,1,1,1&dim=1,1&center=1,1&q=nacka" }}

var searchConditionNegativeTests = []searchMatchNeg { {SearchCondition{Q: "nacka", Center: "1,1"}, "Missing Dim!"},
													  {SearchCondition{Q: "nacka", Dim: "1,1"}, "Missing Center!"}, 
													  {SearchCondition{Q: "nacka", Dim: "-1,1", Center: "1,1"}, "Negative Dim!"},
													  {SearchCondition{Q: "nacka", Dim: "f,1", Center: "1,1"}, "Non number input Dim!"},
													  {SearchCondition{Q: "nacka", Dim: "1,1,1", Center: "1,1"}, "To many args to Dim!"},
													  {SearchCondition{Q: "nacka", Dim: "1", Center: "1,1"}, "To few args to Dim!"},
													  {SearchCondition{Q: "nacka", Dim: "1,1", Center: "-91,1"}, "Lat must be between -90 to 90!"},
													  {SearchCondition{Q: "nacka", Dim: "1,1", Center: "91,1"}, "Lat must be between -90 to 90!"},
													  {SearchCondition{Q: "nacka", Dim: "1,1", Center: "1,-181"}, "Long must be between -180 to 180!"},
													  {SearchCondition{Q: "nacka", Dim: "1,1", Center: "1,181"}, "Long must be between -180 to 180!"},
													  {SearchCondition{Q: "nacka", MaxPrice: -20}, "MaxPrice can not be negative!"}}
												 
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
	
	// Test missing vital search condition
	_, err = GetResultImpl(SearchCondition{}, "xxx", "xxx", &MockHttpGet{UrlMatch: "http://api.booli.se/listings?offset=0&limit=3&q=nacka", TestType: GetCheckSearchCond})
	if err == nil {
		t.Errorf("Error, should alert for missing search conditions!")
	}
	
	// Test	search conditions positive tests
	for i := range searchConditionPositiveTests {
		_, err = GetResultImpl(searchConditionPositiveTests[i].SearchCond, "xxx", "xxx", &MockHttpGet{UrlMatch: searchConditionPositiveTests[i].UrlMatch, TestType: GetCheckSearchCond})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}
	
	// Test	search conditions negative tests
	for i := range searchConditionNegativeTests {
		_, err = GetResultImpl(searchConditionNegativeTests[i].SearchCond, "xxx", "xxx", &MockHttpGet{TestType: GetOk})
		if err == nil {
			t.Errorf("The current test should produce error: " + searchConditionNegativeTests[i].ExpectedError)
		}
	}
	
	// Test wrong auth
	_, err = GetResultImpl(SearchCondition{Q:"nacka"}, "xxx", "xxx", &MockHttpGet{UrlMatch: "http://api.booli.se/listings?offset=0&limit=3&q=nacka", TestType: GetNonAuth})
	if err == nil {
		t.Errorf("Error, should alert for wrong auth!")
	}
}
