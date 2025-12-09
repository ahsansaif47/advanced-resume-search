package parser

import (
	"github.com/otiai10/gosseract/v2"
)

func InitClient() *gosseract.Client {
	client := gosseract.NewClient()
	return client
}
func GetText(client *gosseract.Client, path string) (string, error) {
	client.SetImage(path)
	text, err := client.Text()
	return text, err
}
