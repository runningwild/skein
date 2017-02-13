package mcoej_test

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"testing"

	"github.com/runningwild/skein/mcoej"
	tf1024 "github.com/runningwild/skein/threefish/1024"
	tf256 "github.com/runningwild/skein/threefish/256"
	tf512 "github.com/runningwild/skein/threefish/512"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMcOEJ(t *testing.T) {
	Convey("mcoej", t, func() {
		tests := []struct {
			enc, dec mcoej.TweakableBlockCipher
			size     int
		}{
			{tf256.Encrypt, tf256.Decrypt, 256},
			{tf512.Encrypt, tf512.Decrypt, 512},
			{tf1024.Encrypt, tf1024.Decrypt, 1024},
		}
		datas := [][]byte{
			[]byte(""),
			[]byte("Data"),
			[]byte("lessthan128bits"),
			[]byte("exactly 128 bits"),
			[]byte("morethan 128 bits"),
			[]byte("less than 256 bits-------------"),
			[]byte("exactly 256 bits----------------"),
			[]byte("more than 256 bits---------------"),
			[]byte("This is a bunch of data, just to make sure it is longer than the other things."),
		}
		for _, test := range tests {
			Convey(fmt.Sprintf("with tf %d", test.size), func() {
				mc, err := mcoej.New(test.enc, test.dec, test.size)
				So(err, ShouldBeNil)
				So(mc, ShouldNotBeNil)
				key := make([]byte, test.size/8)
				nonce := make([]byte, test.size/8)
				key[0], key[1], key[2] = 10, 20, 30
				nonce[0], nonce[1], nonce[2] = 40, 50, 60
				for _, publicData := range datas {
					for _, plaintext := range datas {
						Convey(fmt.Sprintf("public data: %q, plaintext: %q", publicData, plaintext), func() {
							ciphertext := mc.Lock(key, nonce, []byte(publicData), []byte(plaintext), nil)
							Convey("decrypt should reverse encryption", func() {
								output, err := mc.Unlock(key, nonce, publicData, ciphertext, nil)
								So(err, ShouldBeNil)
								So(output, ShouldResemble, plaintext)
							})
							Convey("decrypt should fail if the public data has been changed", func() {
								publicData = append(publicData, 0)
								_, err := mc.Unlock(key, nonce, publicData, ciphertext, nil)
								publicData = publicData[0 : len(publicData)-1]
								So(err, ShouldNotBeNil)
							})
							Convey("decrypt should fail if the ciphertext has been changed", func() {
								if len(ciphertext) == 0 {
									ciphertext = append(ciphertext, 0)
								} else {
									ciphertext[0]++
								}
								ciphertext = append(ciphertext, 0)
								_, err := mc.Unlock(key, nonce, publicData, ciphertext, nil)
								if len(ciphertext) == 0 {
									ciphertext = ciphertext[0 : len(ciphertext)-1]
								} else {
									ciphertext[0]--
								}
								So(err, ShouldNotBeNil)
							})
							Convey("decrypt should fail if the wrong key is used", func() {
								key[0]++
								_, err := mc.Unlock(key, nonce, publicData, ciphertext, nil)
								key[0]--
								So(err, ShouldNotBeNil)
							})
							Convey("decrypt should fail if the wrong nonce is used", func() {
								nonce[0]++
								_, err := mc.Unlock(key, nonce, publicData, ciphertext, nil)
								nonce[0]--
								So(err, ShouldNotBeNil)
							})
						})
					}
				}
			})
		}
	})
}

