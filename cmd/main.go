package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/hkhait/crypto/encryption/pbe"
)

func main() {

	password := flag.String("password", "", "password")
	ciphertext := flag.String("ciphertext", "-", "ciphertext")
	plaintext := flag.String("plaintext", "-", "plaintext")
	encrypt := flag.Bool("encrypt", false, "encrypt")
	decrypt := flag.Bool("decrypt", false, "decrypt")
	flag.Parse()

	if *encrypt && *decrypt {
		panic(fmt.Errorf("encrypt and decrypt are mutually exclusive"))
	}

	encryptor := pbe.NewStandardPBEStringEncryptor()
	encryptor.SetPassword(*password)

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	if *encrypt {
		if *plaintext == "-" {
			r := bufio.NewReader(os.Stdin)
			*plaintext, _ = r.ReadString(0)
		}
		ciphertext, err := encryptor.Encrypt(*plaintext)
		if err != nil {
			panic(err)
		}
		_, err = w.WriteString(ciphertext)
		if err != nil {
			panic(err)
		}
	}

	if *decrypt {
		if *ciphertext == "-" {
			r := bufio.NewReader(os.Stdin)
			*ciphertext, _ = r.ReadString(0)
		}
		plaintext, err := encryptor.Decrypt(*ciphertext)
		if err != nil {
			panic(err)
		}
		_, err = w.WriteString(plaintext)
		if err != nil {
			panic(err)
		}
	}

}
