package crypto

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/stretchr/testify/assert"
)

func TestCryptoSmoke(t *testing.T) {
	c := SimpleSymmetricCryptoRelay{}

	data := []byte("hello world")
	key := []byte("6EMHrPKuHIVaM1b5ss2MfahtoidxbGJ4")

	eData, err := c.encrypt(key, data)
	assert.Nil(t, err)

	dData, err := c.decrypt(key, eData)
	assert.Nil(t, err)

	assert.Equal(t, data, dData)
	assert.NotEqual(t, data, eData)
}

func TestCryptoDirectionSmoke(t *testing.T) {
	c := SimpleSymmetricCryptoRelay{}

	key := []byte("6EMHrPKuHIVaM1b6ss2MfahtoidxbGJ4")
	strData := "hello world"
	eStr, err := c.applyEncryptionStr(encryptDirection, key, strData)
	assert.Nil(t, err)

	dStr, err := c.applyEncryptionStr(decryptDirection, key, eStr)
	assert.Nil(t, err)

	assert.Equal(t, strData, dStr)
	assert.NotEqual(t, strData, eStr)
}

func TestCryptoErrIfKeyWrong(t *testing.T) {
	c := SimpleSymmetricCryptoRelay{}

	key1 := []byte("6EMHrPKuHIVaM1b6ss2MfahtoidxbGJ4")
	key2 := []byte("777HrPKuHIVaM1b6ss2MfahtoidxbGJ4")
	strData := "hello world"
	eStr, err := c.applyEncryptionStr(encryptDirection, key1, strData)
	assert.Nil(t, err)

	dStr, err := c.applyEncryptionStr(decryptDirection, key2, eStr)
	// cipher: message authentication failed
	assert.Error(t, err)
	assert.Equal(t, "", dStr)
}

func TestSmoke(t *testing.T) {
	eN := SimpleSymmetricCryptoRelay{
		Mode:     "encrypt",
		Password: "N7X92q5R2CFuP6utEZMrzsaJdDjECXwt",
	}

	dN := SimpleSymmetricCryptoRelay{
		Mode:     "decrypt",
		Password: "N7X92q5R2CFuP6utEZMrzsaJdDjECXwt",
	}

	ch1 := make(chan []events.Event)
	ch2 := make(chan []events.Event)

	results := []events.Event{}
	ch3 := make(chan []events.Event)
	go func() {
		for {
			e := <-ch3
			results = append(results, e...)
		}
	}()

	go func() {
		err := eN.RDo(context.TODO(), ch1, ch2)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := dN.RDo(context.TODO(), ch2, ch3)
		if err != nil {
			panic(err)
		}
	}()

	numEvents := 100
	for i := 0; i < numEvents; i++ {
		count := fmt.Sprintf("%d", i)
		ch1 <- []events.Event{events.EventMsg{
			ID: count,
			Headers: map[string][]string{
				count: []string{count},
			},
			URI:     fmt.Sprintf("/%s", count),
			Content: []byte(count),
		}}
	}

	// messages can take some time.
	time.Sleep(time.Second * 5)

	assert.Equal(t, numEvents, len(results))
	assert.Equal(t, "/0", results[0].GetURI())
	assert.Equal(t, "0", string(results[0].GetContent()))
}
