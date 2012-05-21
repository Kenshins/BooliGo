package booli

import ( 	"net/http"
			"time"
			"errors"
			"crypto/sha1"
			"crypto/rand"
			"encoding/json"
			"io"
			"io/ioutil"
			"strconv"
			"fmt"
		)
			
const 	(
			BooliHttp = "http://api.booli.se/listings?"
			Abc = "abcdefghijklmnopqrstuvxyzABCDEFGHIJKLMNOPQRSTXYZ1234567890"
			
		)
		
// Error types

type AuthError struct {
	ErrString string
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("Received http response code 403 with error message: %s", e.ErrString)
}

type MissingArgumentError struct {
	ErrString string
}

func (e *MissingArgumentError) Error() string {
	return fmt.Sprintf("Missing argument with error message: %s", e.ErrString)
}

// Data received from Booli json is parsed into the following struct
type Result struct {
    TotalCount  int
    Count int
	Listings []HouseList
}

type HouseList struct {
	BooliId int64
	ListPrice int64
	Published string
	ListPriceChangeDate string
	Objecttype string
	Location LocationObject
	Source SourceObject
	Rooms float64
    LivingArea float64
    Rent int64
    Floor int64
    IsNewConstruction int64
    Url string
}

type LocationObject struct {
	Region RegionObject
	Address AddressObject
	NamedAreas []string
	Position PositionObject
}

type RegionObject struct {
    MunicipalityName string
    CountyName string
}

type AddressObject struct {
	City string
	StreetAddress string
}

type PositionObject struct {
	Latitude float64
	Longitude float64
}

type SourceObject struct {
	Name string
	Url string
	Type string
}

// Searchconditions for Booli search.
type SearchCondition struct {
	Q string
	Center string
	Dim string
	Bbox string
	AreaId string
	MinPrice int
	MaxPrice int
	MinRooms int
	MaxRooms int
	MaxRent int
	MinLivingArea int
	MaxLivingArea int
	MinPlotArea int
	MaxPlotArea int
	ObjectType string
	MinCreated string
	MaxCreated string
	Limit int
	Offset int
}

type IHttpGet interface {
	Get(url string) (r *http.Response, err error)
}

func (s *SearchCondition) getSearchString() (searchString string, err error) {

	if s.Q == "" && s.Center == "" && s.AreaId == "" {
	return "", errors.New("Need Q, Center or AreaId to perform a search!")	
	}

	if s.Offset != 0 {
		searchString += "offset=" + strconv.FormatInt(int64(s.Offset),10)
	} else {
		searchString += "offset=0"
	}
	
	if s.Limit != 0 {
		searchString += "&limit=" + strconv.FormatInt(int64(s.Limit),10)
	} else {
		searchString += "&limit=3"
	}
	
	if s.MaxCreated != "" {
		// Todo: Check for bad date
		searchString += "&maxCreated=" + s.MaxCreated
	}
	
	if s.MinCreated != "" {
		// Todo: Check for bad date
		searchString += "&minCreated=" + s.MinCreated
	}
		
	if s.ObjectType != "" {
		// Check for bad objecttype
		searchString += "&objectType=" + s.ObjectType
	}
	
	if s.MaxPlotArea != 0 {
		searchString += "&maxPlotArea=" +  strconv.FormatInt(int64(s.MaxPlotArea),10)
	}
	
	if s.MinPlotArea != 0 {
		searchString += "&minPlotArea=" +  strconv.FormatInt(int64(s.MinPlotArea),10)
	}
		
	if s.MaxLivingArea != 0 {
		searchString += "&maxLivingArea=" +  strconv.FormatInt(int64(s.MaxLivingArea),10)
	}	
	
	if s.MinLivingArea != 0 {
		searchString += "&minLivingArea=" +  strconv.FormatInt(int64(s.MinLivingArea),10)
	}
	
	if s.MaxRent != 0 {
		searchString += "&maxRent=" +  strconv.FormatInt(int64(s.MaxRent),10)
	}
	
	if s.MaxRooms != 0 {
		searchString += "&maxRooms=" +  strconv.FormatInt(int64(s.MaxRooms),10)
	}
	
	if s.MinRooms != 0 {
		searchString += "&minRooms=" +  strconv.FormatInt(int64(s.MinRooms),10)
	}	
	
	if s.MaxPrice != 0 {
		searchString += "&maxPrice=" +  strconv.FormatInt(int64(s.MaxPrice),10)
	}

	if s.MinPrice != 0 {
		searchString += "&minPrice=" +  strconv.FormatInt(int64(s.MinPrice),10)
	}
		
	if s.AreaId != "" {
		// Todo: Check input to conform to 33,44...
		searchString += "&areaId=" + s.AreaId
	}
	
	if s.Bbox != "" {
		// Todo: Check input to conform to 1,1,1,1
		searchString += "&bbox=" + s.Bbox
	}
		
	if s.Dim != "" {
		// Todo: Check input to conform to 1,1
		searchString += "&dim=" + s.Dim	
	}
	
	if s.Center != "" {
		// Todo: Check input to conform to 1,1
		searchString += "&center=" + s.Center
	} 
	
	if s.Q != "" {
		searchString += "&q=" + s.Q
	} 
	
	return searchString, nil
}

