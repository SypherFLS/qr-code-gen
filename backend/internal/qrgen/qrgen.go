package qrgen

import (
	"github.com/skip2/go-qrcode"
)

func Gen(url string) ([]byte, error) {
	png, err := qrcode.Encode(
		url,
		qrcode.Medium,
		256,
	)

	if err != nil {
		return nil, err
	}

	return png, nil
}
