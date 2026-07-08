package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kauandias747474-hue/titandb/internal/rpc"
)

func main() {
	serverAddr := flag.String("server", ":8001", "endereço do servidor TitanDB")
	cmd := flag.String("cmd", "", "comando: PUT, GET, DELETE, 3SUM")
	key := flag.String("key", "", "chave (para PUT/GET/DELETE)")
	val := flag.String("val", "", "valor (para PUT)")
	flag.Parse()

	if *cmd == "" {
		fmt.Println("Erro: comando é obrigatório")
		os.Exit(1)
	}

	client := rpc.NewClient()
	err := client.Connect(*serverAddr)
	if err != nil {
		fmt.Printf("Erro de conexão: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	switch *cmd {
	case "PUT":
		if *key == "" || *val == "" {
			fmt.Println("Erro: PUT requer --key e --val")
			os.Exit(1)
		}
		status, err := client.Put(*key, *val)
		if err != nil {
			fmt.Printf("Erro: %v\n", err)
		} else {
			fmt.Printf("Sucesso (status %d)\n", status)
		}

	case "GET":
		if *key == "" {
			fmt.Println("Erro: GET requer --key")
			os.Exit(1)
		}
		value, status, err := client.Get(*key)
		if err != nil {
			fmt.Printf("Erro: %v\n", err)
		} else if status == rpc.StatusNotFound {
			fmt.Println("Chave não encontrada")
		} else {
			fmt.Printf("Valor: %s\n", value)
		}

	case "DELETE":
		if *key == "" {
			fmt.Println("Erro: DELETE requer --key")
			os.Exit(1)
		}
		status, err := client.Delete(*key)
		if err != nil {
			fmt.Printf("Erro: %v\n", err)
		} else {
			fmt.Printf("Deletado (status %d)\n", status)
		}

	case "3SUM":
		result, err := client.ThreeSum()
		if err != nil {
			fmt.Printf("Erro: %v\n", err)
		} else {
			fmt.Println("Triplas que somam zero:")
			fmt.Println(result)
		}

	default:
		fmt.Printf("Comando desconhecido: %s\n", *cmd)
		os.Exit(1)
	}
}