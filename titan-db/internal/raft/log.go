package raft

// LogEntry é a única declaração (não duplicada).
type LogEntry struct {
	Term    int
	Command []byte
}