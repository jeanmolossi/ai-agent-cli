package str

import "crypto/rand"

var letters = "123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Random(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	for i, v := range b {
		b[i] = letters[v%byte(len(letters))]
	}

	return string(b)
}
