package udecimal

// release under the terms of file LICENSE and LICENSE-FIXED

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// Decimal is a decimal precision 38.24 number (supports 11.7 digits).
type Decimal struct {
	fp uint64
}

// the following constants can be changed to configure a different number of decimal places - these are
// the only required changes.
const nPlaces = 8
const scale = uint64(10 * 10 * 10 * 10 * 10 * 10 * 10 * 10)
const zeros = "00000000"
const MAX = float64(99999999999.99999999)

var Zero = Decimal{fp: 0}

var errTooLarge = errors.New("significand too large")
var errFormat = errors.New("invalid encoding")

func Parse(s string) (Decimal, error) {
	if strings.ContainsAny(s, "eE") {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return Zero, err
		}
		return ParseFloat(f)
	}
	period := strings.Index(s, ".")
	var i uint64
	var f uint64
	var err error
	if period == -1 {
		i, err = strconv.ParseUint(s, 10, 64)
		if err != nil {
			return Zero, errors.New("cannot parse")
		}
	} else {
		if len(s[:period]) > 0 {
			i, err = strconv.ParseUint(s[:period], 10, 64)
			if err != nil {
				return Zero, errors.New("cannot parse")
			}
		}
		fs := s[period+1:]
		fs = fs + zeros[:max(0, nPlaces-len(fs))]
		f, err = strconv.ParseUint(fs[0:nPlaces], 10, 64)
		if err != nil {
			return Zero, errors.New("cannot parse")
		}
	}
	if float64(i) > MAX {
		return Zero, errTooLarge
	}
	return Decimal{fp: (i*scale + f)}, nil
}

