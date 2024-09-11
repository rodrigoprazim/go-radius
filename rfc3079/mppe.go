package rfc3079

import (
	"crypto/sha1"
	"errors"

	"github.com/rodrigoprazim/go-radius/rfc2759"
)

// KeyLength is the length of keys involved with the functions below
type KeyLength uint

const (
	// KeyLength40Bit - 40-bit
	KeyLength40Bit = KeyLength(8)
	// KeyLength128Bit - 128-bit
	KeyLength128Bit = KeyLength(16)
)

var (
	shaPad1 = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	shaPad2 = []byte{
		0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2,
		0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2,
		0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2,
		0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2,
	}

	magic1 = []byte{
		0x54, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x74,
		0x68, 0x65, 0x20, 0x4d, 0x50, 0x50, 0x45, 0x20, 0x4d,
		0x61, 0x73, 0x74, 0x65, 0x72, 0x20, 0x4b, 0x65, 0x79,
	}

	magic2 = []byte{
		0x4f, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x63, 0x6c, 0x69,
		0x65, 0x6e, 0x74, 0x20, 0x73, 0x69, 0x64, 0x65, 0x2c, 0x20,
		0x74, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x74, 0x68,
		0x65, 0x20, 0x73, 0x65, 0x6e, 0x64, 0x20, 0x6b, 0x65, 0x79,
		0x3b, 0x20, 0x6f, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x73,
		0x65, 0x72, 0x76, 0x65, 0x72, 0x20, 0x73, 0x69, 0x64, 0x65,
		0x2c, 0x20, 0x69, 0x74, 0x20, 0x69, 0x73, 0x20, 0x74, 0x68,
		0x65, 0x20, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x20,
		0x6b, 0x65, 0x79, 0x2e,
	}

	magic3 = []byte{
		0x4f, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x63, 0x6c, 0x69,
		0x65, 0x6e, 0x74, 0x20, 0x73, 0x69, 0x64, 0x65, 0x2c, 0x20,
		0x74, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x74, 0x68,
		0x65, 0x20, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x20,
		0x6b, 0x65, 0x79, 0x3b, 0x20, 0x6f, 0x6e, 0x20, 0x74, 0x68,
		0x65, 0x20, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x20, 0x73,
		0x69, 0x64, 0x65, 0x2c, 0x20, 0x69, 0x74, 0x20, 0x69, 0x73,
		0x20, 0x74, 0x68, 0x65, 0x20, 0x73, 0x65, 0x6e, 0x64, 0x20,
		0x6b, 0x65, 0x79, 0x2e,
	}
)

// GetMasterKey - rfc3079, 3.4
func GetMasterKey(passwordHashHash, ntResponse []byte) []byte {
	sha := sha1.New()
	sha.Write(passwordHashHash)
	sha.Write(ntResponse)
	sha.Write(magic1)
	digest := sha.Sum(nil)

	return digest[:16]
}

// GetAsymmetricStartKey - rfc3079, 3.4
func GetAsymmetricStartKey(masterKey []byte, sessionKeyLength KeyLength, isSend bool) ([]byte, error) {
	if len(masterKey) != 16 {
		return nil, errors.New("masterKey must be 16 bytes long")
	}

	sha := sha1.New()
	sha.Write(masterKey)
	sha.Write(shaPad1)
	if isSend {
		sha.Write(magic3)
	} else {
		sha.Write(magic2)
	}
	sha.Write(shaPad2)
	digest := sha.Sum(nil)

	return digest[:sessionKeyLength], nil
}

// MakeKey - rfc2548, 2.4.2
func MakeKey(ntResponse, password []byte, isSend bool) ([]byte, error) {
	if len(ntResponse) != 24 {
		return nil, errors.New("ntResponse must be 24 bytes in size")
	}

	ucs2Password, err := rfc2759.ToUTF16(password)
	if err != nil {
		return nil, err
	}

	passwordHash := rfc2759.NTPasswordHash(ucs2Password)
	passwordHashHash := rfc2759.NTPasswordHash(passwordHash)
	masterKey := GetMasterKey(passwordHashHash, ntResponse)

	return GetAsymmetricStartKey(masterKey, KeyLength128Bit, isSend)
}
