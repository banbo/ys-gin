package main

import (
	"flag"
	"fmt"
	"github.com/banbo/ys-gin/utils/aes"
)

var str = flag.String("i", "", "stream to encrypt or decrypt")
var kind = flag.String("type", "encode", "type to encrypt or decrypt;eg:encode is encrypt,decode is decrypt")

func main() {
	flag.Parse()

	switch *kind {
	case "encode":
		xpass, err := aes.Encrypt([]byte(*str))
		if err != nil {
			fmt.Println(err)
			fmt.Println("-i stream for encrypt or decrypt\n-type encode is encrypt,decode is decrypt")
			return
		}
		fmt.Printf("加密后: %s\n", string(xpass))
	case "decode":
		tpass, err := aes.Decrypt([]byte(*str))
		if err != nil {
			fmt.Println(err)
			fmt.Println("-i stream for encrypt or decrypt\n-type encode is encrypt,decode is decrypt")
			return
		}
		fmt.Printf("解密后: %s\n", tpass)
	default:
		fmt.Println("-i stream for encrypt or decrypt\n-type encode is encrypt,decode is decrypt")
	}
	return
}
