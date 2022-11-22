package api

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func EvalStock(key string) string {
	stockServiceUrl := fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", url.QueryEscape(key))
	log.Println("info : processing", stockServiceUrl)

	response, err := http.Get(stockServiceUrl)
	if err != nil {
		log.Println("error :", err)
		return "Stock service is not available"
	}

	if response.StatusCode == http.StatusOK {
		content, err := csv.NewReader(response.Body).ReadAll()
		if err != nil {
			//handle error
			log.Println("error :", err)
			return "Stock service CSV error"
		}
		//now reading first field of second record.
		//Symbol,Date,Time,Open,High,Low,Close,Volume
		symbol := content[1][0]
		close := content[1][6]
		log.Println("content:", content)
		if close == "N/D" {
			return fmt.Sprintf("%s quote is not available", strings.ToUpper(symbol))
		}
		return fmt.Sprintf("%s quote is $%s per share", strings.ToUpper(symbol), close)
	}

	log.Println("error : response.StatusCode is ", response.StatusCode)
	return "Stock service is not available"
}
