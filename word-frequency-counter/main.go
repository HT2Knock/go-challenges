package main

import (
	"fmt"
	"strings"
	"unicode"
)

// string manipulation with char
func CountWordFrequency(text string) map[string]int {
	freq := make(map[string]int)

	var word strings.Builder
	for i := 0; i < len(text); i++ {
		b := text[i]

		if 'A' <= b && b <= 'Z' {
			word.WriteByte(b | ' ')
		} else if ('0' <= b && b <= '9') || ('a' <= b && b <= 'z') {
			word.WriteByte(b)
		} else if b == '\'' {
			continue
		} else {
			if word.Len() > 0 {
				freq[word.String()]++
				word.Reset()
			}
		}
	}

	if word.Len() > 0 {
		freq[word.String()]++
	}

	return freq
}

func IdiomaticCountFreq(text string) map[string]int {
	freq := make(map[string]int)

	f := func(c rune) bool {
		return !unicode.IsDigit(c) && !unicode.IsLetter(c)
	}

	words := strings.FieldsFunc(text, f)

	for _, word := range words {
		freq[strings.ToLower(word)]++
	}

	return freq
}

func main() {
	fmt.Println(CountWordFrequency("Hello, hello! How are you doing today? Today is a great day."))
}
