package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type CepServer1 struct {
	Cep     string `json:"code"`
	State   string `json:"state"`
	City    string `json:"city"`
	Address string `json:"address"`
	Server  string
}

type CepServer2 struct {
	Cep     string `json:"cep"`
	State   string `json:"uf"`
	City    string `json:"cidade"`
	Address string `json:"logradouro"`
	Server  string
}

func main() {

	var cep string = "80320-040"

	c1 := make(chan CepServer1)
	c2 := make(chan CepServer2)

	go func() {
		result, err := GetCepServer1(cep)
		if err != nil {
			fmt.Println(err)
			return
		}
		c1 <- *result
	}()

	go func() {
		result2, err := GetCepServer2(cep)
		if err != nil {
			fmt.Println(err)
			return
		}
		c2 <- *result2
	}()

	select {
	case servidor1 := <-c1:
		fmt.Println(servidor1)
	case servidor2 := <-c2:
		fmt.Println(servidor2)
	case <-time.After(time.Second):
		println("timeout")
	}

}

func GetCepServer2(cep string) (*CepServer2, error) {

	cep = strings.Replace(cep, "-", "", -1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://viacep.com.br/ws/"+cep+"/json/", nil)

	if err != nil {
		fmt.Println("Erro ao criar a solicitação:", err)
		return nil, err
	}

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao fazer a solicitação:", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Erro ao ler o corpo da resposta:", err)
		return nil, err
	}

	var data2 CepServer2
	err = json.Unmarshal(body, &data2)
	if err != nil {
		fmt.Println("Erro ao fazer o parse da resposta:", err)
	}

	data2.Server = "Server 2"

	return &data2, nil
}

func GetCepServer1(cep string) (*CepServer1, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://cdn.apicep.com/file/apicep/"+cep+".json", nil)

	if err != nil {
		fmt.Println("Erro ao criar a solicitação:", err)
		return nil, err
	}

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao fazer a solicitação:", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Erro ao ler o corpo da resposta:", err)
		return nil, err
	}

	var data1 CepServer1
	err = json.Unmarshal(body, &data1)
	if err != nil {
		fmt.Println("Erro ao fazer o parse da resposta:", err)
	}

	data1.Server = "Server 1"

	return &data1, nil
}
