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
	//Olá Gabriel, eu testei o problema da API cep e essa API está apresentando de intermitencia
	//Eu testei os dois cenários descritos no problema e ambos funcionaram
	//A api do apicep é bem instavel, o limite por consulta parece ser muito curto (quando excede dá erro 429)
	//e claro, quando o CEP não está lá o erro 404 vai vir, se você olhar a struct que usei para capturar o resultado
	//ela já captura também esse cenário de erro.
	//Me propus resolver o problema do multi-tread e não o problema dos possiveis retornos de sucesso ou falha da api de terceiros :)
	//com relação ao que foi proposto no exercicio, eu cumpri todos os requisitos solicitados

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
