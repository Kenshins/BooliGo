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

type searchCondList struct {
	sCond SearchCondition
	lCond ListingsExtendedSearchCondition
}

type searchCondSold struct {
	sCond SearchCondition
	seCond SoldExtendedSearchCondition
}

type searchMatchPos struct {
	SearchCond SearchCondition
	UrlMatch string
}

type searchMatchListPos struct {
	 sCondList searchCondList
	 UrlMatch string
}

type searchMatchSoldPos struct {
	sCondSold searchCondSold
	UrlMatch string
}

type searchMatchNeg struct {
	SearchCond SearchCondition
	ExpectedError string
}

var searchConditionPositiveTests = []searchMatchPos { {SearchCondition{Q: "nacka"}, "http://api.booli.se/listings?offset=0&limit=3&q=nacka"},
												 {SearchCondition{Q: "svapasjarvi"}, "http://api.booli.se/listings?offset=0&limit=3&q=svapasjarvi" }, 
												 {SearchCondition{Q: "nacka", Center: "20,20", Dim: "300,300", Bbox: "1,1,1,1"}, "http://api.booli.se/listings?offset=0&limit=3&bbox=1,1,1,1&dim=300,300&center=20,20&q=nacka" },
												 {SearchCondition{Q: "nacka", Center: "1,1", Dim: "1,1", Bbox: "-1,1,1,-1", AreaId: "1,2,3", Rooms: Rooms{MinRooms: 2, MaxRooms: 4}, MaxRent: 500, LivingArea: LivingArea{MinLivingArea: 10, MaxLivingArea: 500}, MinPlotArea: 200, MaxPlotArea: 6000, ObjectType: "villa, radhus", MinPublished: "20100101", MaxPublished: "20100115", Limit:0, Offset:0}, "http://api.booli.se/listings?offset=0&limit=3&maxPublished=20100115&minPublished=20100101&objectType=villa, radhus&maxPlotArea=6000&minPlotArea=200&maxLivingArea=500&minLivingArea=10&maxRent=500&maxRooms=4&minRooms=2&areaId=1,2,3&bbox=-1,1,1,-1&dim=1,1&center=1,1&q=nacka" }}

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
													  {SearchCondition{Q: "nacka", MaxPublished: "2001-11-02"}, "MaxCreated not conforming to YYYYMMDD!"},
													  {SearchCondition{Q: "nacka", MinPublished: "2001-11-02"}, "MinCreated not conforming to YYYYMMDD!"},
													  {SearchCondition{Q: "nacka", AreaId: "ff,1,22"}, "AreaID must be in the format 55,44...!"},
													  {SearchCondition{Q: "nacka", AreaId: "ff"}, "AreaID must be in the format 55,44...!"},
													  {SearchCondition{Q: "nacka", Bbox: "1,1,1"}, "Bbox must be 1,1,1,1!"},
													  {SearchCondition{Q: "nacka", Bbox: "1,1,1,f"}, "Bbox must be 1,1,1,1!"},
													  {SearchCondition{Q: "nacka", Bbox: "-91,1,1,1"}, "Bbox must be 1,1,1,1!"},
													  {SearchCondition{Q: "nacka", Bbox: "1,1,1,-181"}, "Bbox must be 1,1,1,1!"}}
													  
var searchConditionListingsTests = []searchMatchListPos {{searchCondList{sCond: SearchCondition{Q: "nacka"},lCond: ListingsExtendedSearchCondition{PriceListings: PriceListings{MaxListPrice: 5000000, MinListPrice: 500000}}}, "http://api.booli.se/listings?offset=0&limit=3&q=nacka&maxListPrice=5000000&minListPrice=500000"},
														{searchCondList{sCond: SearchCondition{Q: "nacka"},lCond: ListingsExtendedSearchCondition{PriceDecrease: true}}, "http://api.booli.se/listings?offset=0&limit=3&q=nacka&priceDecrease=1"}}
																										