// MustParse creates a new Fixed from a string, and panics if the string could not be parsed
func MustParse(s string) Decimal {
	f, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return f
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ParseFloat creates a Decimal from an float64, rounding at the 8th decimal place
func ParseFloat(f float64) (Decimal, error) {
	if math.IsNaN(f) {
		return Zero, errors.New("invalid input")
	}
	if f >= MAX || f < 0 {
		return Zero, errors.New("invalid input")
	}
	round := .5
	if f < 0 {
		round = -0.5
	}

	return Decimal{fp: uint64(f*float64(scale) + round)}, nil
}

// MustParseFloat creates a new Fixed from a string, and panics if the string could not be parsed
func MustParseFloat(f float64) Decimal {
	r, err := ParseFloat(f)
	if err != nil {
		panic(err)
	}
	return r
}

// New returns a new fixed-point decimal, value * 10 ^ exp.
func New(value uint64, exp int32) Decimal {
	if exp >= 0 {
		mul := uint64(math.Pow10(int(exp)))
		return NewI(value, 0).Mul(NewI(mul, 0))
	}

	return NewI(value, uint(exp*-1))
}

// NewI creates a Decimal for an integer, moving the decimal point n places to the left
// For example, NewI(123,1) becomes 12.3. If n > 7, the value is truncated
func NewI(i uint64, n uint) Decimal {
	if n > nPlaces {
		i = i / uint64(math.Pow10(int(n-nPlaces)))
		n = nPlaces
	}

	i = i * uint64(math.Pow10(int(nPlaces-n)))

	return Decimal{fp: i}
}

func (f Decimal) IsZero() bool {
	return f.Equal(Zero)
}

// Float converts the Decimal to a float64
func (f Decimal) Float() float64 {
	return float64(f.fp) / float64(scale)
}

// Add adds f0 to f producing a Decimal.
func (f Decimal) Add(f0 Decimal) Decimal {
	if f0.fp > math.MaxUint64-f.fp {
		panic("decimal overflow")
	}

	return Decimal{fp: f.fp + f0.fp}
}

// Sub subtracts f0 from f producing a Decimal.
func (f Decimal) Sub(f0 Decimal) Decimal {
	if f.fp < f0.fp {
		panic("decimal overflow")
	}

	return Decimal{fp: f.fp - f0.fp}
}

// Mul multiplies f by f0 returning a Decimal.
func (f Decimal) Mul(f0 Decimal) Decimal {
	fp_a := f.fp / scale
	fp_b := f.fp % scale

	fp0_a := f0.fp / scale
	fp0_b := f0.fp % scale

	var result uint64

	if fp0_a != 0 {
		result = fp_a * fp0_a
		if float64(result) > MAX {
			panic("decimal overflow")
		}

		result = result*scale + fp_b*fp0_a
	}

	if fp0_b != 0 {
		result = result + (fp_a * fp0_b) + ((fp_b)*fp0_b)/scale
	}

	return Decimal{fp: result}
}

// Div divides f by f0 returning a Decimal.
func (f Decimal) Div(f0 Decimal) Decimal {
	res, err := ParseFloat(f.Float() / f0.Float())
	if err != nil {
		panic("decimal overflow")
	}
	return res
}

// Round returns a rounded (half-up, away from zero) to n decimal places
func (f Decimal) Round(n int) Decimal {
	round := .5
	if f.fp < 0 {
		round = -0.5
	}

	f0 := f.Frac()
	f0 = f0*math.Pow10(n) + round
	f0 = float64(int(f0)) / math.Pow10(n)

	res, err := ParseFloat(float64(f.Int()) + f0)
	if err != nil {
		panic("decimal overflow")
	}
	return res
}

// Equal returns true if the f == f0.
func (f Decimal) Equal(f0 Decimal) bool {
	if f.fp == f0.fp {
		return true
	}
	return false
}

// GreaterThan returns true if the f > f0.
func (f Decimal) GreaterThan(f0 Decimal) bool {
	if f.fp > f0.fp {
		return true
	}
	return false
}

// GreaterThaOrEqual returns true if the f >= f0.
func (f Decimal) GreaterThanOrEqual(f0 Decimal) bool {
	if f.fp >= f0.fp {
		return true
	}
	return false
}

// LessThan returns true if the f < f0.
func (f Decimal) LessThan(f0 Decimal) bool {
	if f.fp < f0.fp {
		return true
	}
	return false
}

// LessThan returns true if the f <= f0.
func (f Decimal) LessThanOrEqual(f0 Decimal) bool {
	if f.fp <= f0.fp {
		return true
	}
	return false
}

// Cmp compares two Decimal. If f == f0, return 0. If f > f0, return 1. If f < f0, return -1.
func (f Decimal) Cmp(f0 Decimal) int {
	if f.fp == f0.fp {
		return 0
	}
	if f.fp < f0.fp {
		return -1
	}
	return 1
}

// String converts a Decimal to a string, dropping trailing zeros
func (f Decimal) String() string {
	s, point := f.tostr()
	if point == -1 {
		return s
	}
	index := len(s) - 1
	for ; index != point; index-- {
		if s[index] != '0' {
			return s[:index+1]
		}
	}
	return s[:point]
}

// StringN converts a Decimal to a String with a specified number of decimal places, truncating as required
func (f Decimal) StringN(decimals int) string {
	s, point := f.tostr()

	if point == -1 {
		return s
	}
	if decimals == 0 {
		return s[:point]
	} else {
		return s[:point+decimals+1]
	}
}

func (f Decimal) tostr() (string, int) {
	fp := f.fp
	if fp == 0 {
		return "0." + zeros, 1
	}

	b := make([]byte, 24)
	b = itoa(b, fp)

	return string(b), len(b) - nPlaces - 1
}

func itoa(buf []byte, val uint64) []byte {
	i := len(buf) - 1
	idec := i - nPlaces
	for val >= 10 || i >= idec {
		buf[i] = byte(val%10 + '0')
		i--
		if i == idec {
			buf[i] = '.'
			i--
		}
		val /= 10
	}
	buf[i] = byte(val + '0')
	return buf[i:]
}

// Int return the integer portion of the Decimal
func (f Decimal) Int() uint64 {
	return f.fp / scale
}

// Frac return the fractional portion of the Decimal
func (f Decimal) Frac() float64 {
	return float64(f.fp%scale) / float64(scale)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (f *Decimal) UnmarshalBinary(data []byte) error {
	fp, n := binary.Uvarint(data)
	if n < 0 {
		return errFormat
	}
	f.fp = fp
	return nil
}

// UnmarshalBinaryData Unmarshals data and returns n
func (f *Decimal) UnmarshalBinaryData(data []byte) (rem []byte, err error) {
	fp, n := binary.Uvarint(data)
	if n < 0 {
		return data, errFormat
	}
	f.fp = fp
	return data[n:], err
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (f Decimal) MarshalBinary() (data []byte, err error) {
	var buffer [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(buffer[:], f.fp)
	return buffer[:n], nil
}

// WriteTo write the Decimal to an io.Writer, returning the number of bytes written
func (f Decimal) WriteTo(w io.ByteWriter) error {
	return writeUvarint(w, f.fp)
}

// ReadFrom reads a Decimal from an io.Reader
func ReadFrom(r io.ByteReader) (Decimal, error) {
	fp, err := binary.ReadUvarint(r)
	if err != nil {
		return Zero, err
	}
	return Decimal{fp: fp}, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (f *Decimal) UnmarshalJSON(bytes []byte) error {
	s := string(bytes)
	if s == "null" {
		return nil
	}

	decimal, err := Parse(s)
	*f = decimal
	if err != nil {
		return fmt.Errorf("Error decoding string '%s': %s", s, err)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (f Decimal) MarshalJSON() ([]byte, error) {
	buffer := make([]byte, 24)
	return itoa(buffer, f.fp), nil
}
