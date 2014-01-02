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
			"strings"
			"fmt"
		)
			
const 	(
			BooliHttp = "http://api.booli.se/" // listings? is now prefix
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

type IncorrectArgumentError struct {
	ErrString string
}

func (e *IncorrectArgumentError) Error() string {
	return fmt.Sprintf("Incorrect argument with error message: %s", e.ErrString)
}

// Data received from Booli json is parsed into the following struct
type Result struct {
    TotalCount  int64
    Count int64
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

type SearchConditionId struct {
	Prefix string
	Id string
}

func (s *SearchConditionId) getSearchString() (searchString string, err error) {

	
	
	
	return "", nil // Todo
}

type SearchConditionArea struct {
	Q string
	LatLong string
}

func (s *SearchConditionArea) getSearchString() (searchString string, err error) {
	
	if s.Q == "" && s.LatLong == "" {
	return "", &MissingArgumentError{ErrString: "Need Q OR LatLong to perform a search!"}
	}
	
	if s.Q != "" && s.LatLong != "" {
	return "", &MissingArgumentError{ErrString: "Need Q OR LatLong to perform a search, not both!"}
	}
	
	if s.LatLong != "" {
		val, err := formatCheck(s.LatLong,2,"Latitude must be between 90 and -90 and Longitude must be between 180 and -180 and be in the format 1.0,1.0!","LatLongCheck")
		if err != nil {
			return "", err
		}
		split := strings.Split(val, ",")
		searchString += "lat=" + split[0] + "&lng=" + split[1]
	}
		
	if s.Q != "" {
		searchString += "q=" + s.Q
	} 
	
	return searchString, nil
}

type ListingsExtendedSearchCondition struct {
	PriceListings
	PriceDecrease bool
}

func (s *ListingsExtendedSearchCondition) getSearchString() (searchString string, err error) {
	
	if s.PriceListings.MaxListPrice != 0 {
		if s.PriceListings.MaxListPrice < 0 {
			return "", &IncorrectArgumentError{ErrString: "MaxListPrice can not be negative!" }
		}
		searchString += "&maxListPrice=" +  strconv.FormatInt(int64(s.PriceListings.MaxListPrice),10)
	}
	
	if s.PriceListings.MinListPrice != 0 {
		if s.PriceListings.MinListPrice < 0 {
			return "", &IncorrectArgumentError{ErrString: "MinListPrce can not be negative!" }
		}
		searchString += "&minListPrice=" +  strconv.FormatInt(int64(s.PriceListings.MinListPrice),10)
	}
	
	if s.PriceDecrease == true {
		searchString += "&priceDecrease=1"
	}

	return searchString, nil
}

type SoldExtendedSearchCondition struct {
	PriceSold
	MinSoldDate string
	MaxSoldDate string
}

func (s *SoldExtendedSearchCondition) getSearchString() (searchString string, err error) {

	if s.PriceSold.MaxSoldPrice != 0 {
		if s.PriceSold.MaxSoldPrice < 0 {
			return "", &IncorrectArgumentError{ErrString: "MaxSoldPrice can not be negative!" }
		}
		searchString += "&maxSoldPrice=" +  strconv.FormatInt(int64(s.PriceSold.MaxSoldPrice),10)
	}
	
	if s.PriceSold.MinSoldPrice != 0 {
		if s.PriceSold.MinSoldPrice < 0 {
			return "", &IncorrectArgumentError{ErrString: "MinSoldPrice can not be negative!" }
		}
		searchString += "&minSoldPrice=" +  strconv.FormatInt(int64(s.PriceSold.MinSoldPrice),10)
	}
	
	if s.MaxSoldDate != "" {
		_, err := time.Parse("20060102", s.MaxSoldDate)
		
		if err != nil {
			return "", &IncorrectArgumentError{ErrString: "MaxSoldDate is not in the format 20060102, YYYYMMDD!" }	
		}
		searchString += "&maxSoldDate=" + s.MaxSoldDate
	}
		
	if s.MinSoldDate != "" {
		_, err := time.Parse("20060102", s.MinSoldDate)
		
		if err != nil {
			return "", &IncorrectArgumentError{ErrString: "MinSoldDate is not in the format 20060102, YYYYMMDD!" }	
		}
		searchString += "&minSoldDate=" + s.MinSoldDate
	}
	
	return searchString, nil
}

type SearchCondition struct {
	Q string
	Center string
	Dim string
	Bbox string
	AreaId string
	Rooms
	MaxRent int
	LivingArea
	MinPlotArea int
	MaxPlotArea int
	ObjectType string
	MinPublished string
	MaxPublished string
	MinConstructionYear string
	MaxConstructionYear string
	minSqmPrice int
	maxSqmPrice int
	Limit int
	Offset int
}

type PriceListings struct {
	MaxListPrice int
	MinListPrice int
}

type PriceSold struct {
	MaxSoldPrice int
	MinSoldPrice int
}

type Rooms struct {
	MaxRooms int
	MinRooms int
}

type LivingArea struct {
	MaxLivingArea int
	MinLivingArea int
}

type IHttpGet interface {
	Get(url string) (r *http.Response, err error)
}

func (s *SearchCondition) getSearchString() (searchString string, err error) {

	if s.Q == "" && s.Center == "" && s.AreaId == "" {
	return "", &MissingArgumentError{ErrString: "Need Q, Center or AreaId to perform a search!"}
	}

	if s.Offset != 0 {
		searchString += "offset=" + strconv.FormatInt(int64(s.Offset),10)
	} else {
		searchString += "offset=0"
	}
	
	if s.Limit != 0 {
		if s.Limit < 0 {
			return "", &IncorrectArgumentError{ErrString: "Limit can not be negative!" }
		}
		searchString += "&limit=" + strconv.FormatInt(int64(s.Limit),10)
	} else {
		searchString += "&limit=3"
	}
	
	if s.MaxPublished != "" {
		_, err := time.Parse("20060102", s.MaxPublished)
		
		if err != nil {
			return "", &IncorrectArgumentError{ErrString: "MaxCreated is not in the format 20060102, YYYYMMDD!" }	
		}
		searchString += "&maxPublished=" + s.MaxPublished
	}
	
	if s.MinPublished != "" {
		_, err := time.Parse("20060102", s.MinPublished)
		
		if err != nil {
			return "", &IncorrectArgumentError{ErrString: "MinCreated is not in the format 20060102, YYYYMMDD!" }	
		}
		
		searchString += "&minPublished=" + s.MinPublished
	}
	
	if s.MaxConstructionYear != "" {
		_, err := time.Parse("20060102", s.MaxConstructionYear)
		
		if err != nil {
			return "", &IncorrectArgumentError{ErrString: "MaxConstructionYear is not in the format 20060102, YYYYMMDD!"}
		}
		
		searchString += "&maxConstructionYear=" + s.MaxConstructionYear
	}
	
	if s.MinConstructionYear != "" {
		_, err := time.Parse("20060102", s.MinConstructionYear)
		
		if err != nil {
			return "", &IncorrectArgumentError{ErrString: "MinConstructionYear is not in the format 20060102, YYYYMMDD!"}
		}
		
		searchString += "&minConstructionYear=" + s.MinConstructionYear
	}
		
	if s.ObjectType != "" {
		// Check for bad objecttype
		searchString += "&objectType=" + s.ObjectType
	}
	
	if s.MaxPlotArea != 0 {
		if s.MaxPlotArea < 0 {
			return "", &IncorrectArgumentError{ErrString: "MaxPlotArea can not be negative!" }
		}
		searchString += "&maxPlotArea=" +  strconv.FormatInt(int64(s.MaxPlotArea),10)
	}
	
	if s.MinPlotArea != 0 {
		if s.MinPlotArea < 0 {
			return "", &IncorrectArgumentError{ErrString: "MinPlotArea can not be negative!" }
		}
		searchString += "&minPlotArea=" +  strconv.FormatInt(int64(s.MinPlotArea),10)
	}
		
	if s.LivingArea.MaxLivingArea != 0 {
		if s.LivingArea.MaxLivingArea < 0 {
			return "", &IncorrectArgumentError{ErrString: "MaxLivingArea can not be negative!" }
		}
		searchString += "&maxLivingArea=" +  strconv.FormatInt(int64(s.LivingArea.MaxLivingArea),10)
	}	
	
	if s.LivingArea.MinLivingArea != 0 {
		if s.LivingArea.MinLivingArea < 0 {
			return "", &IncorrectArgumentError{ErrString: "MinLivingArea can not be negative!" }
		}
		searchString += "&minLivingArea=" +  strconv.FormatInt(int64(s.LivingArea.MinLivingArea),10)
	}
	
	if s.MaxRent != 0 {
		if s.MaxRent < 0 {
			return "", &IncorrectArgumentError{ErrString: "MaxRent can not be negative!" }
		}
		searchString += "&maxRent=" +  strconv.FormatInt(int64(s.MaxRent),10)
	}
	
	if s.Rooms.MaxRooms != 0 {
		if s.Rooms.MaxRooms < 0 {
			return "", &IncorrectArgumentError{ErrString: "MaxRooms can not be negative!" }
		}
		searchString += "&maxRooms=" +  strconv.FormatInt(int64(s.Rooms.MaxRooms),10)
	}
	
	if s.Rooms.MinRooms != 0 {
		if s.Rooms.MinRooms < 0 {
			return "", &IncorrectArgumentError{ErrString: "MinRooms can not be negative!" }
		}
		searchString += "&minRooms=" +  strconv.FormatInt(int64(s.Rooms.MinRooms),10)
	}
		
	if s.AreaId != "" {
		val, err := formatCheck(s.AreaId,0,"AreaId must be one or more positive numbers in the format 66,78...!","GreaterThenZeroCheck")
		if err != nil {
			return "", err
		}
		searchString += "&areaId=" + val
	}
	
	if s.Bbox != "" {
		val, err := formatCheck(s.Bbox,4,"Bbox must be two lat-long pairs, on the form lat_lo,long_lo,lat_hi,long_hi where lo is south west and hi is north east!","BboxCheck")
		if err != nil {
			return "", err
		}
		searchString += "&bbox=" + val
	}
		
	if s.Dim != "" {
		if s.Center == "" {
		return "", &MissingArgumentError{ErrString: "Need Center if Dim is used!"}
		}
		val, err := formatCheck(s.Dim,2,"Dim must be two positive numbers in the format 1,1!","GreaterThenZeroCheck")
		if err != nil {
			return "", err
		}
		searchString += "&dim=" + val
	}
	
	if s.Center != "" {
		if s.Dim == "" {
		return "", &MissingArgumentError{ErrString: "Need Dim if Center is used!"}
		}
		val, err := formatCheck(s.Center,2,"Latitude must be between 90 and -90 and Longitude must be between 180 and -180 and be in the format 1.0,1.0!","LatLongCheck")
		if err != nil {
			return "", err
		}
		searchString += "&center=" + val
	} 
	
	if s.Q != "" {
		searchString += "&q=" + s.Q
	} 
	
	return searchString, nil
}

// Returns a result for Result Listings from Booli or a empty result and a error if a problem was encountered
func GetResultListings(searchCond SearchCondition, listCond ListingsExtendedSearchCondition, callerId string, key string) (booliRes Result, err error) {

	scond, err := searchCond.getSearchString()
	if err != nil {
		return booliRes, err
	}
	
	lcond, err := listCond.getSearchString()
	if err != nil {
		return booliRes, err
	}
	
	return GetResultImpl("listings?" + scond + lcond, callerId, key, &http.Client{})
}

// Returns a result for Sold Listings from Booli or a empty result and a error if a problem was encountered
func GetResultSold(prefix string, searchCond SearchCondition, soldCond SoldExtendedSearchCondition, callerId string, key string) (booliRes Result, err error) {

	scond, err := searchCond.getSearchString()
	if err != nil {
		return booliRes, err
	}
	
	soldcond, err := soldCond.getSearchString()
	if err != nil {
		return booliRes, err
	}
	
	return GetResultImpl("sold?" + scond + soldcond, callerId, key, &http.Client{})
}

// Returns a result for SearchArea Listings from Booli or a empty result and a error if a problem was encountered
func GetResultSearchArea(prefix string, searchCond SearchConditionArea, callerId string, key string) (booliRes Result, err error) {

	scond, err := searchCond.getSearchString()
	if err != nil {
		return booliRes, err
	}
	
	return GetResultImpl("areas?" + scond, callerId, key, &http.Client{})
}

// Add a implementation function to be able to feed a custom Http.get function for unit testing
func GetResultImpl(cond string, callerId string, key string, httpGet IHttpGet) (booliRes Result, err error) {
	if callerId == "" {
		return booliRes, &MissingArgumentError{ErrString: "Caller Id empty!"}
	}
	
	if key == "" {
		return booliRes, &MissingArgumentError{ErrString: "Key empty!"}
	}

	searchStr, err := searchStr(cond, callerId, key)
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

func searchStr(cond string, callerId string, key string) (outstr string, err error) {

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

func formatCheck(instr string, length int, errorMsg string, errorType string) (outstr string, err error) {
	split := strings.Split(instr, ",")
	
	if len(split) != length && length != 0 {
		return "", &IncorrectArgumentError{ErrString: errorMsg }
	}
	for i, v := range split {
		val, err := strconv.ParseFloat(v,64)
		if err != nil {
			return "", &IncorrectArgumentError{ErrString: errorMsg }
		}
		switch errorType {
		case "GreaterThenZeroCheck":
		if val < 0 {
			return "", &IncorrectArgumentError{ErrString: errorMsg }
		}
		case "LatLongCheck":
			if i == 0 {
				if checkLatOutOfBound(val) {
					return "", &IncorrectArgumentError{ErrString: errorMsg }
				}
		} else {
				if checkLongOutOfBound(val) {
					return "", &IncorrectArgumentError{ErrString: errorMsg }
				}
		}
		case "BboxCheck":
			switch i {
				case 0:
					if checkLatOutOfBound(val) {
						return "", &IncorrectArgumentError{ErrString: errorMsg }
					}
				case 1:
					if checkLongOutOfBound(val) {
						return "", &IncorrectArgumentError{ErrString: errorMsg }
					}
				case 2:
					if checkLatOutOfBound(val) {
						return "", &IncorrectArgumentError{ErrString: errorMsg }
					}
				case 3:
					if checkLongOutOfBound(val) {
						return "", &IncorrectArgumentError{ErrString: errorMsg }
					}
			}
		}
	}
	return instr, nil
}

func checkLatOutOfBound(val float64) bool {
	if val > 90 || val < -90 {
		return true
	}
	return false
}

func checkLongOutOfBound(val float64) bool {
	if val > 180 || val < -180 {
		return true
	}
	return false
}