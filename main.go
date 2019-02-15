package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Details struct {
	Store             string `json:"store"`
	ProductCode       string `json:"productcode"`
	Description       string `json:"description"`
	ProductExpiryDate string `json:"expiry"`
	Price             int    `json:"price"`
}

var data []Details
var expiring []map[string]interface{}
var expired []map[string]interface{}
var active []map[string]interface{}
var new map[string][]int
var pos int
var sum int
var high int
var low int
var priceCount int
var average int
var collec map[string][]int

func main() {

	r := gin.Default()
	r.POST("/add", addProduct)
	r.GET("/admin", adminView)
	r.PUT("/change", changePrice)
	r.GET("/find", findProduct)
	r.GET("/view", viewProduct)
	r.Run(":8080") // listen and serve on 0.0.0.0:8080

}

//This function Views all the Products Categorized on Their expiry
//Parameters Should be provided as Post Form data in Postman
/*
Parameters:
store string
*/
func adminView(c *gin.Context) {

	cDate := time.Now().Format("2006/01/02")

	Year, month, day := getDate(cDate)

	for i := 0; i < len(data); i++ {
		if data[i].Store == c.PostForm("store") {
			pYear, pmonth, pday := getDate(data[i].ProductExpiryDate)
			start := time.Date(pYear, pmonth, pday, 0, 0, 0, 0, time.UTC)
			end := time.Date(Year, month, day, 0, 0, 0, 0, time.UTC)
			difference := start.Sub(end).String()
			difftostr := strings.Split(difference, "h")
			hours, _ := strconv.Atoi(difftostr[0])
			if hours == 0 || hours < 0 {
				expired = append(expired, map[string]interface{}{"Info": data[i],
					"status": "Expired",
				})

			} else if hours > 0 && hours < 168 {
				expiring = append(expiring, map[string]interface{}{"Info": data[i],
					"status": "Expiring",
				})
			} else {
				active = append(active, map[string]interface{}{"Info": data[i],
					"status": "Still to Expire",
				})

			}
		}
	}
	c.JSON(200, gin.H{
		"expired":  expired,
		"active":   active,
		"expiring": expiring,
	})
}

//This function finds and give details of a product Based on Unique ProductId
//Parameters Should be provided as Query String in Postman
/*
Parameters:
productcode string
*/
func findProduct(c *gin.Context) {
	code := c.Query("productcode")
	for i := 0; i < len(data); i++ {
		if data[i].ProductCode == code {
			c.JSON(200, gin.H{
				"Result": data[i],
			})
		}
	}
}

//This function adds a Product with a Unique Product id
//Parameters Should be provided as Form data in Postman
/*
Parameters:
store string
productcode string
description string
expiry (2019/02/15) in this format
price int
*/
func addProduct(c *gin.Context) {
	price, _ := strconv.Atoi(c.PostForm("price"))
	for i := 0; i < len(data); i++ {
		if data[i].ProductCode == c.PostForm("productcode") {
			c.JSON(200, gin.H{
				"Error": "Pls Enter A Unique Product Code",
			})
			return
		}
	}
	data = append(data, Details{
		Store:             c.PostForm("store"),
		ProductCode:       c.PostForm("productcode"),
		Description:       c.PostForm("description"),
		ProductExpiryDate: c.PostForm("expiry"),
		Price:             price,
	})
	c.JSON(200, gin.H{
		"Success": "Response Recorded!",
	})

}

//This function is used to enter a new price for an existing product
//Parameters should be provided as Form data through Postman
/*
Parameters:
productcode string
*/
func changePrice(c *gin.Context) {
	newprice, _ := strconv.Atoi(c.PostForm("newprice"))
	code := c.PostForm("productcode")

	for i := 0; i < len(data); i++ {
		if data[i].ProductCode == code {
			pos = i
			sum = data[i].Price
			high = data[i].Price
			low = data[i].Price
			sum = newprice + sum
			priceCount = priceCount + 1

			if !(high > newprice) {
				high = newprice
			}
			if !(low < newprice) {
				low = newprice
			}
			new = map[string][]int{
				code: append(new[code], newprice),
			}
			collec = new

		}
	}
	average = (sum / (priceCount + 1))
	c.JSON(200, gin.H{
		"Info":          data[pos],
		"Average Price": average,
		"High Price":    high,
		"Low Price":     low,
	})

}

//This function splits the date into year month and days
func getDate(value string) (int, time.Month, int) {
	date := strings.Split(value, "/")
	Year, _ := strconv.Atoi(date[0])
	month, _ := strconv.Atoi(date[1])
	Month := time.Month(month)
	day, _ := strconv.Atoi(date[2])
	return Year, Month, day
}

//This function is used to view an existing product
//Parameters should be provided as Query data through Postman
/*
Parameters:
productcode string
*/
func viewProduct(c *gin.Context) {
	code := c.Query("productcode")
	for i := 0; i < len(data); i++ {
		if data[i].ProductCode == code {
			val := idealsort(code)
			c.JSON(200, gin.H{
				"Info":             data[i],
				"Average Price":    average,
				"Lowest Price":     low,
				"Highest Price":    high,
				"Ideal Price":      val,
				"Prices Collected": priceCount,
			})
			return
		}
	}
}

func idealsort(code string) float32 {

	for k := 1; k < len(collec[code]); k++ {
		if err := collec[code][k-1] != 0; collec[code][k-1] > collec[code][k] && err {
			temp := collec[code][k]
			collec[code][k] = collec[code][k-1]
			collec[code][k-1] = temp
		}
	}
	var total int
	for i := 2; i < len(collec[code])-2; i++ {
		total = collec[code][i] + total
	}
	ideal := float32(total)
	return ideal * 0.2
}
