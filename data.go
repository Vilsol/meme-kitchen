package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"google.golang.org/protobuf/proto"
	"io"
	"memekitchen/data"
)

func EncodeData(payload *data.Payload) (string, error) {
	marshaled, err := proto.Marshal(payload)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	writer := gzip.NewWriter(buf)

	if _, err := writer.Write(marshaled); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	if buf.Len() > len(marshaled) {
		// If compression made it bigger, don't use compression
		return base64.RawURLEncoding.EncodeToString(marshaled), err
	}

	return base64.RawURLEncoding.EncodeToString(buf.Bytes()), err
}

func DecodeData(str string) (*data.Payload, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}

	// Might just not be compressed
	out := &data.Payload{}
	if err := proto.Unmarshal(decoded, out); err == nil {
		return out, nil
	}

	println(len(decoded), string(decoded))

	reader, err := gzip.NewReader(bytes.NewReader(decoded))
	if err != nil {
		if errors.Is(err, gzip.ErrHeader) {
			println("B")
			return out, nil
		}

		println("C")
		return nil, err
	}
	defer reader.Close()

	marshaled, err := io.ReadAll(reader)
	if err != nil {
		println("D")
		return nil, err
	}

	if err := proto.Unmarshal(marshaled, out); err != nil {
		return nil, err
	}

	return out, nil
}
