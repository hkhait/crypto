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

	password := flag.String("password", os.Getenv("DECRYPT_PASSWORD"), "the password used to encrypt/decrypt")
	ciphertext := flag.String("ciphertext", "-", "the ciphertext to be decrypted, read from stdin if not specified")
	plaintext := flag.String("plaintext", "-", "the plaintext to be encrypted, read from stdin if not specified")
	encrypt := flag.Bool("encrypt", false, "encrypt mode")
	decrypt := flag.Bool("decrypt", false, "decrypt mode")
	serverMode := flag.Bool("server", false, "run decryptor as GRPC server")
	serverPort := flag.Int("port", 8080, "the server port when running as GRPC server")
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
			s := bufio.NewScanner(r)
			for s.Scan() {
				encryptor := pbe.NewStandardPBEStringEncryptor()
				encryptor.SetPassword(*password)
				ciphertext, err := encryptor.Encrypt(s.Text())
				if err != nil {
					panic(err)
				}
				fmt.Fprintln(w, ciphertext)
				if err != nil {
					panic(err)
				}
			}
		} else {
			ciphertext, err := encryptor.Encrypt(*plaintext)
			if err != nil {
				panic(err)
			}
			fmt.Fprintln(w, ciphertext)
			if err != nil {
				panic(err)
			}
		}
	}

	if *decrypt {
		if *ciphertext == "-" {
			r := bufio.NewReader(os.Stdin)
			s := bufio.NewScanner(r)
			for s.Scan() {
				decryptor := pbe.NewStandardPBEStringEncryptor()
				decryptor.SetPassword(*password)
				plaintext, err := decryptor.Decrypt(s.Text())
				if err != nil {
					panic(err)
				}
				fmt.Fprintln(w, plaintext)
				if err != nil {
					panic(err)
				}
			}
		} else {
			plaintext, err := encryptor.Decrypt(*ciphertext)
			if err != nil {
				panic(err)
			}
			fmt.Fprintln(w, plaintext)
			if err != nil {
				panic(err)
			}
		}
	}

	if *serverMode {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *serverPort))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		cryptoservice.RegisterCryptoServiceServer(s, cryptogrpc.NewCryptoServer(*password))
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}

}
