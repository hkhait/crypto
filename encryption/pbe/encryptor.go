package pbe

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"

	"github.com/hkhait/crypto/encryption"
	pkcs12 "github.com/hkhait/crypto/iv"
	"github.com/hkhait/crypto/salt"
)

type PBEStringPasswordEncryptor interface {
	encryption.StringEncryptor
	SetPassword(password string)
}

type standardPBEStringEncryptor struct {
	SaltGenerator salt.SaltGenerator
	KeyGenerator  pkcs12.PKCS12ParameterGenerator
	Password      string
}

func NewStandardPBEStringEncryptor() PBEStringPasswordEncryptor {
	saltGenerator := salt.NewDefaultRandomSaltGenerator()
	keyGenerator := pkcs12.NewPKCS12ParameterGenerator(sha256.New(), pkcs12.DEFAULT_KEY_SIZE, pkcs12.DEFAULT_IV_SIZE)

	return &standardPBEStringEncryptor{
		SaltGenerator: saltGenerator,
		KeyGenerator:  keyGenerator,
	}
}

func (encryptor *standardPBEStringEncryptor) Encrypt(text string) (string, error) {
	// generate a 16 byte salt which is used to generate key material and iv
	salt, err := encryptor.SaltGenerator.GenerateSalt()
	if err != nil {
		return "", err
	}

	// generate key material
	key, iv := encryptor.KeyGenerator.GenerateDerivedParameters(encryptor.Password, salt, pkcs12.DEFAULT_ITERATIONS)

	// setup AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCEncrypter(block, iv)

	// pad the plain text secret to AES block size
	plaintext := pad(block.BlockSize(), []byte(text))

	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	// concatenate salt + encrypted message
	res := append(salt, ciphertext...)

	return base64.StdEncoding.EncodeToString(res), nil
}

func (encryptor *standardPBEStringEncryptor) Decrypt(ciphertext string) (string, error) {
	nCipherBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	saltBytes := nCipherBytes[:salt.DEFAULT_BLOCK_SIZE]

	// create reverse key material
	key, iv := encryptor.KeyGenerator.GenerateDerivedParameters(encryptor.Password, saltBytes, pkcs12.DEFAULT_ITERATIONS)

	// setup AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	nCipherMessage := nCipherBytes[salt.DEFAULT_BLOCK_SIZE:]

	plaintext := make([]byte, len(nCipherMessage))
	mode.CryptBlocks(plaintext, nCipherMessage)

	return string(unpad(plaintext)), nil
}

func (encryptor *standardPBEStringEncryptor) SetPassword(password string) {
	encryptor.Password = password
}

func pad(blockSize int, s []byte) []byte {
	padSize := blockSize - (len(s) % blockSize)
	pad := bytes.Repeat([]byte{byte(padSize)}, padSize)
	return append(s, pad...)
}

func unpad(s []byte) []byte {
	length := len(s)
	unpadSize := int(s[length-1])
	return s[:(length - unpadSize)]
}
