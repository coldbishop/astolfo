package main

import (
	"fmt"
	"unicode/utf8"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/blake2b"
)

const (
	iter       = 1 << 4
	memorySize = 1 << 15
	thread     = 1 << 2
	keySize    = 1 << 6

	algorithmID = "com.github.coldbishop.astolfo"
)

func algorithm(b []byte, passLength uint8, upper, lower, punct, digit bool) string {
	// the magic number 94 is the total amount of printable ASCII characters
	// minus the "space" character
	s := make([]byte, 0, 94)
	includeChars(&s, upper, lower, punct, digit)

	lenS := len(s)
	lenB := len(b)
	divS := lenB / int(passLength)
	modS := lenB % int(passLength)
	var generatedPassword string
	var i, j int

	for {
		acc := 0
		k := j + divS
		if modS != 0 {
			k += 1
			modS -= 1
		}
		for ; i < k; i++ {
			acc += int(b[i])
		}
		generatedPassword += string(s[(acc % lenS)])
		if i == lenB {
			break
		}
		j = i
	}
	return generatedPassword
}

func generateKey(userName string, password []byte) []byte {
	salt := []byte(fmt.Sprintf("%s%d%s", algorithmID, utf8.RuneCountInString(userName), userName))
	return argon2.Key(password, salt, iter, memorySize, thread, keySize)
}

func generatePassword(userName, siteName string, password []byte, passLength uint8, counter uint, upper, lower, punct, digit bool) (string, error) {
	key := generateKey(userName, password)
	seed, err := generateSeed(key, siteName, counter)
	if err != nil {
		return "", err
	}

	return algorithm(seed, passLength, upper, lower, punct, digit), nil
}

func generateSeed(key []byte, siteName string, counter uint) ([]byte, error) {
	hash, err := blake2b.New512(key)
	if err != nil {
		return nil, err
	}
	siteNamePlusCounter := []byte(fmt.Sprintf("%s%d%s%d", algorithmID, utf8.RuneCountInString(siteName), siteName, counter))
	hash.Write(siteNamePlusCounter)

	return hash.Sum(nil), nil
}

func includeChars(b *[]byte, upper, lower, digit, punct bool) {
	if upper {
		for i := 65; i <= 90; i++ {
			*b = append(*b, byte(i))
		}
	}
	if lower {
		for i := 97; i <= 122; i++ {
			*b = append(*b, byte(i))
		}
	}
	if digit {
		for i := 48; i <= 57; i++ {
			*b = append(*b, byte(i))
		}
	}
	if punct {
		for i := 33; i <= 47; i++ {
			*b = append(*b, byte(i))
		}
		for i := 58; i <= 64; i++ {
			*b = append(*b, byte(i))
		}
		for i := 91; i <= 96; i++ {
			*b = append(*b, byte(i))
		}
		for i := 123; i <= 126; i++ {
			*b = append(*b, byte(i))
		}
	}
}
