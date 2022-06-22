package grpc

import (
	"context"

	"github.com/hkhait/crypto/encryption/pbe"
	"github.com/hkhait/crypto/grpc/cryptoservice"
)

type CryptoServer struct {
	cryptoservice.UnimplementedCryptoServiceServer
	encryptor pbe.PBEStringPasswordEncryptor
}

func NewCryptoServer(encyptor pbe.PBEStringPasswordEncryptor) *CryptoServer {
	return &CryptoServer{encryptor: encyptor}
}

func (s *CryptoServer) Decrypt(ctx context.Context, in *cryptoservice.DecryptRequest) (*cryptoservice.DecryptResponse, error) {
	plaintext, err := s.encryptor.Decrypt(in.GetCiphertext())
	return &cryptoservice.DecryptResponse{Plaintext: plaintext}, err
}
