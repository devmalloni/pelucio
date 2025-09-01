package pelucio

import (
	"math/big"
)

func MustFromString(s string, precision int) *big.Int {
	b, err := FromString(s, precision)
	if err != nil {
		panic(err)
	}

	return b
}

func FromString(s string, precision int) (*big.Int, error) {
	f, _, err := big.ParseFloat(s, 10, 256, big.ToNearestEven)
	if err != nil {
		return nil, err
	}

	scale := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precision)), nil))

	res := new(big.Int)
	f.Mul(f, scale).Int(res)

	return res, nil
}

func ToString(b *big.Int, precision int) string {
	rat := new(big.Rat).SetInt(b)

	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precision)), nil)
	rat.Quo(rat, new(big.Rat).SetInt(scale))

	return rat.FloatString(precision)
}
