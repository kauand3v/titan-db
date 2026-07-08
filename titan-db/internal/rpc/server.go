package rpc

import (
	"fmt"
	"io"
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/kauandias747474-hue/titandb/internal/kv"
	"github.com/kauandias747474-hue/titandb/internal/raft"
	"github.com/kauandias747474-hue/titandb/pkg/logger"
)

type Server struct {
	address  string
	store    *kv.Store
	raftNode *raft.Node
	listener net.Listener
}

func NewServer(address string, store *kv.Store, raftNode *raft.Node) *Server {
	return &Server{
		address:  address,
		store:    store,
		raftNode: raftNode,
	}
}

func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("falha ao iniciar servidor RPC: %w", err)
	}
	logger.Info("Servidor TitanRPC ouvindo em %s", s.address)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			logger.Info("Erro ao aceitar conexão: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		payloadSize, msgID, cmd, err := DecodeHeader(conn)
		if err != nil {
			if err != io.EOF {
				logger.Info("Erro ao ler cabeçalho: %v", err)
			}
			return
		}

		payload := make([]byte, payloadSize)
		if payloadSize > 0 {
			_, err = io.ReadFull(conn, payload)
			if err != nil {
				logger.Info("Erro ao ler payload: %v", err)
				return
			}
		}

		var respData []byte
		var status uint8 = StatusOK

		switch cmd {
		case CMD_PUT:
			key, value, err := DecodePutRequest(payload)
			if err != nil {
				status = StatusError
				respData = []byte(err.Error())
			} else {
				s.store.Put(key, value)
				logger.Info("PUT key=%s value=%s", key, value)
			}

		case CMD_GET:
			key, err := DecodeGetRequest(payload)
			if err != nil {
				status = StatusError
				respData = []byte(err.Error())
			} else {
				val, ok := s.store.Get(key)
				if ok {
					respData = []byte(val)
				} else {
					status = StatusNotFound
				}
			}

		case CMD_DELETE:
			key, err := DecodeDeleteRequest(payload)
			if err != nil {
				status = StatusError
				respData = []byte(err.Error())
			} else {
				s.store.Delete(key)
				logger.Info("DELETE key=%s", key)
			}

		case CMD_3SUM:
			triplas := s.handleThreeSum()
			respData = []byte(formatTriplets(triplas))

		default:
			status = StatusError
			respData = []byte("comando desconhecido")
		}

		resp := EncodeResponse(status, respData)
		_, err = conn.Write(resp)
		if err != nil {
			logger.Info("Erro ao enviar resposta: %v", err)
			return
		}
		_ = msgID // evita warning
	}
}

func (s *Server) handleThreeSum() [][]int {
	allVals := s.store.GetAllValues()
	nums := make([]int, 0, len(allVals))
	for _, v := range allVals {
		if n, err := strconv.Atoi(v); err == nil {
			nums = append(nums, n)
		}
	}
	return threeSum(nums)
}

func threeSum(nums []int) [][]int {
	sort.Ints(nums)
	n := len(nums)
	res := make([][]int, 0)

	for i := 0; i < n-2; i++ {
		if nums[i] > 0 {
			break
		}
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		left, right := i+1, n-1
		for left < right {
			sum := nums[i] + nums[left] + nums[right]
			if sum == 0 {
				res = append(res, []int{nums[i], nums[left], nums[right]})
				for left < right && nums[left] == nums[left+1] {
					left++
				}
				for left < right && nums[right] == nums[right-1] {
					right--
				}
				left++
				right--
			} else if sum < 0 {
				left++
			} else {
				right--
			}
		}
	}
	return res
}

func formatTriplets(triplets [][]int) string {
	if len(triplets) == 0 {
		return "nenhuma tripla encontrada"
	}
	strs := make([]string, len(triplets))
	for i, t := range triplets {
		strs[i] = fmt.Sprintf("%d,%d,%d", t[0], t[1], t[2])
	}
	return strings.Join(strs, "; ")
}

func (s *Server) Stop() {
	if s.listener != nil {
		s.listener.Close()
	}
}