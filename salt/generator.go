package salt

import (
	"crypto/rand"
)

const (
	DEFAULT_BLOCK_SIZE = 16
)

type SaltGenerator interface {
	GenerateSalt() ([]byte, error)
}

type randomSaltGenerator struct {
	blockSize int
}

func NewDefaultRandomSaltGenerator() SaltGenerator {
	return &randomSaltGenerator{blockSize: DEFAULT_BLOCK_SIZE}
}

func NewRandomSaltGenerator(blockSize int) SaltGenerator {
	return &randomSaltGenerator{blockSize: blockSize}
}

func (g *randomSaltGenerator) GenerateSalt() ([]byte, error) {
	b := make([]byte, g.blockSize)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
