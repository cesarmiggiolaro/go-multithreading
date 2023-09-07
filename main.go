package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Cep struct {
	Cep     string
	State   string
	City    string
	Address string
	Server  string
}

func main() {

	var cep string = "80320-040"
	var server1 string = "https://viacep.com.br/ws/" + cep + "/json/"
	var server2 string = "https://cdn.apicep.com/file/apicep/" + cep + ".json"

	print(server1)

	c1 := make(chan Cep)
	c2 := make(chan Cep)

	go func() {
		result, err := GetCepServer(1, server1)
		if err != nil {
			fmt.Println(err)
			return
		}
		c1 <- *result
	}()

	go func() {
		result2, err := GetCepServer(2, server2)
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

func GetCepServer(serverId int, host string) (*Cep, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", host, nil)

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

	var data map[string]interface{}
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println("Erro ao decodificar JSON:", err)
		return nil, err
	}

	if serverId == 1 {
		return &Cep{
			Cep:     data["cep"].(string),
			State:   data["uf"].(string),
			City:    data["localidade"].(string),
			Address: data["logradouro"].(string),
			Server:  "Server 1",
		}, nil

	}

	return &Cep{
		Cep:     data["code"].(string),
		State:   data["state"].(string),
		City:    data["city"].(string),
		Address: data["address"].(string),
		Server:  "Server 2",
	}, nil
}
