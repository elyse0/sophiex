package crypto

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func TestAesSimpleDecryption(t *testing.T) {
	ciphertext, err := hex.DecodeString("e17624bb9e0c4d580f29fbe6edf302dc")
	if err != nil {
		t.Error(err)
	}

	key, err := hex.DecodeString("6d5a7134743677397a24432646294a40")
	if err != nil {
		t.Error(err)
	}

	iv, err := hex.DecodeString("68f6f2484cd62c62601b15ff669ab63a")
	if err != nil {
		t.Error(err)
	}

	plaintext, err := AesDecrypt(ciphertext, key, iv)

	expected := []byte{'t', 'e', 's', 't', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if !reflect.DeepEqual(plaintext, expected) {
		t.Errorf("Decrypted plaintext is %s", string(plaintext))
	}

}