func TestMcOEJ512(t *testing.T) {
	Convey("mcoej 512", t, func() {
		datas := [][]byte{
			[]byte(""),
			[]byte("Data"),
			[]byte("lessthan128bits"),
			[]byte("exactly 128 bits"),
			[]byte("morethan 128 bits"),
			[]byte("less than 256 bits-------------"),
			[]byte("exactly 256 bits----------------"),
			[]byte("more than 256 bits---------------"),
			[]byte("This is a bunch of data, just to make sure it is longer than the other things."),
		}
		mc := mcoej.New512()
		So(mc, ShouldNotBeNil)
		key := make([]byte, 64)
		nonce := make([]byte, 64)
		key[0], key[1], key[2] = 10, 20, 30
		nonce[0], nonce[1], nonce[2] = 40, 50, 60
		for _, publicData := range datas {
			for _, plaintext := range datas {
				Convey(fmt.Sprintf("public data: %q, plaintext: %q", publicData, plaintext), func() {
					ciphertext := mc.Lock(key, nonce, publicData, plaintext, nil)
					Convey("decrypt should reverse encryption", func() {
						output, err := mc.Unlock(key, nonce, publicData, ciphertext, nil)
						So(err, ShouldBeNil)
						So(output, ShouldResemble, plaintext)
					})
					Convey("decrypt should fail if the public data has been changed", func() {
						publicData = append(publicData, 0)
						_, err := mc.Unlock(key, nonce, publicData, ciphertext, nil)
						publicData = publicData[0 : len(publicData)-1]
						So(err, ShouldNotBeNil)
					})
					Convey("decrypt should fail if the ciphertext has been changed", func() {
						if len(ciphertext) == 0 {
							ciphertext = append(ciphertext, 0)
						} else {
							ciphertext[0]++
						}
						ciphertext = append(ciphertext, 0)
						_, err := mc.Unlock(key, nonce, publicData, ciphertext, nil)
						if len(ciphertext) == 0 {
							ciphertext = ciphertext[0 : len(ciphertext)-1]
						} else {
							ciphertext[0]--
						}
						So(err, ShouldNotBeNil)
					})
					Convey("decrypt should fail if the wrong key is used", func() {
						key[0]++
						_, err := mc.Unlock(key, nonce, publicData, ciphertext, nil)
						key[0]--
						So(err, ShouldNotBeNil)
					})
					Convey("decrypt should fail if the wrong nonce is used", func() {
						nonce[0]++
						_, err := mc.Unlock(key, nonce, publicData, ciphertext, nil)
						nonce[0]--
						So(err, ShouldNotBeNil)
					})
					Convey("ciphertext should match ciphertext from general McOEJ object", func() {
						mcgen, _ := mcoej.New(tf512.Encrypt, tf512.Decrypt, 512)
						ciphertext2 := mcgen.Lock(key, nonce, publicData, plaintext, nil)
						So(ciphertext, ShouldResemble, ciphertext2)
					})
				})
			}
		}
	})
}

var (
	benchmarkPublicData   []byte
	benchmarkPlaintext32b []byte
	benchmarkPlaintext1k  []byte
	benchmarkPlaintext1M  []byte
)

func init() {
	benchmarkPublicData = []byte("this is public data, rawr       ") // 32 bytes
	benchmarkPlaintext32b = make([]byte, 32)
	for i := range benchmarkPlaintext32b {
		benchmarkPlaintext32b[i] = byte(i)
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

func BenchmarkMcOEJLock_32b(b *testing.B) {
	b.StopTimer()
	mc, _ := mcoej.New(tf512.Encrypt, tf512.Decrypt, 512)
	key := [64]byte{0, 1, 2, 3, 4, 5}
	nonce := [64]byte{5, 4, 3, 2, 1}
	dst := make([]byte, 1000)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mc.Lock(key[:], nonce[:], benchmarkPublicData, benchmarkPlaintext32b, dst[0:0])
	}
}

func BenchmarkMcOEJUnlock_32b(b *testing.B) {
	b.StopTimer()
	mc, _ := mcoej.New(tf512.Encrypt, tf512.Decrypt, 512)
	key := [64]byte{0, 1, 2, 3, 4, 5}
	nonce := [64]byte{5, 4, 3, 2, 1}

	ciphertext := mc.Lock(key[:], nonce[:], benchmarkPublicData, benchmarkPlaintext32b, nil)
	if _, err := mc.Unlock(key[:], nonce[:], benchmarkPublicData, ciphertext, nil); err != nil {
		panic(err)
	}
	dst := make([]byte, 1000)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mc.Unlock(key[:], nonce[:], benchmarkPublicData, ciphertext, dst[0:0])
	}
}

