package password

import (
	"crypto/rand"
	"math/big"
)

const (
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberBytes  = "0123456789"
	symbolBytes  = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	defaultLength = 16
)

type Options struct {
	Length       int
	IncludeOr    bool // if true, ensure at least one of chosen types
	IncludeNums  bool
	IncludeSyms  bool
}

func Generate(opts Options) (string, error) {
	if opts.Length <= 0 {
		opts.Length = defaultLength
	}

	charSet := letterBytes
	if opts.IncludeNums {
		charSet += numberBytes
	}
	if opts.IncludeSyms {
		charSet += symbolBytes
	}

	b := make([]byte, opts.Length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			return "", err
		}
		b[i] = charSet[num.Int64()]
	}

	return string(b), nil
}
