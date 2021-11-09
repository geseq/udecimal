package udecimal

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
)

func BenchmarkAddDecimal(b *testing.B) {
	f0 := NewF(1)
	f1 := NewF(1)

	for i := 0; i < b.N; i++ {
		f1 = f1.Add(f0)
	}
}
func BenchmarkAddShopspringDecimal(b *testing.B) {
	f0 := decimal.NewFromFloat(1)
	f1 := decimal.NewFromFloat(1)

	for i := 0; i < b.N; i++ {
		f1 = f1.Add(f0)
	}
}
func BenchmarkAddBigInt(b *testing.B) {
	f0 := big.NewInt(1)
	f1 := big.NewInt(1)

	for i := 0; i < b.N; i++ {
		f1 = f1.Add(f1, f0)
	}
}
func BenchmarkAddBigFloat(b *testing.B) {
	f0 := big.NewFloat(1)
	f1 := big.NewFloat(1)

	for i := 0; i < b.N; i++ {
		f1 = f1.Add(f1, f0)
	}
}

func BenchmarkMulDecimal(b *testing.B) {
	f0 := NewF(123456789.0)
	f1 := NewF(1234.0)

	for i := 0; i < b.N; i++ {
		f0.Mul(f1)
	}
}
func BenchmarkMulShopspringDecimal(b *testing.B) {
	f0 := decimal.NewFromFloat(123456789.0)
	f1 := decimal.NewFromFloat(1234.0)

	for i := 0; i < b.N; i++ {
		f0.Mul(f1)
	}
}
func BenchmarkMulBigInt(b *testing.B) {
	f0 := big.NewInt(123456789)
	f1 := big.NewInt(1234)

	var x big.Int
	for i := 0; i < b.N; i++ {
		x.Mul(f0, f1)
	}
}
func BenchmarkMulBigFloat(b *testing.B) {
	f0 := big.NewFloat(123456789.0)
	f1 := big.NewFloat(1234.0)

	var x big.Float
	for i := 0; i < b.N; i++ {
		x.Mul(f0, f1)
	}
}

func BenchmarkDivDecimal(b *testing.B) {
	f0 := NewF(123456789.0)
	f1 := NewF(1234.0)

	for i := 0; i < b.N; i++ {
		f0.Div(f1)
	}
}
func BenchmarkDivShopspringDecimal(b *testing.B) {
	f0 := decimal.NewFromFloat(123456789.0)
	f1 := decimal.NewFromFloat(1234.0)

	for i := 0; i < b.N; i++ {
		f0.Div(f1)
	}
}
func BenchmarkDivBigInt(b *testing.B) {
	f0 := big.NewInt(123456789)
	f1 := big.NewInt(1234)

	var x big.Int
	for i := 0; i < b.N; i++ {
		x.Div(f0, f1)
	}
}
func BenchmarkDivBigFloat(b *testing.B) {
	f0 := big.NewFloat(123456789.0)
	f1 := big.NewFloat(1234.0)

	var x big.Float
	for i := 0; i < b.N; i++ {
		x.Quo(f0, f1)
	}
}

func BenchmarkCmpDecimal(b *testing.B) {
	f0 := NewF(123456789.0)
	f1 := NewF(1234.0)

	for i := 0; i < b.N; i++ {
		f0.Cmp(f1)
	}
}
func BenchmarkCmpShopspringDecimal(b *testing.B) {
	f0 := decimal.NewFromFloat(123456789.0)
	f1 := decimal.NewFromFloat(1234.0)

	for i := 0; i < b.N; i++ {
		f0.Cmp(f1)
	}
}
func BenchmarkCmpBigInt(b *testing.B) {
	f0 := big.NewInt(123456789)
	f1 := big.NewInt(1234)

	for i := 0; i < b.N; i++ {
		f0.Cmp(f1)
	}
}
func BenchmarkCmpBigFloat(b *testing.B) {
	f0 := big.NewFloat(123456789.0)
	f1 := big.NewFloat(1234.0)

	for i := 0; i < b.N; i++ {
		f0.Cmp(f1)
	}
}

func BenchmarkStringDecimal(b *testing.B) {
	f0 := NewF(123456789.12345)

	for i := 0; i < b.N; i++ {
		f0.String()
	}
}
func BenchmarkStringNDecimal(b *testing.B) {
	f0 := NewF(123456789.12345)

	for i := 0; i < b.N; i++ {
		f0.StringN(5)
	}
}
func BenchmarkStringShopspringDecimal(b *testing.B) {
	f0 := decimal.NewFromFloat(123456789.12345)

	for i := 0; i < b.N; i++ {
		f0.String()
	}
}
func BenchmarkStringBigInt(b *testing.B) {
	f0 := big.NewInt(123456789)

	for i := 0; i < b.N; i++ {
		f0.String()
	}
}
func BenchmarkStringBigFloat(b *testing.B) {
	f0 := big.NewFloat(123456789.12345)

	for i := 0; i < b.N; i++ {
		f0.String()
	}
}

func BenchmarkWriteTo(b *testing.B) {
	f0 := NewF(123456789.0)

	buf := new(bytes.Buffer)

	for i := 0; i < b.N; i++ {
		f0.WriteTo(buf)
	}
}

var res bool

func BenchmarkEqualDecimal(b *testing.B) {
	f0 := NewF(1)
	f1 := NewF(1)

	for i := 0; i < b.N; i++ {
		res = f1.Equal(f0)
	}
}

func BenchmarkLessThanDecimal(b *testing.B) {
	f0 := NewF(1)
	f1 := NewF(1)

	for i := 0; i < b.N; i++ {
		res = f1.LessThan(f0)
	}
}

func BenchmarkLessThanOrEqualDecimal(b *testing.B) {
	f0 := NewF(1)
	f1 := NewF(1)

	for i := 0; i < b.N; i++ {
		res = f1.LessThanOrEqual(f0)
	}
}

func BenchmarkGreaterThanDecimal(b *testing.B) {
	f0 := NewF(1)
	f1 := NewF(1)

	for i := 0; i < b.N; i++ {
		res = f1.GreaterThan(f0)
	}
}

func BenchmarkGreaterThanOrEqualDecimal(b *testing.B) {
	f0 := NewF(1)
	f1 := NewF(1)

	for i := 0; i < b.N; i++ {
		res = f1.GreaterThanOrEqual(f0)
	}
}
