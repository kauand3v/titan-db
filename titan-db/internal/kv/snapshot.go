package kv

type Snapshot struct {
	Data map[string]string
}

func (s *Store) TakeSnapshot() *Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	dataCopy := make(map[string]string, len(s.data))
	for k, v := range s.data {
		dataCopy[k] = v
	}
	return &Snapshot{Data: dataCopy}
}

func (s *Store) RestoreSnapshot(snap *Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]string)
	for k, v := range snap.Data {
		s.data[k] = v
	}
}