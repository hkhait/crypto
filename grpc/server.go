package grpc

import (
	"context"

	"github.com/hkhait/crypto/encryption/pbe"
	"github.com/hkhait/crypto/grpc/cryptoservice"
)

type CryptoServer struct {
	cryptoservice.UnimplementedCryptoServiceServer
	password string
}

func NewCryptoServer(password string) *CryptoServer {
	return &CryptoServer{password: password}
}

func (s *CryptoServer) Decrypt(ctx context.Context, in *cryptoservice.DecryptRequest) (*cryptoservice.DecryptResponse, error) {
	encryptor := pbe.NewStandardPBEStringEncryptor()
	encryptor.SetPassword(s.password)
	plaintext, err := encryptor.Decrypt(in.GetCiphertext())
	return &cryptoservice.DecryptResponse{Plaintext: plaintext}, err
}
