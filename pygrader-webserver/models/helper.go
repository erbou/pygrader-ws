package models

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"pygrader-webserver/uti"
)

const (
	ErrKidInvalid uti.ErrorCode = 2001 + iota
	ErrKidUnknown
	ErrNoBody
	ErrInvalidInput
	ErrDeadline
	ErrMaxTry
	ErrNotFound
)

type Signature struct {
	Keyid  string `json:"kid"`
	Nonce  string `json:"nonce"`
	Ecdsa  string `json:"ecdas"` // ECDAS Sign(Sha256(Digest))
	Issuer *User
}

type SignedBody struct {
	Signature Signature       `json:"signature"`
	Payload   json.RawMessage `json:"payload"`
}

var GlCache cache.Cache

func init() {
	var err error
	GlCache, err = cache.NewReadThroughCache(
		cache.NewRandomExpireCache(cache.NewMemoryCache()),
		5*time.Minute, // expiration
		// load data from database if the key is absent.
		func(ctx context.Context, key string) (any, error) {
			o := orm.NewOrm()
			u := User{}

			if err := o.RawWithCtx(ctx, "SELECT * FROM user WHERE kid = ?", key).QueryRow(&u); err != nil {
				return nil, err
			}

			// logs.Debug("Cache KID %v", key)
			return &u, nil
		})

	if err != nil {
		logs.Error("Failed to initialize the cache - %v", err)
	}
}

func Unmarshal[E interface{}](data []byte, requireBody bool, e ...*E) (*E, *Signature, error) {
	var h SignedBody
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, nil, err
	}

	if h.Payload == nil {
		if requireBody {
			return nil, nil, uti.Errorf(ErrNoBody, "Payload required")
		}

		return nil, &h.Signature, nil
	}

	var pe *E
	if len(e) > 0 && e[0] != nil {
		pe = e[0]
	} else {
		var ze E
		pe = &ze
	}

	if err := json.Unmarshal(h.Payload, pe); err != nil {
		return nil, nil, err
	}

	return pe, &h.Signature, nil
}

func Verify[E interface{}](data []byte, params ...interface{}) (*E, *Signature, error) {
	var e *E

	var _params []string

	for _, p := range params {
		if s, ok := p.(string); ok {
			_params = append(_params, s)
		} else if i, ok := p.(int64); ok {
			_params = append(_params, strconv.FormatInt(i, 10))
		} else if ee, ok := p.(*E); ok {
			e = ee
		}
	}

	if e, s, err := Unmarshal[E](data, false, e); err != nil {
		logs.Warning("%v", "Parsing error")

		return e, nil, err
	} else {
		if _u, err := GlCache.Get(context.Background(), s.Keyid); err != nil {
			return e, s, uti.Errorf(ErrKidUnknown, "KID not Found %v, %w", s.Keyid, err)
		} else if u, ok := _u.(*User); !ok {
			return e, s, uti.Errorf(uti.ErrSystemError, "System error")
		} else if s.Issuer = u; e == nil {
			return e, s, nil
		} else if ecdsaSign, err := base64.StdEncoding.DecodeString(s.Ecdsa); err != nil {
			return e, s, uti.Errorf(uti.ErrInvalidECDSA, "Cannot decode signature, %w", err)
		} else if keyDer, err := base64.StdEncoding.DecodeString(u.Key); err != nil {
			return e, s, uti.Errorf(uti.ErrInvalidKey, "Cannot decode public key kid %v, %w", s.Keyid, err)
		} else if key, err := x509.ParsePKIXPublicKey(keyDer); err != nil {
			return e, s, uti.Errorf(uti.ErrInvalidKey, "Cannot deserialize public key %v, %w", s.Keyid, err)
		} else if pubKey, ok := key.(*ecdsa.PublicKey); !ok {
			return e, s, uti.Errorf(uti.ErrInvalidKeyType, "Key %v, Unsupported algorithm", s.Keyid)
		} else {
			h := strings.Join(append(_params, s.Nonce, uti.HexDigest(e)), ":")
			hash := sha256.Sum256([]byte(h))
			// logs.Debug("%v -> %v", h, hex.EncodeToString(hash[:]))
			if ecdsa.VerifyASN1(pubKey, hash[:], ecdsaSign) {
				return e, s, nil
			}

			return e, s, uti.Errorf(uti.ErrFailedECDSA, "Verification failed with key %s", s.Keyid)
		}
	}
}
