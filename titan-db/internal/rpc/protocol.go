package rpc

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Comandos
const (
	CMD_PUT    = 1
	CMD_GET    = 2
	CMD_DELETE = 3
	CMD_3SUM   = 4
)

// Status
const (
	StatusOK       = 0
	StatusNotFound = 1
	StatusError    = 2
)

// ----- Codificação das requisições -----
func EncodePutRequest(key, value string) []byte {
	buf := make([]byte, 2+len(key)+2+len(value))
	binary.BigEndian.PutUint16(buf[0:2], uint16(len(key)))
	copy(buf[2:], key)
	binary.BigEndian.PutUint16(buf[2+len(key):], uint16(len(value)))
	copy(buf[2+len(key)+2:], value)
	return buf
}

func EncodeGetRequest(key string) []byte {
	buf := make([]byte, 2+len(key))
	binary.BigEndian.PutUint16(buf[0:2], uint16(len(key)))
	copy(buf[2:], key)
	return buf
}

func EncodeDeleteRequest(key string) []byte {
	return EncodeGetRequest(key) // mesmo formato
}

// ----- Decodificação do cabeçalho (para o servidor) -----
func DecodeHeader(reader io.Reader) (payloadSize uint32, msgID uint32, cmd uint8, err error) {
	header := make([]byte, 9)
	_, err = io.ReadFull(reader, header)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("erro ao ler cabeçalho: %w", err)
	}
	payloadSize = binary.BigEndian.Uint32(header[0:4])
	msgID = binary.BigEndian.Uint32(header[4:8])
	cmd = header[8]
	return payloadSize, msgID, cmd, nil
}

// ----- Decodificação dos payloads (para o servidor) -----
func DecodePutRequest(payload []byte) (key, value string, err error) {
	if len(payload) < 2 {
		return "", "", fmt.Errorf("payload muito curto para PUT")
	}
	keyLen := binary.BigEndian.Uint16(payload[0:2])
	if int(2+keyLen) > len(payload) {
		return "", "", fmt.Errorf("payload PUT: key_len excede tamanho")
	}
	key = string(payload[2 : 2+keyLen])
	rest := payload[2+keyLen:]
	if len(rest) < 2 {
		return "", "", fmt.Errorf("payload PUT: falta value_len")
	}
	valLen := binary.BigEndian.Uint16(rest[0:2])
	if int(2+valLen) > len(rest) {
		return "", "", fmt.Errorf("payload PUT: value_len excede tamanho")
	}
	value = string(rest[2 : 2+valLen])
	return key, value, nil
}

func DecodeGetRequest(payload []byte) (key string, err error) {
	if len(payload) < 2 {
		return "", fmt.Errorf("payload muito curto para GET/DELETE")
	}
	keyLen := binary.BigEndian.Uint16(payload[0:2])
	if int(2+keyLen) != len(payload) {
		return "", fmt.Errorf("tamanho do payload não confere com key_len")
	}
	key = string(payload[2 : 2+keyLen])
	return key, nil
}

func DecodeDeleteRequest(payload []byte) (key string, err error) {
	return DecodeGetRequest(payload)
}

// ----- Codificação da resposta (usada pelo servidor) -----
func EncodeResponse(status uint8, data []byte) []byte {
	resp := make([]byte, 9+len(data))
	binary.BigEndian.PutUint32(resp[0:4], uint32(len(data)))
	binary.BigEndian.PutUint32(resp[4:8], 0) // msgID (não usado)
	resp[8] = status
	copy(resp[9:], data)
	return resp
}