var searchConditionSoldTests = []searchMatchSoldPos {{searchCondSold{sCond: SearchCondition{Q: "nacka"},seCond: SoldExtendedSearchCondition{PriceSold: PriceSold{MaxSoldPrice: 5000000, MinSoldPrice: 500000}, MinSoldDate: "20130101", MaxSoldDate: "20131010"}}, "http://api.booli.se/sold?offset=0&limit=3&q=nacka&maxSoldPrice=5000000&minSoldPrice=500000&maxSoldDate=20131010&minSoldDate=20130101"}}
												 
func TestGetResultImpl (t *testing.T) {
	// Test caller id empty
	_, err := GetResultImpl("listings?", "", "xxx", nil)
	if err == nil {
		t.Errorf("Should be missing caller id error!")
	}
	
	// Test key empty
	_, err = GetResultImpl("listings?", "xxx", "", nil)
	if err == nil {
		t.Errorf("Should be key empty error!")
	}
	
	// Test missing vital search conditions
	searchC := SearchCondition{}
	_, err = searchC.getSearchString()
	if err == nil {
		t.Errorf("Should be missing search condition error")
	}
	
	// Test missing vital search condition
	_, err = GetResultImpl("listings?", "xxx", "xxx", &MockHttpGet{UrlMatch: "http://api.booli.se/listings?offset=0&limit=3&q=nacka", TestType: GetCheckSearchCond})
	if err == nil {
		t.Errorf("Error, should alert for missing search conditions!")
	}
	
	// Test	search conditions positive tests
	for i := range searchConditionPositiveTests {
		searchRes, err := searchConditionPositiveTests[i].SearchCond.getSearchString()
		_, err = GetResultImpl("listings?" + searchRes, "xxx", "xxx", &MockHttpGet{UrlMatch: searchConditionPositiveTests[i].UrlMatch, TestType: GetCheckSearchCond})
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}
	
	// Test	search conditions negative tests
	for i := range searchConditionNegativeTests {
		_, err := searchConditionNegativeTests[i].SearchCond.getSearchString()
		if err == nil {
			t.Errorf("The current test should produce error: " + searchConditionNegativeTests[i].ExpectedError)
		}
	}
	
	// Test positive list conditions tests
	for i := range searchConditionListingsTests {
		searchRes, err := searchConditionListingsTests[i].sCondList.sCond.getSearchString()
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		
		extlistRes, err := searchConditionListingsTests[i].sCondList.lCond.getSearchString()		
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		
		_, err = GetResultImpl("listings?" + searchRes + extlistRes, "xxx", "xxx", &MockHttpGet{UrlMatch: searchConditionListingsTests[i].UrlMatch, TestType: GetCheckSearchCond})
		
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}
	
	// Test positive sold conditions tests
	for i := range searchConditionSoldTests {
		searchRes, err := searchConditionSoldTests[i].sCondSold.sCond.getSearchString()
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		
		extSoldRes, err := searchConditionSoldTests[i].sCondSold.seCond.getSearchString()		
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		
		_, err = GetResultImpl("sold?" + searchRes + extSoldRes, "xxx", "xxx", &MockHttpGet{UrlMatch: searchConditionSoldTests[i].UrlMatch, TestType: GetCheckSearchCond})
		
		if err != nil {
			t.Errorf("%s", err.Error())
		}
	}
	
	// Test negative list condition tests
	
	var listExtCondNegMax = ListingsExtendedSearchCondition{PriceListings: PriceListings{MaxListPrice: -5000000, MinListPrice: 500000}}
	_, err = listExtCondNegMax.getSearchString()
	if (err == nil) {
		t.Errorf("The current test should produce error: MaxListPrice can not be negative!")
	}
	
	var listExtCondNegMin = ListingsExtendedSearchCondition{PriceListings: PriceListings{MaxListPrice: 5000000, MinListPrice: -500000}}
	_, err = listExtCondNegMin.getSearchString()
	if (err == nil) {
		t.Errorf("The current test should produce error: MinListPrice can not be negative!")
	}
	
	// Test negative sold condition tests	
	var soldExtCondNegMax = SoldExtendedSearchCondition{PriceSold: PriceSold{MaxSoldPrice: -5000000, MinSoldPrice: 500000}}
	_, err = soldExtCondNegMax.getSearchString()
	if (err == nil) {
		t.Errorf("The current test should produce error: MaxSoldPrice can not be negative!")
	}
	
	var soldExtCondNegMin = SoldExtendedSearchCondition{PriceSold: PriceSold{MaxSoldPrice: 5000000, MinSoldPrice: -500000}}
	_, err = soldExtCondNegMin.getSearchString()
	if (err == nil) {
		t.Errorf("The current test should produce error: MinSoldPrice can not be negative!")
	}
	
	var soldExtCondIncorrectMaxSoldDate = SoldExtendedSearchCondition{PriceSold: PriceSold{MaxSoldPrice: 5000000, MinSoldPrice: 500000}, MinSoldDate: "20130101", MaxSoldDate: "10102013"}
	_, err = soldExtCondIncorrectMaxSoldDate.getSearchString()
	if (err == nil) {
		t.Errorf("The current test should produce error: MaxSoldDate is not in the format 20060102, YYYYMMDD!")
	}
	
    var soldExtCondIncorrectMinSoldDate = SoldExtendedSearchCondition{PriceSold: PriceSold{MaxSoldPrice: 5000000, MinSoldPrice: 500000}, MinSoldDate: "2013-01-01", MaxSoldDate: "20131010"}
	_, err = soldExtCondIncorrectMinSoldDate.getSearchString()
	if (err == nil) {
		t.Errorf("The current test should produce error: MinSoldDate is not in the format 20060102, YYYYMMDD!")
	}
	
	// Test positive search condition area test
	var searchCondAreaQTest = SearchConditionArea{Q: "nacka"}
	outstr, err := searchCondAreaQTest.getSearchString()		
	if err != nil {
		t.Errorf("%s", err.Error())
	}
		
	_, err = GetResultImpl("areas?" + outstr, "xxx", "xxx", &MockHttpGet{UrlMatch: "http://api.booli.se/areas?q=nacka", TestType: GetCheckSearchCond})
		
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	
	var searchCondAreaLatLongTest = SearchConditionArea{LatLong: "59.334979,18.065579"}
	outstr, err = searchCondAreaLatLongTest.getSearchString()		
	if err != nil {
		t.Errorf("%s", err.Error())
	}
		
	_, err = GetResultImpl("areas?" + outstr, "xxx", "xxx", &MockHttpGet{UrlMatch: "http://api.booli.se/areas?lat=59.334979&lng=18.065579", TestType: GetCheckSearchCond})
		
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	
	// Test negative search condition area test
	
	var searchCondAreaNoArgs = SearchConditionArea{}
	outstr, err = searchCondAreaNoArgs.getSearchString()		
	if err == nil {
		t.Errorf("The current test should produce error: Need Q OR LatLong to perform a search, not both!")
	}
	
	var searchCondAreaQandLatLong = SearchConditionArea{Q: "nacka", LatLong: "59.334979,18.065579"}
	outstr, err = searchCondAreaQandLatLong.getSearchString()		
	if err == nil {
		t.Errorf("The current test should produce error: Need Q OR LatLong to perform a search, not both!")
	}
	
	var searchCondAreaLat = SearchConditionArea{LatLong: "590.334979,18.065579"}
	outstr, err = searchCondAreaLat.getSearchString()		
	if err == nil {
		t.Errorf("The current test should produce error: Latitude must be between 90 and -90 and Longitude must be between 180 and -180 and be in the format 1.0,1.0!")
	}
	
	var searchCondAreaLong = SearchConditionArea{LatLong: "59.334979,1845.065579"}
	outstr, err = searchCondAreaLong.getSearchString()		
	if err == nil {
		t.Errorf("The current test should produce error: Latitude must be between 90 and -90 and Longitude must be between 180 and -180 and be in the format 1.0,1.0!")
	}
	
	// Test wrong auth
	_, err = GetResultImpl("listings?offset=0&limit=3&q=nacka", "xxx", "xxx", &MockHttpGet{UrlMatch: "http://api.booli.se/listings?offset=0&limit=3&q=nacka", TestType: GetNonAuth})
	if err == nil {
		t.Errorf("Error, should alert for wrong auth!")
	}
}
