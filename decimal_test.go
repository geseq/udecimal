package udecimal_test

import (
	"bytes"
	"encoding/json"
	"math"
	"testing"

	. "github.com/geseq/udecimal"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	f0 := MustParse("123.456")
	f1 := MustParse("123.456")

	assert.Equal(t, f0, f1)
	assert.True(t, f0.Equal(f1))
	assert.True(t, f1.Equal(f0))
	assert.Equal(t, uint64(123), f0.Int())
	assert.Equal(t, uint64(123), f1.Int())
	assert.Equal(t, "123.456", f0.String())
	assert.Equal(t, "123.456", f1.String())

	f0 = MustParseFloat(1)
	f1 = MustParseFloat(.5).Add(MustParseFloat(.5))
	f2 := MustParseFloat(.3).Add(MustParseFloat(.3)).Add(MustParseFloat(.4))

	assert.True(t, f0.Equal(f1))
	assert.True(t, f0.Equal(f2))

	f0 = MustParseFloat(.999)
	assert.Equal(t, "0.999", f0.String())
}

func TestEqual(t *testing.T) {
	f0 := Zero
	f1 := MustParse("123.456")
	assert.NotEqual(t, f0, f1)
	assert.False(t, f0.Equal(f1))
	assert.False(t, f1.Equal(f0))

	f1 = Zero
	assert.True(t, f0.Equal(f1))

	f0 = Zero
	assert.True(t, f0.Equal(f1))
	assert.Equal(t, uint64(0), f0.Int())
}

func TestNew(t *testing.T) {
	f := New(123, 1)
	assert.Equal(t, "1230", f.String())

	f = New(123, 0)
	assert.Equal(t, "123", f.String())

	//	f = New(123456789012, 9)
	//	assert.Equal(t, "123", f.String())

	f = New(123, -1)
	assert.Equal(t, "12.3", f.String())

	f = New(123456789001, -9)
	assert.Equal(t, "123.456789", f.String())

	f = New(123456789012, -9)
	assert.Equal(t, "123.4567890", f.StringN(7))

	f = New(123456789012, -9)
	assert.Equal(t, "123.45678901", f.StringN(8))
}

func TestParse(t *testing.T) {
	_, err := Parse("123")
	assert.NoError(t, err)

	_, err = Parse("123,456")
	assert.Error(t, err)

	_, err = Parse("abc")
	assert.Error(t, err)
}

func TestMustParse(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")

		}

	}()

	_ = MustParse("abc")
}

func TestNewI(t *testing.T) {
	f := NewI(123, 1)
	assert.Equal(t, "12.3", f.String())

	f = NewI(123, 0)
	assert.Equal(t, "123", f.String())

	f = NewI(123456789001, 9)
	assert.Equal(t, "123.456789", f.String())

	f = NewI(123456789012, 9)
	assert.Equal(t, "123.4567890", f.StringN(7))
	assert.Equal(t, "123.45678901", f.String())

	f = NewI(123456789012, 9)
	assert.Equal(t, "123.45678901", f.StringN(8))
}

func TestMaxValue(t *testing.T) {
	f0 := MustParse("12345678901")
	assert.Equal(t, f0.String(), "12345678901")
	assert.Panics(t, func() { f0 = MustParse("123456789012") })
	assert.Panics(t, func() { f0 = MustParse("-12345678901") })
	assert.Panics(t, func() { f0 = MustParse("-123456789012") })

	f0 = MustParse("99999999999")
	assert.Equal(t, f0.String(), "99999999999")

	f0 = MustParse("9.99999999")
	assert.Equal(t, f0.String(), "9.99999999")

	f0 = MustParse("99999999999.99999999")
	assert.Equal(t, f0.String(), "99999999999.99999999")

	f0 = MustParse("99999999999.12345678901234567890")
	assert.Equal(t, f0.String(), "99999999999.12345678")
}

