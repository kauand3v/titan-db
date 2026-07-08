package raft

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/kauandias747474-hue/titandb/internal/kv"
	"github.com/kauandias747474-hue/titandb/pkg/logger"
)

const (
	Follower  = 0
	Candidate = 1
	Leader    = 2
)

const (
	msgRequestVote    = 10
	msgVoteResponse   = 11
	msgAppendEntries  = 12
	msgAppendResponse = 13
)

type Node struct {
	id          int
	store       *kv.Store
	peers       map[int]string
	mu          sync.Mutex

	currentTerm int
	votedFor    int
	role        int
	leaderID    int

	electionTimer  *time.Timer
	heartbeatTimer *time.Ticker

	stopCh chan struct{}
}

func NewNode(id int, store *kv.Store, peers map[int]string) *Node {
	return &Node{
		id:    id,
		store: store,
		peers: peers,
		role:  Follower,
		stopCh: make(chan struct{}),
	}
}

func (n *Node) Start() {
	logger.Info("Nó Raft %d iniciado como Follower", n.id)
	n.resetElectionTimer()
	go n.run()
}

func (n *Node) run() {
	for {
		select {
		case <-n.electionTimer.C:
			n.startElection()
		case <-n.stopCh:
			return
		}
	}
}

func (n *Node) resetElectionTimer() {
	if n.electionTimer != nil {
		n.electionTimer.Stop()
	}
	timeout := time.Duration(150+rand.Intn(150)) * time.Millisecond
	n.electionTimer = time.NewTimer(timeout)
}

func (n *Node) sendRaftMessage(conn net.Conn, msgType uint8, payload []byte) {
	header := make([]byte, 5)
	binary.BigEndian.PutUint32(header[0:4], uint32(len(payload)))
	header[4] = msgType
	conn.Write(header)
	if len(payload) > 0 {
		conn.Write(payload)
	}
}

func (n *Node) readRaftMessage(conn net.Conn) (msgType uint8, payload []byte, err error) {
	header := make([]byte, 5)
	_, err = io.ReadFull(conn, header)
	if err != nil {
		return 0, nil, err
	}
	payloadSize := binary.BigEndian.Uint32(header[0:4])
	msgType = header[4]
	payload = make([]byte, payloadSize)
	if payloadSize > 0 {
		_, err = io.ReadFull(conn, payload)
	}
	return
}

func (n *Node) HandleRaftConnection(conn net.Conn) {
	defer conn.Close()
	msgType, payload, err := n.readRaftMessage(conn)
	if err != nil {
		return
	}
	switch msgType {
	case msgRequestVote:
		n.handleVoteRequest(conn, payload)
	case msgAppendEntries:
		n.handleAppendEntries(conn, payload)
	}
}

func (n *Node) handleVoteRequest(conn net.Conn, payload []byte) {
	if len(payload) < 16 {
		return
	}
	term := binary.BigEndian.Uint64(payload[0:8])
	candidateID := int(binary.BigEndian.Uint64(payload[8:16]))

	n.mu.Lock()
	grant := false
	if term > uint64(n.currentTerm) {
		n.currentTerm = int(term)
		n.votedFor = -1
		n.role = Follower
	}
	if term == uint64(n.currentTerm) && (n.votedFor == -1 || n.votedFor == candidateID) {
		grant = true
		n.votedFor = candidateID
		n.resetElectionTimer()
	}
	n.mu.Unlock()

	resp := []byte{0}
	if grant {
		resp[0] = 1
	}
	n.sendRaftMessage(conn, msgVoteResponse, resp)
}

func (n *Node) startHeartbeats() {
	n.heartbeatTimer = time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			select {
			case <-n.heartbeatTimer.C:
				n.sendHeartbeats()
			case <-n.stopCh:
				n.heartbeatTimer.Stop()
				return
			}
		}
	}()
}

func (n *Node) sendHeartbeats() {
	n.mu.Lock()
	if n.role != Leader {
		n.mu.Unlock()
		return
	}
	term := n.currentTerm
	n.mu.Unlock()

	for peerID, addr := range n.peers {
		if peerID == n.id {
			continue
		}
		go func(address string) {
			conn, err := net.DialTimeout("tcp", address, 200*time.Millisecond)
			if err != nil {
				return
			}
			defer conn.Close()
			payload := make([]byte, 36)
			binary.BigEndian.PutUint64(payload[0:8], uint64(term))
			binary.BigEndian.PutUint64(payload[8:16], uint64(n.id))
			n.sendRaftMessage(conn, msgAppendEntries, payload)
		}(addr)
	}
}

func (n *Node) handleAppendEntries(conn net.Conn, _ []byte) {
	n.mu.Lock()
	n.resetElectionTimer()
	n.role = Follower
	n.mu.Unlock()
	n.sendRaftMessage(conn, msgAppendResponse, []byte{1})
}

func (n *Node) SubmitCommand(cmd []byte) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.role != Leader {
		return fmt.Errorf("não é o líder")
	}
	n.applyCommand(cmd)
	return nil
}

func (n *Node) applyCommand(cmd []byte) {
	if len(cmd) < 1 {
		return
	}
	cmdType := cmd[0]
	payload := cmd[1:]
	switch cmdType {
	case 1: // PUT
		if len(payload) < 2 {
			return
		}
		keyLen := binary.BigEndian.Uint16(payload[0:2])
		key := string(payload[2 : 2+keyLen])
		valStart := int(2 + keyLen + 2)
		if len(payload) < valStart {
			return
		}
		valLen := binary.BigEndian.Uint16(payload[2+keyLen : valStart])
		val := string(payload[valStart : int(valStart)+int(valLen)])
		n.store.Put(key, val)
	case 3: // DELETE
		if len(payload) < 2 {
			return
		}
		keyLen := binary.BigEndian.Uint16(payload[0:2])
		key := string(payload[2 : 2+keyLen])
		n.store.Delete(key)
	}
}

func (n *Node) Shutdown() {
	close(n.stopCh)
	logger.Info("Nó Raft %d encerrado", n.id)
}