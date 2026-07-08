package raft

import (
	"encoding/binary"
	"net"
	"time"

	"github.com/kauandias747474-hue/titandb/pkg/logger"
)

func (n *Node) startElection() {
	n.mu.Lock()
	n.currentTerm++
	n.votedFor = n.id
	n.role = Candidate
	term := n.currentTerm
	logger.Info("Nó %d tornou-se Candidate para o termo %d", n.id, term)
	n.mu.Unlock()

	votes := make(chan bool, len(n.peers)-1)
	for peerID, addr := range n.peers {
		if peerID == n.id {
			continue
		}
		go func(id int, address string) {
			ok := n.requestVote(id, address, term)
			votes <- ok
		}(peerID, addr)
	}

	granted := 1
	for i := 0; i < len(n.peers)-1; i++ {
		if <-votes {
			granted++
		}
	}

	n.mu.Lock()
	defer n.mu.Unlock()
	if granted > len(n.peers)/2 && n.role == Candidate && n.currentTerm == term {
		n.role = Leader
		n.leaderID = n.id
		logger.Info("Nó %d tornou-se Líder para o termo %d", n.id, term)
		n.startHeartbeats()
	} else {
		n.role = Follower
		n.resetElectionTimer()
	}
}

func (n *Node) requestVote(_ int, addr string, term int) bool {
	conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
	if err != nil {
		return false
	}
	defer conn.Close()

	payload := make([]byte, 16)
	binary.BigEndian.PutUint64(payload[0:8], uint64(term))
	binary.BigEndian.PutUint64(payload[8:16], uint64(n.id))
	n.sendRaftMessage(conn, msgRequestVote, payload)

	respType, respPayload, err := n.readRaftMessage(conn)
	if err != nil || respType != msgVoteResponse {
		return false
	}
	if len(respPayload) < 1 {
		return false
	}
	return respPayload[0] == 1
}