func TestFloat(t *testing.T) {
	f0 := MustParse("123.456")
	f1 := MustParseFloat(123.456)

	if !f0.Equal(f1) {
		t.Error("should be equal", f0, f1)
	}

	f1 = MustParseFloat(0.0001)

	if f1.String() != "0.0001" {
		t.Error("should be equal", f1.String(), "0.0001")
	}

	f1 = MustParse(".1")
	f2 := MustParse(MustParseFloat(f1.Float()).String())
	if !f1.Equal(f2) {
		t.Error("should be equal", f1, f2)
	}
}

func TestInfinite(t *testing.T) {
	f0 := MustParse("0.10")
	f1 := MustParseFloat(0.10)

	if !f0.Equal(f1) {
		t.Error("should be equal", f0, f1)
	}

	f2 := MustParseFloat(0.0)
	for i := 0; i < 3; i++ {
		f2 = f2.Add(MustParseFloat(.10))
	}
	if f2.String() != "0.3" {
		t.Error("should be equal", f2.String(), "0.3")
	}

	f2 = MustParseFloat(0.0)
	for i := 0; i < 10; i++ {
		f2 = f2.Add(MustParseFloat(.10))
	}
	if f2.String() != "1" {
		t.Error("should be equal", f2.String(), "1")
	}

}

func TestAddSub(t *testing.T) {
	f0 := MustParse("1")
	f1 := MustParse("0.3333333")

	f2 := f0.Sub(f1)
	f2 = f2.Sub(f1)
	f2 = f2.Sub(f1)

	if f2.String() != "0.0000001" {
		t.Error("should be equal", f2.String(), "0.0000001")
	}
	f2 = f2.Sub(MustParse("0.0000001"))
	if f2.String() != "0" {
		t.Error("should be equal", f2.String(), "0")
	}

	f0 = MustParse("0")
	for i := 0; i < 10; i++ {
		f0 = f0.Add(MustParse("0.1"))
	}
	if f0.String() != "1" {
		t.Error("should be equal", f0.String(), "1")
	}

}

func TestMulDiv(t *testing.T) {
	f0 := MustParse("123.456")
	f1 := MustParse("1000")

	f2 := f0.Mul(f1)
	if f2.String() != "123456" {
		t.Error("should be equal", f2.String(), "123456")
	}
	f0 = MustParse("123456")
	f1 = MustParse("0.0001")

	f2 = f0.Mul(f1)
	if f2.String() != "12.3456" {
		t.Error("should be equal", f2.String(), "12.3456")
	}

	f0 = MustParse("10000.1")
	f1 = MustParse("10000")

	f2 = f0.Mul(f1)
	if f2.String() != "100001000" {
		t.Error("should be equal", f2.String(), "100001000")
	}

	f2 = f2.Div(f1)
	if !f2.Equal(f0) {
		t.Error("should be equal", f0, f2)
	}

	f0 = MustParse("2")
	f1 = MustParse("3")

	f2 = f0.Div(f1)
	if f2.String() != "0.66666667" {
		t.Error("should be equal", f2.String(), "0.66666667")
	}

	f0 = MustParse("1000")
	f1 = MustParse("10")

	f2 = f0.Div(f1)
	if f2.String() != "100" {
		t.Error("should be equal", f2.String(), "100")
	}

	f0 = MustParse("1000")
	f1 = MustParse("0.1")

	f2 = f0.Div(f1)
	if f2.String() != "10000" {
		t.Error("should be equal", f2.String(), "10000")
	}

	f0 = MustParse("1")
	f1 = MustParse("0.1")

	f2 = f0.Mul(f1)
	if f2.String() != "0.1" {
		t.Error("should be equal", f2.String(), "0.1")
	}

}

func TestNegatives(t *testing.T) {
	assert.Panics(t, func() { MustParse("-1") })

	f0 := MustParse("99")
	f1 := MustParse("100")

	assert.Panics(t, func() { f0.Sub(f1) })

	f0 = MustParse(".001")
	f1 = MustParse(".002")

	assert.Panics(t, func() { f0.Sub(f1) })
}

