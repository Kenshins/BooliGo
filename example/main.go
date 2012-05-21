package main

import  (
		 "github.com/Kenshins/BooliGo"
		 "fmt"
		 "strconv"
		 )
		 
const (
		BooliId = "yourBooliId" // Supplied from booli, http://www.booli.se/api/key
		BooliKey = "P8rhkeJzKOXgHj3XZ1npRXVQG2kHPmXpd5NZetKJ" // Supplied from booli, http://www.booli.se/api/key
	)

func main() {
	booliRes, err := booli.GetResult(booli.SearchCondition{Q: "nacka", MaxPrice: 3000000, MinPrice: 300000, MinLivingArea: 65, Limit: 5}, BooliId, BooliKey)
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