// Returns a result from Booli or a empty result and a error if a problem was encountered
func GetResult(searchCond SearchCondition, callerId string, key string) (booliRes Result, err error) {
	return GetResultImpl(searchCond, callerId, key, &http.Client{})
}

// Add a implementation function to be able to feed a custom Http.get function for unit testing
func GetResultImpl(searchCond SearchCondition, callerId string, key string, httpGet IHttpGet) (booliRes Result, err error) {
	if callerId == "" {
		return booliRes, &MissingArgumentError{ErrString: "Caller Id empty!"}
	}
	
	if key == "" {
		return booliRes, &MissingArgumentError{ErrString: "Key empty!"}
	}

	searchStr, err := searchStr(searchCond, callerId, key)
	if err != nil {
		return booliRes, err
	}
	
	resp, err := httpGet.Get(searchStr)
	if err != nil {
		return booliRes, err
	}
	
	if resp.StatusCode == 403 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return booliRes, errors.New("When reading the body from http response code 403 the following error occurred: " + err.Error())
		}
		return booliRes, &AuthError{ErrString: string(body)}
	}
	
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return booliRes, err
	}
	err = json.Unmarshal(body, &booliRes)
	if err != nil {
		return booliRes, err
	}
	return booliRes, nil
}

func sha1String(instr string) (outStr string) {
	h := sha1.New()
	io.WriteString(h, instr)	
	for _, v := range h.Sum(nil) {
		tmpVal := strconv.FormatInt(int64(v),16)
		if len(tmpVal) == 1 {
			tmpVal = "0" + tmpVal
		}
		outStr = outStr + tmpVal
	}
	return outStr
}

func searchStr(searchCond SearchCondition, callerId string, key string) (outstr string, err error) {
	
	cond, err := searchCond.getSearchString()
	if err != nil {
		return outstr, err
	}
	
	time := strconv.FormatInt(int64(time.Now().Unix()),10)
	
	unique, err := unique()
	if err != nil {
		return outstr, err
	}
	
	hash := sha1String(callerId + time + key + unique)
	outstr = BooliHttp + cond + "&callerId=" + callerId + "&time=" + time + "&unique=" + unique + "&hash=" + hash
	return
}

func unique() (outstr string, err error) {
	randbytes := make([]byte, 16)
	if _, err := rand.Read(randbytes); err == nil {
		for i := 0; i < 16; i++ {
			tmpindex := int(randbytes[i]) % len(Abc)
			outstr += string(Abc[tmpindex])
		}
	} else {
		return outstr, err
	}
	return outstr, err
}