package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {

	returnViaCep := make(chan map[string]string)
	returnBrasil := make(chan map[string]string)

	go getCepByViaCep(returnViaCep, "14412402")
	go getCepByBrasilApi(returnBrasil, "14412402")

	select {
	case response := <-returnViaCep:
		fmt.Println("Resposta mais rápida do VIA CEP, response:")
		fmt.Println(response)
	case response := <-returnBrasil:
		fmt.Println("Resposta mais rápida do Brasil API, response:")
		fmt.Println(response)

	case <-time.After(time.Second * 5):
		println("timeout")

	}

}

func getCepByViaCep(ch chan map[string]string, cep string) error {
	resp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		log.Println(err)
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	response := map[string]string{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	ch <- response

	return nil
}

func getCepByBrasilApi(ch chan map[string]string, cep string) error {
	resp, err := http.Get(fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep))
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	response := map[string]string{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	ch <- response

	return nil
}
