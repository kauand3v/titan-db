package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/kauandias747474-hue/titandb/internal/kv"
	"github.com/kauandias747474-hue/titandb/internal/raft"
	"github.com/kauandias747474-hue/titandb/internal/rpc"
	"github.com/kauandias747474-hue/titandb/pkg/logger"
)

func main() {
	idFlag := flag.Int("id", 0, "ID do nó (obrigatório)")
	addressFlag := flag.String("address", "", "Endereço de escuta (ex: :8001)")
	clusterFlag := flag.String("cluster", "", "Lista de peers: 1@:8001,2@:8002,3@:8003")
	flag.Parse()

	if *idFlag == 0 || *addressFlag == "" || *clusterFlag == "" {
		fmt.Println("Erro: --id, --address e --cluster são obrigatórios")
		os.Exit(1)
	}

	logger.Init(*idFlag)
	logger.Info("Iniciando nó %d em %s", *idFlag, *addressFlag)

	// Parse da lista de peers
	peers := make(map[int]string)
	parts := strings.Split(*clusterFlag, ",")
	foundSelf := false
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		idAddr := strings.SplitN(p, "@", 2)
		if len(idAddr) != 2 {
			fmt.Printf("Formato inválido para peer: %s\n", p)
			os.Exit(1)
		}
		id, err := strconv.Atoi(strings.TrimSpace(idAddr[0]))
		if err != nil {
			fmt.Printf("ID inválido em peer: %s\n", p)
			os.Exit(1)
		}
		addr := strings.TrimSpace(idAddr[1])
		peers[id] = addr
		if id == *idFlag {
			foundSelf = true
		}
	}
	if !foundSelf {
		fmt.Printf("ID %d não encontrado na lista de peers\n", *idFlag)
		os.Exit(1)
	}

	store := kv.NewStore()
	raftNode := raft.NewNode(*idFlag, store, peers)
	// rpc.NewServer agora aceita 3 argumentos: address, store, raftNode
	rpcServer := rpc.NewServer(*addressFlag, store, raftNode)

	raftNode.Start()

	go func() {
		if err := rpcServer.Start(); err != nil {
			logger.Info("Erro ao iniciar servidor RPC: %v", err)
			os.Exit(1)
		}
	}()

	logger.Info("Nó %d pronto. Aguardando conexões...", *idFlag)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("Encerrando nó %d...", *idFlag)
	rpcServer.Stop()
	raftNode.Shutdown()
	store.Close()
}