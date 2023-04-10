package butin

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

func GenSignature(timestamp int64, secretKey string) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, timestamp); err != nil {
		return nil, err
	}

	key := fmt.Sprintf("%d%s", timestamp, secretKey)
	keyByte := md5.Sum(hmac.New(sha256.New, nil).Sum([]byte(key)))

	if err := binary.Write(buf, binary.BigEndian, keyByte); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func CheckSignature(signature []byte, secretKey string) error {
	if len(signature) != (64+128)/8 {
		return fmt.Errorf("signature length error:%d", len(signature))
	}

	buf := bytes.NewBuffer([]byte(signature))

	var timestamp int64
	if err := binary.Read(buf, binary.BigEndian, &timestamp); err != nil {
		return err
	}

	// if time.Now().Unix()-timestamp > 60 {
	// 	return fmt.Errorf("signature timeout")
	// }

	keyByte := make([]byte, 128/8)
	if err := binary.Read(buf, binary.BigEndian, &keyByte); err != nil {
		return err
	}

	key := fmt.Sprintf("%d%s", timestamp, secretKey)
	keyByteC := md5.Sum(hmac.New(sha256.New, nil).Sum([]byte(key)))

	if string(keyByte) != string(keyByteC[:]) {
		return fmt.Errorf("signature error")
	}

	return nil
}
