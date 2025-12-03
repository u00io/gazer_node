package utils

import (
	"errors"
	"math/big"
	"strings"
)

const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func Base58ToBytes(str string) ([]byte, error) {
	if len(str) == 0 {
		return []byte{}, nil
	}

	num := big.NewInt(0)
	base := big.NewInt(58)

	for _, c := range str {
		index := int64(strings.IndexRune(base58Alphabet, c))
		if index == -1 {
			return nil, errors.New("invalid Base58 character")
		}
		num.Mul(num, base)
		num.Add(num, big.NewInt(index))
	}

	// Получаем байты из BigInt
	bytes := num.Bytes()

	// Добавляем ведущие нули для каждого символа '1'
	leadingZeros := 0
	for _, c := range str {
		if c == '1' {
			leadingZeros++
		} else {
			break
		}
	}

	if leadingZeros > 0 {
		result := make([]byte, leadingZeros+len(bytes))
		copy(result[leadingZeros:], bytes)
		return result, nil
	}

	return bytes, nil
}

func BytesToBase58(bytes []byte) string {
	if len(bytes) == 0 {
		return ""
	}

	// Конвертируем в BigInt для работы с большими числами
	var num big.Int
	num.SetBytes(bytes)

	var result string
	base := big.NewInt(58)
	zero := big.NewInt(0)

	for num.Cmp(zero) > 0 {
		remainder := new(big.Int)
		num.DivMod(&num, base, remainder)
		result = string(base58Alphabet[remainder.Int64()]) + result
	}

	// Добавляем ведущие '1' для каждого нулевого байта
	for i := 0; i < len(bytes) && bytes[i] == 0; i++ {
		result = "1" + result
	}

	return result
}
