package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"io"
)

// Interface to be implemented for an encryption and encoding scheme
type PasswordTranscoder interface {
	// DecodePassword takes an encoded password as an input, and decode if then decypts it.
	DecodePassword(pass string) ([]byte, error)
	// EncodePassword takes a password as an input, and encrypts then encode it to printable characters. If an error occurs, it should raise it.
	EncodePassword(pass string) ([]byte, error)
}

type transcoder struct {
	passphrase string
	*base64.Encoding
}

func NewTranscoder(s string) PasswordTranscoder {
	return &transcoder{s, base64.StdEncoding}
}

func (d *transcoder) DecodePassword(pass string) ([]byte, error) {
	ciphertext, err := d.DecodeString(pass)
	if err != nil {
		return nil, err
	}

	key := sha512.Sum512_256([]byte(d.passphrase))

	block, _ := aes.NewCipher(key[:])
	if len(ciphertext) < aes.BlockSize {
		panic("No IV in the ciphertext !!!")
	}

	iv, ciphertext := ciphertext[:aes.BlockSize], ciphertext[aes.BlockSize:]
	stream := cipher.NewCTR(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

func (e *transcoder) EncodePassword(pass string) ([]byte, error) {
	key := sha512.Sum512_256([]byte(e.passphrase))

	block, _ := aes.NewCipher(key[:])

	ciphertext := make([]byte, aes.BlockSize+len(pass))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(pass))

	res := make([]byte, e.EncodedLen(len(ciphertext)))
	e.Encode(res, ciphertext)

	return res, nil
}
