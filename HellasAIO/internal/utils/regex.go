package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

// LRRegex does not work in golang
func LRRegex(lValue, rValue, text string) string {
	regex := regexp.MustCompile(fmt.Sprintf(`(?<=%s)(.*)(?=%s)`, lValue, rValue))
	return regex.FindString(text)
}

func GetAFOrderID(body string) int {
	regex := regexp.MustCompile(`"OrderId":(\d+)`)
	orderId := regex.FindStringSubmatch(body)
	if len(orderId) > 1 {
		orderIdInt, _ := strconv.Atoi(orderId[1])
		return orderIdInt
	} else {
		return -1
	}
}

func GetBuzzProductId(body string) string {
	regex := regexp.MustCompile(`athlitika-papoutsia..(\d{4})`)
	productId := regex.FindStringSubmatch(body)

	return productId[1]
}
