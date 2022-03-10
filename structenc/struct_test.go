package structenc

import (
	"testing"

	"github.com/hkhait/crypto/encryption/pbe"
)

const password = "password"

type TestStruct struct {
	A string  `ha:"encrypt"`
	B *string `ha:"encrypt"`
	C string
	D *string
}

func TestEncryptStruct(t *testing.T) {
	b, d := "b", "d"
	s := TestStruct{
		A: "a",
		B: &b,
		C: "c",
		D: &d,
	}
	e := pbe.NewStandardPBEStringEncryptor()
	e.SetPassword(password)
	err := EncryptStruct(e, &s)
	if err != nil {
		t.Errorf("EncryptStruct error: %v", err)
	}
	deca, err := e.Decrypt(s.A)
	if err != nil {
		t.Errorf("Fail to decrypt encrypted string: %v", err)
	}
	if deca != "a" {
		t.Errorf("Fail to decrypt encrypted string: expected %s, found %s", "a", deca)
	}
	decb, err := e.Decrypt(*s.B)
	if err != nil {
		t.Errorf("Fail to decrypt encrypted string pointer: %v", err)
	}
	if decb != "b" {
		t.Errorf("Fail to decrypt encrypted string pointer: expected %s, found %s", "b", decb)
	}
	if s.C != "c" {
		t.Errorf("String field should not be encrypted: expected %s, found %s", "c", s.C)
	}
	if *s.D != "d" {
		t.Errorf("String pointer field should not be encrypted: expected %s, found %s", "d", *s.D)
	}
}

func TestDecryptStruct(t *testing.T) {
	b, d := "xCd93tBGCUzAnzEw9kCKh+iE6p7BICDoMas3sbGz3iY=", "l5GzbIWadaqx0XAAFRa4MFl2LPy6mUFOvesW/nh7mZc="
	s := TestStruct{
		A: "5HRJrJDadhU118emSH3T2oRKHx5wb1nk/RpZNlDAtHc=",
		B: &b,
		C: "M5VJJ+/8EvO7dePyTsYnItlP74supi/mhkgl9W4p5/o=",
		D: &d,
	}
	e := pbe.NewStandardPBEStringEncryptor()
	e.SetPassword(password)
	err := DecryptStruct(e, &s)
	if err != nil {
		t.Errorf("DecryptStruct error: %v", err)
	}
	if s.A != "a" {
		t.Errorf("Fail to decrypt encrypted string: expected %s, found %s", "a", s.A)
	}
	if *s.B != "b" {
		t.Errorf("Fail to decrypt encrypted string pointer: expected %s, found %s", "b", *s.B)
	}
	if s.C != "M5VJJ+/8EvO7dePyTsYnItlP74supi/mhkgl9W4p5/o=" {
		t.Errorf("String field should not be encrypted: expected %s, found %s", "M5VJJ+/8EvO7dePyTsYnItlP74supi/mhkgl9W4p5/o=", s.C)
	}
	if *s.D != "l5GzbIWadaqx0XAAFRa4MFl2LPy6mUFOvesW/nh7mZc=" {
		t.Errorf("String pointer field should not be encrypted: expected %s, found %s", "l5GzbIWadaqx0XAAFRa4MFl2LPy6mUFOvesW/nh7mZc=", *s.D)
	}
}