func BenchmarkMcOEJ512Lock_32b(b *testing.B) {
	b.StopTimer()
	mc := mcoej.New512()
	key := [64]byte{0, 1, 2, 3, 4, 5}
	nonce := [64]byte{5, 4, 3, 2, 1}
	dst := make([]byte, 1000)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mc.Lock(key[:], nonce[:], benchmarkPublicData, benchmarkPlaintext32b, dst[0:0])
	}
}

func BenchmarkMcOEJ512Unlock_32b(b *testing.B) {
	b.StopTimer()
	mc := mcoej.New512()
	key := [64]byte{0, 1, 2, 3, 4, 5}
	nonce := [64]byte{5, 4, 3, 2, 1}

	ciphertext := mc.Lock(key[:], nonce[:], benchmarkPublicData, benchmarkPlaintext32b, nil)
	if _, err := mc.Unlock(key[:], nonce[:], benchmarkPublicData, ciphertext, nil); err != nil {
		panic(err)
	}
	dst := make([]byte, 1000)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mc.Unlock(key[:], nonce[:], benchmarkPublicData, ciphertext, dst[0:0])
	}
}

func BenchmarkThreefish512_CBC_64b(b *testing.B) {
	b.StopTimer()
	c := tf512.MakeCipher([64]byte{})
	en := cipher.NewCBCEncrypter(c, make([]byte, 64))
	data := make([]byte, 64)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		en.CryptBlocks(data, data)
	}
}

func BenchmarkAESCBC_32b(b *testing.B) {
	b.StopTimer()
	c, err := aes.NewCipher(make([]byte, 16))
	if err != nil {
		panic(err)
	}
	en := cipher.NewCBCEncrypter(c, make([]byte, 16))
	data := make([]byte, len(benchmarkPublicData)+len(benchmarkPlaintext32b))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		en.CryptBlocks(data, data)
	}
}

func BenchmarkAESGCMSeal_32b(b *testing.B) {
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
		gcm.Seal(nil, nonce[:], benchmarkPlaintext32b, benchmarkPublicData)
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

func BenchmarkMcOEJ512Lock_1k(b *testing.B) {
	b.StopTimer()
	mc := mcoej.New512()
	key := [64]byte{0, 1, 2, 3, 4, 5}
	nonce := [64]byte{5, 4, 3, 2, 1}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mc.Lock(key[:], nonce[:], benchmarkPublicData, benchmarkPlaintext1k, nil)
	}
}

func BenchmarkMcOEJ512Unlock_1k(b *testing.B) {
	b.StopTimer()
	mc := mcoej.New512()
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

func BenchmarkThreefish512_CBC_1k(b *testing.B) {
	b.StopTimer()
	c := tf512.MakeCipher([64]byte{})
	en := cipher.NewCBCEncrypter(c, make([]byte, 64))
	data := make([]byte, 1024+64)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		en.CryptBlocks(data, data)
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
	var dst []byte
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		dst = mc.Lock(key[:], nonce[:], benchmarkPublicData, benchmarkPlaintext1M, dst[0:0])
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

func BenchmarkMcOEJ512Lock_1M(b *testing.B) {
	b.StopTimer()
	mc := mcoej.New512()
	key := [64]byte{0, 1, 2, 3, 4, 5}
	nonce := [64]byte{5, 4, 3, 2, 1}
	var dst []byte
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		dst = mc.Lock(key[:], nonce[:], benchmarkPublicData, benchmarkPlaintext1M, dst[0:0])
	}
}

func BenchmarkMcOEJ512Unlock_1M(b *testing.B) {
	b.StopTimer()
	mc := mcoej.New512()
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

func BenchmarkThreefish512_CBC_1M(b *testing.B) {
	b.StopTimer()
	c := tf512.MakeCipher([64]byte{})
	en := cipher.NewCBCEncrypter(c, make([]byte, 64))
	data := make([]byte, 1024*1024+64)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		en.CryptBlocks(data, data)
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
