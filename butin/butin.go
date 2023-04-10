package butin

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"net"
)

const (
	ProtocolVersion = 0

	ProtocolVersionLen = 1
	ProtocolCmdLen     = 1
	ProtocolDataLen    = 4
	ProtocolMaxLen     = 256 * 256
)

const (
	CmdSignCheck Cmd = 0 + iota
	CmdSignSuccess
	CmdSignError

	CmdPing
	CmdPong

	CmdConnectionNew
	CmdConnectionData
	CmdConnectionClientCloseReq
	CmdConnectionClientCloseRsp
	CmdConnectionServerCloseReq
	CmdConnectionServerCloseRsp
)

type Cmd int

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

func GenBuf() []byte {
	return make([]byte, ProtocolMaxLen)
}

func GenDataBuf() []byte {
	return make([]byte, ProtocolMaxLen-ProtocolVersionLen-ProtocolCmdLen-ProtocolDataLen)
}

func ReadCmd(buf []byte) Cmd {
	return Cmd(buf[ProtocolCmdLen])
}

func ReadData(buf []byte) []byte {
	offset := ProtocolVersionLen + ProtocolCmdLen + ProtocolDataLen

	var dataLen int32
	bytesBuffer := bytes.NewBuffer(buf[offset-ProtocolDataLen : offset])
	binary.Read(bytesBuffer, binary.BigEndian, &dataLen)

	return buf[offset : offset+int(dataLen)]
}

func WriteCmdAndData(cnn net.Conn, cmd Cmd, data []byte) {
	if len(data) == 0 {
		cnn.Write([]byte{ProtocolVersion, byte(cmd)})
		cnn.Write(make([]byte, ProtocolDataLen))
		return
	}

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, int32(len(data)))

	cnn.Write([]byte{ProtocolVersion, byte(cmd)})
	cnn.Write(bytesBuffer.Bytes())
	cnn.Write(data)
}
