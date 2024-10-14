package uti

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type KeyEncoding int

const (
	KeyB64Der KeyEncoding = 1 << iota
	KeyPEM
	KeyCert
	KeyAny = KeyB64Der | KeyPEM | KeyCert
)

const (
	ErrInvalidKeyType ErrorCode = 1001 + iota
	ErrInvalidKey
	ErrInvalidECDSA
	ErrFailedECDSA
)

func CryptoDecodeKey(key string, enc KeyEncoding) ([]byte, error) {
	// Base64 encoded DER?
	if enc&KeyB64Der > 0 {
		if derBytes, err := base64.StdEncoding.DecodeString(key); err == nil {
			pubKey, err := x509.ParsePKIXPublicKey(derBytes)
			if err != nil && enc&KeyCert > 0 {
				if cert, err := x509.ParseCertificate(derBytes); err == nil {
					pubKey = cert.PublicKey
					derBytes, _ = x509.MarshalPKIXPublicKey(pubKey)
				} else {
					return nil, Errorf(ErrInvalidKey, "Invalid key")
				}
			}

			switch pubKey.(type) {
			// reject dsa
			case *rsa.PublicKey, *ecdsa.PublicKey, ed25519.PublicKey:
				return derBytes, nil
			default:
				return nil, Errorf(ErrInvalidKeyType, `Unsupported public key type '%v'`, reflect.TypeOf(pubKey))
			}
		}
	}

	// PEM?
	if enc&KeyPEM > 0 {
		block, _ := pem.Decode([]byte(key))
		if block == nil || block.Type != "PUBLIC KEY" {
			return nil, Errorf(ErrInvalidKey, "Invalid key")
		} else if pubKey, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
			return nil, err
		} else if derEncodedPublicKey, err := x509.MarshalPKIXPublicKey(pubKey); err != nil {
			return nil, err
		} else {
			return derEncodedPublicKey, nil
		}
	}

	return nil, Errorf(ErrInvalidKey, "Invalid key")
}

func CryptoVerifySignature(msg []byte, ecdsaSign []byte, key string, allowEncoding KeyEncoding) (bool, error) {
	if keyDer, err := CryptoDecodeKey(key, allowEncoding); err != nil {
		return false, err
	} else if key, err := x509.ParsePKIXPublicKey(keyDer); err != nil {
		return false, Errorf(ErrInvalidKeyType, "Cannot deserialize public key: %v", err)
	} else if pubKey, ok := key.(*ecdsa.PublicKey); !ok {
		return false, Errorf(ErrInvalidKeyType, "Unsupported signing algorithm")
	} else {
		hash := sha256.Sum256(msg)
		if ecdsa.VerifyASN1(pubKey, hash[:], ecdsaSign) {
			return true, nil
		}

		return false, Errorf(ErrFailedECDSA, "Verification failed")
	}
}

func CryptoGetKeyFingerprint(key string, allowEncoding KeyEncoding, maxlen int) (string, []byte, error) {
	if derBytes, err := CryptoDecodeKey(key, allowEncoding); err == nil {
		hash := sha256.Sum256(derBytes)

		return string([]rune(hex.EncodeToString(hash[:]))[0:maxlen]), derBytes, nil
	} else {
		return "", nil, err
	}
}

func hashField(value reflect.Value) []byte {
	switch value.Kind() {
	case reflect.String:
		return []byte(value.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(value.Int(), 10))
	case reflect.Float32, reflect.Float64: // 3 digits prec!
		return []byte(strconv.FormatFloat(value.Float(), 'f', 3, 64))
	case reflect.Struct, reflect.Interface:
		if t, ok := value.Interface().(time.Time); ok {
			// Hash Unix time
			return []byte(strconv.FormatInt(t.Unix(), 10))
		}

		return HashStruct(value.Interface())
	case reflect.Ptr:
		if value.IsNil() {
			return []byte{}
		}

		return hashField(value.Elem())
	case reflect.Slice:
		hash := sha256.New()
		for i := 0; i < value.Len(); i++ {
			hash.Write(hashField(value.Index(i)))
		}

		return hash.Sum(nil)
	// Add more cases as needed
	default:
		logs.Error(fmt.Sprintf("Unsupported type %v", value.Interface()))

		return []byte(fmt.Sprintf("%v", value.Interface()))
	}
}

func HashStruct(data interface{}) []byte {
	v := reflect.ValueOf(data)
	t := v.Type()

	type FieldHash struct {
		Tag   string
		Value []byte
	}

	var fields []FieldHash

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("hash")

		if tag == "" {
			continue
		}

		fields = append(fields, FieldHash{Tag: tag, Value: hashField(field)})
	}

	// Sort fields based on the 'hash' tag
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Tag < fields[j].Tag
	})

	// Concatenate the field hashes
	hash := sha256.New()
	for _, field := range fields {
		hash.Write(field.Value)
	}

	return hash.Sum(nil)
}

func HexDigest(obj interface{}) string {
	v := reflect.ValueOf(obj)

	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	hash := hashField(v)

	return string([]rune(hex.EncodeToString(hash)))
}
