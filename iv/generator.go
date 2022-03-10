package pkcs12

import (
	"hash"
)

const (
	DEFAULT_ITERATIONS = 1000
	DEFAULT_KEY_SIZE   = 256
	DEFAULT_IV_SIZE    = 128
)

const (
	KEY_MATERIAL = 1
	IV_MATERIAL  = 2
)

type PKCS12ParameterGenerator interface {
	GenerateDerivedParameters(password string, salt []byte, iterations int) ([]byte, []byte)
}

type pkcs12ParameterGenerator struct {
	keySizeBits, ivSizeBits int
	digestFactory           hash.Hash
}

func pkcs12PasswordToBytes(password string) []byte {
	pkcs12Password := make([]byte, (len(password)+1)*2)
	for i := range pkcs12Password {
		pkcs12Password[i] = 0
	}
	for i, v := range password {
		pkcs12Password[i+i] = byte(v >> 8)
		pkcs12Password[i+i+1] = byte(v & 0xff)
	}
	return pkcs12Password
}

func adjust(a []byte, aOff int, b []byte) {
	x := rune(b[len(b)-1]) + rune(a[aOff+len(b)-1]) + 1
	a[aOff+len(b)-1] = byte(x & 0xff)
	x = x >> 8
	for i := len(b) - 2; i >= 0; i-- {
		x += rune(b[i]) + rune(a[aOff+i])
		a[aOff+i] = byte(x & 0xff)
		x = x >> 8
	}
}

func NewPKCS12ParameterGenerator(digestFactory hash.Hash, keySizeBits, ivSizeBits int) PKCS12ParameterGenerator {
	return &pkcs12ParameterGenerator{
		keySizeBits:   keySizeBits,
		ivSizeBits:    ivSizeBits,
		digestFactory: digestFactory,
	}
}

func (g *pkcs12ParameterGenerator) GenerateDerivedParameters(password string, salt []byte, iterations int) ([]byte, []byte) {
	keySize := g.keySizeBits / 8
	ivSize := g.ivSizeBits / 8

	// pkcs12 padded password (unicode byte array with 2 trailing 0x0 bytes)
	passwordBytes := pkcs12PasswordToBytes(password)

	dKey := g.GenerateDerivedKey(passwordBytes, salt, iterations, KEY_MATERIAL, keySize)
	var dIv []byte
	if ivSize > 0 {
		dIv = g.GenerateDerivedKey(passwordBytes, salt, iterations, IV_MATERIAL, ivSize)
	}
	return dKey, dIv
}

func (g *pkcs12ParameterGenerator) GenerateDerivedKey(password, salt []byte, iterations, idByte, keySize int) []byte {
	u := g.digestFactory.Size()
	v := g.digestFactory.BlockSize()

	dKey := make([]byte, keySize)

	// Step 1
	D := make([]byte, v)
	for i := range D {
		D[i] = byte(idByte)
	}

	// Step 2
	var S []byte
	if len(salt) > 0 {
		saltSize := len(salt)
		sSize := v * ((saltSize + v - 1) / v)
		S = make([]byte, sSize)
		for i := range S {
			S[i] = salt[i%saltSize]
		}
	}

	// Step 3
	var P []byte
	if len(password) > 0 {
		passwordSize := len(password)
		pSize := v * ((passwordSize + v - 1) / v)
		P = make([]byte, pSize)
		for i := range P {
			P[i] = password[i%passwordSize]
		}
	}

	// Step 4
	I := append(S, P...)
	B := make([]byte, v)

	// Step 5
	c := (keySize + u - 1) / u

	// Step 6
	for i := 0; i < c; i++ {
		// Step 6a
		g.digestFactory.Reset()
		g.digestFactory.Write(D)
		g.digestFactory.Write(I)

		A := g.digestFactory.Sum(nil)

		for j := 1; j < iterations; j++ {
			g.digestFactory.Reset()
			g.digestFactory.Write(A)
			A = g.digestFactory.Sum(nil)
		}

		// Step 6b
		for k := 0; k < v; k++ {
			B[k] = A[k%u]
		}

		// step 6c
		for j := 0; j < len(I)/v; j++ {
			adjust(I, j*v, B)
		}

		if i+1 == c {
			for j := 0; j < keySize-i*u; j++ {
				dKey[i*u+j] = A[j]
			}
		} else {
			for j := 0; j < u; j++ {
				dKey[i*u+j] = A[j]
			}
		}

	}

	return dKey
}
