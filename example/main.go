package main

import  (
		 //"github.com/Kenshins/BooliGo"
		 "../../BooliGo"
		 "fmt"
		 "strconv"
		 )
		 
const (
		BooliId = "bopren" // Supplied from booli, http://www.booli.se/api/key
		BooliKey = "P8rfkeJvKOXgHjvXZ1npRXVGG2kHPmXpd5NZetHS" // Supplied from booli, http://www.booli.se/api/key
	)

func main() {
	booliRes, err := booli.GetResult("listings?", booli.SearchCondition{Q: "nacka", Price: booli.Price{MaxListPrice: 3000000, MinListPrice: 300000}, LivingArea: booli.LivingArea{MinLivingArea: 65}, Limit: 5}, BooliId, BooliKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Totalcount : " + strconv.FormatInt(booliRes.TotalCount,10))
	fmt.Println("Count : " + strconv.FormatInt(booliRes.Count,10))
	
	for _, v := range booliRes.Listings {
		fmt.Println("")
		fmt.Println("BooliId : " + strconv.FormatInt(v.BooliId,10))
		fmt.Println("Listprice : " + strconv.FormatInt(v.ListPrice,10))
		fmt.Println("Published : " + v.Published)
	}
}
