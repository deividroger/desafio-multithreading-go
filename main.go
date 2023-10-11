package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ViaCepResult struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ApiCEPResult struct {
	Status     int    `json:"status"`
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Ok         bool   `json:"ok"`
	Message    string `json:"message"`
	StatusText string `json:"statusText"`
}

func main() {

	for _, cep := range os.Args[1:] {
		viaCepChannel := make(chan ViaCepResult)
		apiCepChannel := make(chan ApiCEPResult)

		viaCEPRequestEndpoint := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
		apiCepRequestEndpoint := fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cep)

		go func() {
			viaCepChannel <- ApiCall[ViaCepResult](viaCEPRequestEndpoint)
		}()

		go func() {
			apiCepChannel <- ApiCall[ApiCEPResult](apiCepRequestEndpoint)
		}()

		select {
		case viaCepResponse := <-viaCepChannel:
			fmt.Printf("VIA CEP Request endpoint: %v\n\nRsponse:\n\n %v\n",
				viaCEPRequestEndpoint, viaCepResponse)
		case apiCepResponse := <-apiCepChannel:
			fmt.Printf("API CEP Request endpoint: %v\n\nRsponse:\n\n %v\n",
				apiCepRequestEndpoint,
				apiCepResponse)
		case <-time.After(time.Second * 1):
			println("Timeout")
		}
	}
}

func ApiCall[T interface{}](endpoint string) T {
	resp, err := http.Get(endpoint)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var result T
	err = json.Unmarshal(body, &result)

	if err != nil {
		panic(err)
	}

	return result
}
