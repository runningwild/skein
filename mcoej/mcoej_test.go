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
		Convey("decrypt should reverse encryption", func() {
			output, err := mc.Unlock(key[:], nonce[:], publicData, ciphertext, nil)
			So(err, ShouldBeNil)
			So(output, ShouldResemble, plaintext)
		})
		Convey("decrypt should fail if the public data has been changed", func() {
			publicData[0]++
			_, err := mc.Unlock(key[:], nonce[:], publicData, ciphertext, nil)
			publicData[0]--
			So(err, ShouldNotBeNil)
		})
		Convey("decrypt should fail if the ciphertext has been changed", func() {
			ciphertext[0]++
			_, err := mc.Unlock(key[:], nonce[:], publicData, ciphertext, nil)
			ciphertext[0]--
			So(err, ShouldNotBeNil)
		})
		Convey("decrypt should fail if the wrong key is used", func() {
			key[0]++
			_, err := mc.Unlock(key[:], nonce[:], publicData, ciphertext, nil)
			key[0]--
			So(err, ShouldNotBeNil)
		})
		Convey("decrypt should fail if the wrong nonce is used", func() {
			nonce[0]++
			_, err := mc.Unlock(key[:], nonce[:], publicData, ciphertext, nil)
			nonce[0]--
			So(err, ShouldNotBeNil)
		})
	})
}

var (
	benchmarkPublicData   []byte
	benchmarkPlaintext16b []byte
	benchmarkPlaintext1k  []byte
	benchmarkPlaintext1M  []byte
)

func init() {
	benchmarkPublicData = []byte("this is public data, rawr       ") // 32 bytes
	benchmarkPlaintext16b = make([]byte, 16)
	for i := range benchmarkPlaintext16b {
		benchmarkPlaintext16b[i] = byte(i)
	}
	benchmarkPlaintext1k = make([]byte, 1024)
	for i := range benchmarkPlaintext1k {
		benchmarkPlaintext1k[i] = byte(i)
	}
	benchmarkPlaintext1M = make([]byte, 1024*1024)
	for i := range benchmarkPlaintext1M {
		benchmarkPlaintext1M[i] = byte(i)
	}
}

func BenchmarkMcOEJLock_10b(b *testing.B) {
	b.StopTimer()
	mc, _ := mcoej.New(tf512.Encrypt, tf512.Decrypt, 512)
	key := [64]byte{0, 1, 2, 3, 4, 5}
	nonce := [64]byte{5, 4, 3, 2, 1}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mc.Lock(key[:], nonce[:], benchmarkPublicData, benchmarkPlaintext16b, nil)
	}
}

func BenchmarkMcOEJUnlock_10b(b *testing.B) {
	b.StopTimer()
	mc, _ := mcoej.New(tf512.Encrypt, tf512.Decrypt, 512)
	key := [64]byte{0, 1, 2, 3, 4, 5}
	nonce := [64]byte{5, 4, 3, 2, 1}

	ciphertext := mc.Lock(key[:], nonce[:], benchmarkPublicData, benchmarkPlaintext16b, nil)
	if _, err := mc.Unlock(key[:], nonce[:], benchmarkPublicData, ciphertext, nil); err != nil {
		panic(err)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mc.Unlock(key[:], nonce[:], benchmarkPublicData, ciphertext, nil)
	}
}

func BenchmarkAESCBC_10b(b *testing.B) {
	b.StopTimer()
	c, err := aes.NewCipher(make([]byte, 16))
	if err != nil {
		panic(err)
	}
	en := cipher.NewCBCEncrypter(c, make([]byte, 16))
	data := make([]byte, len(benchmarkPublicData)+len(benchmarkPlaintext16b))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		en.CryptBlocks(data, data)
	}
}

func BenchmarkAESGCMSeal_10b(b *testing.B) {
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

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		gcm.Seal(nil, nonce[:], benchmarkPlaintext16b, benchmarkPublicData)
	}
}

func BenchmarkMcOEJLock_1k(b *testing.B) {
	b.StopTimer()
	mc, _ := mcoej.New(tf512.Encrypt, tf512.Decrypt, 512)
	key := [64]byte{0, 1, 2, 3, 4, 5}
	nonce := [64]byte{5, 4, 3, 2, 1}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mc.Lock(key[:], nonce[:], benchmarkPublicData, benchmarkPlaintext1k, nil)
	}
}

func BenchmarkMcOEJUnlock_1k(b *testing.B) {
	b.StopTimer()
	mc, _ := mcoej.New(tf512.Encrypt, tf512.Decrypt, 512)
	key := [64]byte{0, 1, 2, 3, 4, 5}
	nonce := [64]byte{5, 4, 3, 2, 1}

	ciphertext := mc.Lock(key[:], nonce[:], benchmarkPublicData, benchmarkPlaintext1k, nil)
	if _, err := mc.Unlock(key[:], nonce[:], benchmarkPublicData, ciphertext, nil); err != nil {
		panic(err)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mc.Unlock(key[:], nonce[:], benchmarkPublicData, ciphertext, nil)
	}
}

func BenchmarkAESCBC_1k(b *testing.B) {
	b.StopTimer()
	c, err := aes.NewCipher(make([]byte, 16))
	if err != nil {
		panic(err)
	}
	en := cipher.NewCBCEncrypter(c, make([]byte, 16))
	data := make([]byte, len(benchmarkPublicData)+len(benchmarkPlaintext1k))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		en.CryptBlocks(data, data)
	}
}

func BenchmarkAESGCMSeal_1k(b *testing.B) {
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

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		gcm.Seal(nil, nonce[:], benchmarkPlaintext1k, benchmarkPublicData)
	}
}

func BenchmarkMcOEJLock_1M(b *testing.B) {
	b.StopTimer()
	mc, _ := mcoej.New(tf512.Encrypt, tf512.Decrypt, 512)
	key := [64]byte{0, 1, 2, 3, 4, 5}
	nonce := [64]byte{5, 4, 3, 2, 1}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mc.Lock(key[:], nonce[:], benchmarkPublicData, benchmarkPlaintext1M, nil)
	}
}

func BenchmarkMcOEJUnlock_1M(b *testing.B) {
	b.StopTimer()
	mc, _ := mcoej.New(tf512.Encrypt, tf512.Decrypt, 512)
	key := [64]byte{0, 1, 2, 3, 4, 5}
	nonce := [64]byte{5, 4, 3, 2, 1}

	ciphertext := mc.Lock(key[:], nonce[:], benchmarkPublicData, benchmarkPlaintext1M, nil)
	if _, err := mc.Unlock(key[:], nonce[:], benchmarkPublicData, ciphertext, nil); err != nil {
		panic(err)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mc.Unlock(key[:], nonce[:], benchmarkPublicData, ciphertext, nil)
	}
}

func BenchmarkAESCBC_1M(b *testing.B) {
	b.StopTimer()
	c, err := aes.NewCipher(make([]byte, 16))
	if err != nil {
		panic(err)
	}
	en := cipher.NewCBCEncrypter(c, make([]byte, 16))
	data := make([]byte, len(benchmarkPublicData)+len(benchmarkPlaintext1M))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		en.CryptBlocks(data, data)
	}
}

func BenchmarkAESGCMSeal_1M(b *testing.B) {
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

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		gcm.Seal(nil, nonce[:], benchmarkPlaintext1M, benchmarkPublicData)
	}
}
