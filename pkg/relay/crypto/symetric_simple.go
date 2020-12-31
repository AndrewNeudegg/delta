package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/relay"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
)

type encryptionDirection int

const (
	encryptDirection encryptionDirection = iota
	decryptDirection encryptionDirection = iota
)

// SimpleSymmetricCryptoRelay is a reference implementation of encryption in the relay flow.
// Note that I am not a cryptographer and at the time of writing these methods / techniques have not
// been verified as safe.
// They have been copied from the golang AES tests at
//   https://golang.org/src/crypto/cipher/example_test.go
// and while they might be assumed to be safe, the author(s) do not guarantee safety.
// Your millage may vary.
type SimpleSymmetricCryptoRelay struct {
	relay.R

	Mode         string              // Mode encrypt/decrypt
	Password     string              // Password is hardcoded in the config.
	EnvVar       string              // EnvVar overrides Password with an environment variable.
	HashPassword bool                // HashPassword will sha1 hash the password.
	key          []byte              // key is the container for the password Todo: something more sensible with this.
	direction    encryptionDirection // specifies if this is encryption or decryption.
}

func (r SimpleSymmetricCryptoRelay) actionMap(m map[string][]string) (map[string][]string, error) {
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

func (r SimpleSymmetricCryptoRelay) actionString(s string) (string, error) {
	return r.applyEncryptionStr(r.direction, r.key, s)
}

func (r SimpleSymmetricCryptoRelay) actionContent(c []byte) ([]byte, error) {
	return r.applyEncryptionBytes(r.direction, r.key, c)
}

func (r SimpleSymmetricCryptoRelay) applyEncryptionStr(direction encryptionDirection, key []byte, data string) (string, error) {
	if direction == encryptDirection {
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

func (r SimpleSymmetricCryptoRelay) applyEncryptionBytes(direction encryptionDirection, key []byte, data []byte) ([]byte, error) {
	if direction == encryptDirection {
		return r.encrypt(key, data)
	}

	return r.decrypt(key, data)
}

// encrypt will generate a nonce and prepend it to the cipherdata after encryption.
func (r SimpleSymmetricCryptoRelay) encrypt(key []byte, data []byte) ([]byte, error) {
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
func (r SimpleSymmetricCryptoRelay) decrypt(key []byte, data []byte) ([]byte, error) {
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

// Do will pass messages through an intermediary that may perform operations on the data.
func (r *SimpleSymmetricCryptoRelay) Do(ctx context.Context, outbound <-chan events.Event, inbound chan<- events.Event) error {

	if r.EnvVar != "" {
		r.Password = os.Getenv(r.EnvVar)
		if r.Password == "" {
			return fmt.Errorf("expected to populate password from environment variable named '%s' but was empty", r.EnvVar)
		}
	}

	key := []byte{}
	if r.HashPassword {
		hasher := sha512.New()
		hasher.Write(key)
		key = hasher.Sum(nil)
		key = key[:32]
	} else {
		key = []byte(r.Password)
	}

	if len(key) != 32 {
		return fmt.Errorf("key is incorrect size (%d != 32)", len(key))
	}
	r.key = key

	switch r.Mode {
	case "encrypt":
		r.direction = encryptDirection
	case "decrypt":
		r.direction = decryptDirection
	default:
		return fmt.Errorf("SimpleSymmetricCryptoRelay mode '%s' is unknown", r.Mode)
	}

	// Pass all messages from the outbound queue to the inbound queue.
	for ctx.Err() == nil {
		select {
		case e := <-outbound:
			mHeaders, err := r.actionMap(e.GetHeaders())
			if err != nil {
				log.Error(errors.Wrapf(err, "failed to '%s' headers for event '%s'", r.Mode, e.GetMessageID()))
			}

			mURI, err := r.actionString(e.GetURI())
			if err != nil {
				log.Error(errors.Wrapf(err, "failed to '%s' URI for event '%s'", r.Mode, e.GetMessageID()))
			}

			mContent, err := r.actionContent(e.GetContent())
			if err != nil {
				log.Error(errors.Wrapf(err, "failed to '%s' content for event '%s'", r.Mode, e.GetMessageID()))
			}

			inbound <- events.EventMsg{
				ID:      e.GetMessageID(),
				Headers: mHeaders,
				URI:     mURI,
				Content: mContent,
			}
		case _ = <-ctx.Done():
			break
		}
	}
	return ctx.Err()
}
