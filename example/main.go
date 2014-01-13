package main

import  (
		 //"github.com/Kenshins/BooliGo"
		 "BooliGo"
		 "fmt"
		 "strconv"
		 )
		 
const (
		BooliId = "your-booli-id" // Supplied from booli, http://www.booli.se/api/key
		BooliKey = "P9rhkeJvKOKGijvXZ1npRXVKK2kHPmXpdIN3etZS" // Supplied from booli, http://www.booli.se/api/key
	)

func main() {
	booliResL, err := booli.GetResultListings(booli.SearchCondition{Q: "nacka", LivingArea: booli.LivingArea{MinLivingArea: 65}, Limit: 5},booli.ListingsExtendedSearchCondition{PriceListings: booli.PriceListings{MaxListPrice: 3000000, MinListPrice: 300000}}, BooliId, BooliKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Totalcount : " + strconv.FormatInt(booliResL.TotalCount,10))
	fmt.Println("Count : " + strconv.FormatInt(booliResL.Count,10))
	
	for _, v := range booliResL.Listings {
		fmt.Println("")
		fmt.Println("BooliId : " + strconv.FormatInt(v.BooliId,10))
		fmt.Println("Listprice : " + strconv.FormatInt(v.ListPrice,10))
		fmt.Println("Published : " + v.Published)
	}
	fmt.Println("\n")
	
	booliResS, err := booli.GetResultSold(booli.SearchCondition{Q: "nacka", LivingArea: booli.LivingArea{MinLivingArea: 65}, Limit: 5},booli.SoldExtendedSearchCondition{PriceSold: booli.PriceSold{MaxSoldPrice: 3000000, MinSoldPrice: 300000}}, BooliId, BooliKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Totalcount : " + strconv.FormatInt(booliResS.TotalCount,10))
	fmt.Println("Count : " + strconv.FormatInt(booliResS.Count,10))
	
	for _, v := range booliResS.Sold {
		fmt.Println("")
		fmt.Println("BooliId : " + strconv.FormatInt(v.BooliId,10))
		fmt.Println("Soldprice : " + strconv.FormatInt(v.SoldPrice,10))
		fmt.Println("Published : " + v.Published)
	}

	booliResA, err := booli.GetResultSearchArea(booli.SearchConditionArea{Q: "nacka"}, BooliId, BooliKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Totalcount : " + strconv.FormatInt(booliResA.TotalCount,10))
	fmt.Println("Count : " + strconv.FormatInt(booliResA.Count,10))
	
	for _, v := range booliResA.Areas {
		fmt.Println("")
		fmt.Println("BooliId : " + strconv.FormatInt(v.BooliId,10))
		fmt.Println("Name : " + v.Name)
		fmt.Println("Full name : " + v.FullName)
	}
	
	booliResSoldId, err := booli.GetResultSoldId(1386756, BooliId, BooliKey)
	if err != nil {
		fmt.Println(err)
	}
	
	for _, v := range booliResSoldId.Sold {
		fmt.Println("")
		fmt.Println("BooliId : " + strconv.FormatInt(v.BooliId,10))
		fmt.Println("Soldprice : " + strconv.FormatInt(v.SoldPrice,10))
		fmt.Println("Published : " + v.Published)
	}
	
	booliResListId, err := booli.GetResultListId(1574414, BooliId, BooliKey)
	if err != nil {
		fmt.Println(err)
	}
	
	for _, v := range booliResListId.Listings {
		fmt.Println("")
		fmt.Println("BooliId : " + strconv.FormatInt(v.BooliId,10))
		fmt.Println("Listprice : " + strconv.FormatInt(v.ListPrice,10))
		fmt.Println("Published : " + v.Published)
	}
}
