package butin

import (
	"fmt"
	"testing"
	"time"
)

func TestGenSignature(t *testing.T) {
	timestamp := time.Now().Unix()
	secretKey := "jerryzhuo@abcd"
	signature, err := GenSignature(timestamp, secretKey)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(signature))

	if err := CheckSignature(signature, secretKey); err != nil {
		panic(err)
	}
}
