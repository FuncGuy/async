package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type PizzaOrder struct {
	Pizza, Store, Price string
}

func main() {
	var pizza = flag.String("pizza", "", "Pizza to order")
	var store = flag.String("store", "", "Name of the Pizza Store")
	var price = flag.String("price", "", "Price")

	flag.Parse()

	order := PizzaOrder{*pizza, *store, *price}

	body, _ := json.Marshal(order)

	start := time.Now()

	orderChan := make(chan *http.Response)
	paymentChan := make(chan *http.Response)
	storeChan := make(chan *http.Response)

	SendPostAsync("http://localhost:8081", body, orderChan)

	SendPostAsync("http://localhost:8082", body, paymentChan)

	SendPostAsync("http://localhost:8083", body, storeChan)

	orderResponse := <-orderChan
	defer orderResponse.Body.Close()

	bytes, _ := ioutil.ReadAll(orderResponse.Body)
	fmt.Println(string(bytes))

	paymentResponse := <-paymentChan
	defer paymentResponse.Body.Close()
	bytes, _ = ioutil.ReadAll(paymentResponse.Body)
	fmt.Println(string(bytes))

	storeResponse := <-storeChan
	defer storeResponse.Body.Close()
	bytes, _ = ioutil.ReadAll(storeResponse.Body)
	fmt.Println(string(bytes))

	end := time.Now()

	fmt.Printf("Order processed after %v seconds\n", end.Sub(start).Seconds())

}
func SendPostAsync(url string, body []byte, rc chan *http.Response) {
	response, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	rc <- response
}

func SendPostRequest(url string, body []byte) *http.Response {
	response, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	return response
}
