package mcoej_test

import (
	"crypto/aes"
	"crypto/cipher"
	"testing"

	"github.com/runningwild/skein/mcoej"
	tf512 "github.com/runningwild/skein/threefish/512"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSkein1024(t *testing.T) {
	Convey("mcoej", t, func() {
		mc, err := mcoej.New(tf512.Encrypt, tf512.Decrypt, 512)
		So(err, ShouldBeNil)
		So(mc, ShouldNotBeNil)
		key := [64]byte{0, 1, 2, 3, 4, 5}
		nonce := [64]byte{5, 4, 3, 2, 1}
		publicData := []byte("this is public data, rawr")
		plaintext := []byte("buttons so adorable, so adorable i could scream!  Like seriously, the most adorbs evah!")
		ciphertext := mc.Lock(key[:], nonce[:], publicData, plaintext, nil)
		output, err := mc.Unlock(key[:], nonce[:], publicData, ciphertext, nil)
		So(err, ShouldBeNil)
		So(output, ShouldResemble, plaintext)
	})
}

var (
	benchmarkPlaintext1M []byte
)

func init() {
	benchmarkPlaintext1M = make([]byte, 1024*1024)
	for i := range benchmarkPlaintext1M {
		benchmarkPlaintext1M[i] = byte(i)
	}
}

func BenchmarkMcOEJLock(b *testing.B) {
	b.StopTimer()
	mc, _ := mcoej.New(tf512.Encrypt, tf512.Decrypt, 512)
	key := [64]byte{0, 1, 2, 3, 4, 5}
	nonce := [64]byte{5, 4, 3, 2, 1}
	publicData := []byte("this is public data, rawr")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mc.Lock(key[:], nonce[:], publicData, benchmarkPlaintext1M, nil)
	}
}

func BenchmarkMcOEJUnlock(b *testing.B) {
	b.StopTimer()
	mc, _ := mcoej.New(tf512.Encrypt, tf512.Decrypt, 512)
	key := [64]byte{0, 1, 2, 3, 4, 5}
	nonce := [64]byte{5, 4, 3, 2, 1}
	publicData := []byte("this is public data, rawr")
	ciphertext := mc.Lock(key[:], nonce[:], publicData, benchmarkPlaintext1M, nil)
	if _, err := mc.Unlock(key[:], nonce[:], publicData, ciphertext, nil); err != nil {
		panic(err)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mc.Unlock(key[:], nonce[:], publicData, ciphertext, nil)
	}
}

func BenchmarkAESCBC(b *testing.B) {
	b.StopTimer()
	c, err := aes.NewCipher(make([]byte, 16))
	if err != nil {
		panic(err)
	}
	en := cipher.NewCBCEncrypter(c, make([]byte, 16))
	data := make([]byte, 1024*1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		en.CryptBlocks(data, data)
	}
}

func BenchmarkAESGCMSeal(b *testing.B) {
	b.StopTimer()
	c, err := aes.NewCipher(make([]byte, 16))
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		panic(err)
	}
	nonce := [12]byte{5, 4, 3, 2, 1}
	publicData := []byte("this is public data, rawr")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		gcm.Seal(nil, nonce[:], benchmarkPlaintext1M, publicData)
	}
}
