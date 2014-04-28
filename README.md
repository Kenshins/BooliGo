BooliGo
=======

A wrapper written in Go to connect to the Booli.se API.

Installation
=======

go get github.com/Kenshins/BooliGo

Example
=======
<pre>
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
</pre>


License
=======
This Booli Go wrapper is released under the MIT License (MIT).

Copyright (c) 2014 Martin Kleberger

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
