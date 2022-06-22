package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hkhait/crypto/encryption/pbe"
	cryptogrpc "github.com/hkhait/crypto/grpc"
	"github.com/hkhait/crypto/grpc/cryptoservice"
	"google.golang.org/grpc"
)

func main() {

	password := flag.String("password", os.Getenv("DECRYPT_PASSWORD"), "password")
	ciphertext := flag.String("ciphertext", "-", "ciphertext")
	plaintext := flag.String("plaintext", "-", "plaintext")
	encrypt := flag.Bool("encrypt", false, "encrypt")
	decrypt := flag.Bool("decrypt", false, "decrypt")
	serverMode := flag.Bool("server", false, "activate server mode")
	serverPort := flag.Int("port", 8080, "the server port")
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

	if *serverMode {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *serverPort))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		cryptoservice.RegisterCryptoServiceServer(s, cryptogrpc.NewCryptoServer(encryptor))
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}

}
