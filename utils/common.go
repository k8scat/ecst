package utils

import (
	"bytes"
	"fmt"
	"math/rand"
	"mime/multipart"
)

func GeneratePassword() string {
	source := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789()`~!@#$%^&*-_+=|{}[]:;'<>,.?"
	password := fmt.Sprintf("%s%s", RandomString(source, 16), "Vss^1")
	return password
}

func RandomString(source string, n int) string {
	rs := []rune(source)
	b := make([]byte, n)
	for i := range b {
		b[i] = source[rand.Intn(len(rs))]
	}
	return string(b)
}

func ParseFormPayload(data map[string]string) (payload *bytes.Buffer, contentType string, err error) {
	payload = &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	defer writer.Close()
	for k, v := range data {
		err = writer.WriteField(k, v)
		if err != nil {
			return
		}
	}
	contentType = writer.FormDataContentType()
	return
}
