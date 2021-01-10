package crypto1

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/andrewneudegg/delta/pkg/events"
	log "github.com/sirupsen/logrus"
)

const (
	// ID for this collection of resources.
	ID = "utilities/crypto/v1"

	// EncryptMode specifies that this resource will perfrom encryption.
	EncryptMode = "encrypt"
	// DecryptMode specifies that this resource will perform decryption.
	DecryptMode = "decrypt"
)

func doCrypto(eCol events.Collection, c *crypto) (events.Collection, error) {
	log.Debugf("doCrypto on '%d' events, direction '%s'", len(eCol), c.direction)

	rCol := make(events.Collection, len(eCol))

	for i, e := range eCol {
		headers, err := c.actionMap(e.GetHeaders())
		if err != nil {
			return nil, err
		}

		uri, err := c.actionString(e.GetURI())
		if err != nil {
			return nil, err
		}

		cont, err := c.actionContent(e.GetContent())
		if err != nil {
			return nil, err
		}

		rCol[i] = events.EventMsg{
			ID:      ID,
			Headers: headers,
			URI:     uri,
			Content: cont,
		}
	}

	return rCol, nil
}

type crypto struct {
	key       []byte // key is the container for the password Todo: something more sensible with this.
	direction string // specifies if this is encryption or decryption.
}

func (r crypto) actionMap(m map[string][]string) (map[string][]string, error) {
	resultantMap := make(map[string][]string)

	for k, v := range m {
		mK, err := r.actionString(k)
		if err != nil {
			return map[string][]string{}, err
		}
		resultantMap[mK] = []string{}
		for _, s := range v {
			dV, err := r.actionString(s)
			if err != nil {
				return map[string][]string{}, err
			}
			resultantMap[mK] = append(resultantMap[mK], dV)
		}
	}

	return resultantMap, nil
}

func (r crypto) actionString(s string) (string, error) {
	return r.applyEncryptionStr(r.direction, r.key, s)
}

func (r crypto) actionContent(c []byte) ([]byte, error) {
	return r.applyEncryptionBytes(r.direction, r.key, c)
}

func (r crypto) applyEncryptionStr(direction string, key []byte, data string) (string, error) {
	if direction == EncryptMode {
		eBytes, err := r.encrypt(key, []byte(data))
		return hex.EncodeToString(eBytes), err
	}

	eBytes, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}

	bytes, err := r.decrypt(key, eBytes)
	return string(bytes), err
}

func (r crypto) applyEncryptionBytes(direction string, key []byte, data []byte) ([]byte, error) {
	if direction == EncryptMode {
		return r.encrypt(key, data)
	}

	return r.decrypt(key, data)
}

// encrypt will generate a nonce and prepend it to the cipherdata after encryption.
func (r crypto) encrypt(key []byte, data []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key is incorrect size (%d != 32)", len(key))
	}

	// https://golang.org/src/crypto/cipher/example_test.go
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, data, nil)
	ciphertext = append(nonce, ciphertext...)
	return ciphertext, nil
}

// decrypt will extract the nonce and return the decrypted data.
func (r crypto) decrypt(key []byte, data []byte) ([]byte, error) {
	// https://golang.org/src/crypto/cipher/example_test.go
	ciphertext := data[12:]
	nonce := data[:12]
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
