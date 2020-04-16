package tools

import (
	"math/rand"
)

// codeRunes is set of runes that can be used to generate
// readable codes, on the standard Latin alphabet and digits.
// Characters that easy to misread are excluded:
//
// - E, 3
// - I,  1
// - D, O, Q, 0
// - I, J
// - U, V
// - N, M (phoneticaly)
var codeRunes = []rune("ABCFGHKLPRSTWXYZ2456789")

// GenerateCode generates an easily readable random code.
//
// The code consistens of three groups of two characters, This
// grouping is chosen to prevent a group from forming a bad word.
// This gives a total 308 million options, or an entropy of 28 bits.
func GenerateCode() string {
	code := make([]rune, 7)
	for i := range code {
		if i%4 == 3 {
			code[i] = '-'
		} else {
			code[i] = codeRunes[rand.Intn(len(codeRunes))]
		}
	}
	return string(code)
}
