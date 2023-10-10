package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ApiCallResult struct {
	ResponseBody    string
	RequestEndpoint string
}

func main() {

	for _, cep := range os.Args[1:] {
		viaCepChannel := make(chan ApiCallResult)
		apiCepChannel := make(chan ApiCallResult)

		go func(cep string) {
			requestUrl := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
			viaCepChannel <- ApiCall(requestUrl)

		}(cep)

		go func(cep string) {
			requestUrl := fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cep)
			apiCepChannel <- ApiCall(requestUrl)

		}(cep)

		select {
		case viaCepResponse := <-viaCepChannel:
			fmt.Printf("Request endpoint: %v\n\nRsponse:\n\n %v\n",
				viaCepResponse.RequestEndpoint,
				viaCepResponse.ResponseBody)
		case apiCepResponse := <-apiCepChannel:
			fmt.Printf("Request endpoint: %v\n\nRsponse:\n\n %v\n",
				apiCepResponse.RequestEndpoint,
				apiCepResponse.ResponseBody)
		case <-time.After(time.Second * 1):
			println("Timeout")
		}
	}
}

func ApiCall(endpoint string) ApiCallResult {
	resp, err := http.Get(endpoint)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	return ApiCallResult{
		ResponseBody:    string(body),
		RequestEndpoint: endpoint,
	}
}