func TestOverflow(t *testing.T) {
	f0 := MustParseFloat(1.12345678)
	if f0.String() != "1.12345678" {
		t.Error("should be equal", f0.String(), "1.12345678")
	}
	f0 = MustParseFloat(1.123456789123)
	if f0.String() != "1.12345679" {
		t.Error("should be equal", f0.String(), "1.12345679")
	}
	f0 = MustParseFloat(1.0 / 3.0)
	if f0.String() != "0.33333333" {
		t.Error("should be equal", f0.String(), "0.33333333")
	}
	f0 = MustParseFloat(2.0 / 3.0)
	if f0.String() != "0.66666667" {
		t.Error("should be equal", f0.String(), "0.66666667")
	}
}

func TestNaN(t *testing.T) {
	assert.Panics(t, func() { MustParseFloat(math.NaN()) })
	assert.Panics(t, func() { MustParse("NaN") })

	f0 := MustParse("0.0004096")
	if f0.String() != "0.0004096" {
		t.Error("should be equal", f0.String(), "0.0004096")
	}
}

func TestIntFrac(t *testing.T) {
	f0 := MustParseFloat(1234.5678)
	if f0.Int() != 1234 {
		t.Error("should be equal", f0.Int(), 1234)
	}
	if f0.Frac() != .5678 {
		t.Error("should be equal", f0.Frac(), .5678)
	}
}

func TestString(t *testing.T) {
	f0 := MustParseFloat(1234.5678)
	if f0.String() != "1234.5678" {
		t.Error("should be equal", f0.String(), "1234.5678")
	}
	f0 = MustParseFloat(1234.0)
	if f0.String() != "1234" {
		t.Error("should be equal", f0.String(), "1234")
	}
}

func TestStringN(t *testing.T) {
	f0 := MustParse("1.1")
	s := f0.StringN(2)

	if s != "1.10" {
		t.Error("should be equal", s, "1.10")
	}
	f0 = MustParse("1")
	s = f0.StringN(2)

	if s != "1.00" {
		t.Error("should be equal", s, "1.00")
	}

	f0 = MustParse("1.123")
	s = f0.StringN(2)

	if s != "1.12" {
		t.Error("should be equal", s, "1.12")
	}
	f0 = MustParse("1.123")
	s = f0.StringN(2)

	if s != "1.12" {
		t.Error("should be equal", s, "1.12")
	}
	f0 = MustParse("1.123")
	s = f0.StringN(0)

	if s != "1" {
		t.Error("should be equal", s, "1")
	}
}

func TestRound(t *testing.T) {
	f0 := MustParse("1.12345")
	f1 := f0.Round(2)

	if f1.String() != "1.12" {
		t.Error("should be equal", f1, "1.12")
	}

	f1 = f0.Round(5)

	if f1.String() != "1.12345" {
		t.Error("should be equal", f1, "1.12345")
	}
	f1 = f0.Round(4)

	if f1.String() != "1.1235" {
		t.Error("should be equal", f1, "1.1235")
	}
}

func TestEncodeDecode(t *testing.T) {
	b := &bytes.Buffer{}

	f := MustParse("12345.12345")

	f.WriteTo(b)

	f0, err := ReadFrom(b)
	if err != nil {
		t.Error(err)
	}

	if !f.Equal(f0) {
		t.Error("don't match", f, f0)
	}

	data, err := f.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	f1 := MustParseFloat(0)
	f1.UnmarshalBinary(data)

	if !f.Equal(f1) {
		t.Error("don't match", f, f0)
	}
}

type JStruct struct {
	F Decimal `json:"f"`
}

func TestJSON(t *testing.T) {
	j := JStruct{}

	f := MustParse("1234567.1234567")
	j.F = f

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	err := enc.Encode(&j)
	if err != nil {
		t.Error(err)
	}

	j.F = Zero

	dec := json.NewDecoder(&buf)

	err = dec.Decode(&j)
	if err != nil {
		t.Error(err)
	}

	if !j.F.Equal(f) {
		t.Error("don't match", j.F, f)
	}
}
