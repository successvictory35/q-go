package qubit

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"time"
)

func MathRand(seed ...int64) float64 {
	if len(seed) > 0 {
		rand.Seed(seed[0])
		return rand.Float64()
	}

	rand.Seed(time.Now().UnixNano())
	return rand.Float64()
}

func CryptoRand(_ ...int64) float64 {
	n, err := crand.Int(crand.Reader, big.NewInt(1000))
	if err != nil {
		panic(err)
	}

	return float64(n.Int64()) / 1000.0
}
