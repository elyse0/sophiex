package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

func AesDecrypt(ciphertext []byte, key []byte, iv []byte) (plaintext []byte, err error) {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext = make([]byte, len(ciphertext))

	aesCbcMode := cipher.NewCBCDecrypter(aesCipher, iv)
	aesCbcMode.CryptBlocks(plaintext, ciphertext)

	return plaintext, nil
}
