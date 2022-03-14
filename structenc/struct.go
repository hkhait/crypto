package structenc

import (
	"fmt"
	"reflect"

	"github.com/hkhait/crypto/encryption"
)

func EncryptStruct(e encryption.StringEncryptor, v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("type error: v must be a pointer to a struct")
	}
	rv := reflect.ValueOf(v).Elem()
	rvt := rv.Type()
	if rvt.Kind() != reflect.Struct {
		return fmt.Errorf("type error: v must be a pointer to a struct")
	}
	for i := 0; i < rv.NumField(); i++ {
		if rvt.Field(i).Tag.Get("ha") != "encrypted" {
			continue
		}
		switch rv.Field(i).Kind() {
		case reflect.String:
			val := rv.Field(i)
			ciphertext, err := e.Encrypt(val.String())
			if err != nil {
				return fmt.Errorf("encrypt error: %v", err)
			}
			val.SetString(ciphertext)
		case reflect.Ptr:
			val := rv.Field(i).Elem()
			if val.Kind() == reflect.String {
				ciphertext, err := e.Encrypt(val.String())
				if err != nil {
					return fmt.Errorf("encrypt error: %v", err)
				}
				val.SetString(ciphertext)
			}
		}
	}

	return nil
}

func DecryptStruct(e encryption.StringEncryptor, v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("type error: v must be a pointer to a struct")
	}
	rv := reflect.ValueOf(v).Elem()
	rvt := rv.Type()
	if rvt.Kind() != reflect.Struct {
		return fmt.Errorf("type error: v must be a pointer to a struct")
	}
	for i := 0; i < rv.NumField(); i++ {
		if rvt.Field(i).Tag.Get("ha") != "encrypted" {
			continue
		}
		switch rv.Field(i).Kind() {
		case reflect.String:
			val := rv.Field(i)
			decrypted, err := e.Decrypt(val.String())
			if err != nil {
				return fmt.Errorf("decrypt error: %v", err)
			}
			val.SetString(decrypted)
		case reflect.Ptr:
			val := rv.Field(i).Elem()
			if val.Kind() == reflect.String {
				decrypted, err := e.Decrypt(val.String())
				if err != nil {
					return fmt.Errorf("decrypt error: %v", err)
				}
				val.SetString(decrypted)
			}
		}
	}

	return nil
}
