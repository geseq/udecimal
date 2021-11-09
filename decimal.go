package udecimal

// release under the terms of file license.txt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// Decimal is a decimal precision 38.24 number (supports 11.7 digits). It supports NaN.
type Decimal struct {
	fp uint64
}

// the following constants can be changed to configure a different number of decimal places - these are
// the only required changes. only 18 significant digits are supported due to NaN

const nPlaces = 8
const scale = uint64(10 * 10 * 10 * 10 * 10 * 10 * 10 * 10)
const zeros = "00000000"
const MAX = float64(99999999999.99999999)

const nan = uint64(1<<63 - 1)

var NaN = Decimal{fp: nan}
var Zero = Decimal{fp: 0}

var errTooLarge = errors.New("significand too large")
var errFormat = errors.New("invalid encoding")

// NewS creates a new Decimal from a string, returning NaN if the string could not be parsed
func NewS(s string) Decimal {
	f, _ := NewSErr(s)
	return f
}

// NewSErr creates a new Decimal from a string, returning NaN, and error if the string could not be parsed
func NewSErr(s string) (Decimal, error) {
	if strings.ContainsAny(s, "eE") {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return NaN, err
		}
		return NewF(f), nil
	}
	if "NaN" == s {
		return NaN, nil
	}
	period := strings.Index(s, ".")
	var i uint64
	var f uint64
	var err error
	if period == -1 {
		i, err = strconv.ParseUint(s, 10, 64)
		if err != nil {
			return NaN, errors.New("cannot parse")
		}
	} else {
		if len(s[:period]) > 0 {
			i, err = strconv.ParseUint(s[:period], 10, 64)
			if err != nil {
				return NaN, errors.New("cannot parse")
			}
		}
		fs := s[period+1:]
		fs = fs + zeros[:max(0, nPlaces-len(fs))]
		f, err = strconv.ParseUint(fs[0:nPlaces], 10, 64)
		if err != nil {
			return NaN, errors.New("cannot parse")
		}
	}
	if float64(i) > MAX {
		return NaN, errTooLarge
	}
	return Decimal{fp: (i*scale + f)}, nil
}

// Parse creates a new Fixed from a string, returning NaN, and error if the string could not be parsed. Same as NewSErr
// but more standard naming
func Parse(s string) (Decimal, error) {
	return NewSErr(s)

}

// MustParse creates a new Fixed from a string, and panics if the string could not be parsed
func MustParse(s string) Decimal {
	f, err := NewSErr(s)
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

// NewF creates a Decimal from an float64, rounding at the 8th decimal place
func NewF(f float64) Decimal {
	if math.IsNaN(f) {
		return Decimal{fp: nan}
	}
	if f >= MAX || f <= -MAX {
		return NaN
	}
	round := .5
	if f < 0 {
		round = -0.5
	}

	return Decimal{fp: uint64(f*float64(scale) + round)}
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

func (f Decimal) IsNaN() bool {
	return f.fp == nan
}

func (f Decimal) IsZero() bool {
	return f.Equal(Zero)
}

// Float converts the Decimal to a float64
func (f Decimal) Float() float64 {
	if f.IsNaN() {
		return math.NaN()
	}
	return float64(f.fp) / float64(scale)
}

// Add adds f0 to f producing a Decimal. If either operand is NaN, NaN is returned
func (f Decimal) Add(f0 Decimal) Decimal {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}
	return Decimal{fp: f.fp + f0.fp}
}

// Sub subtracts f0 from f producing a Decimal. If either operand is NaN, NaN is returned
func (f Decimal) Sub(f0 Decimal) Decimal {
	if f.IsNaN() || f0.IsNaN() || f.fp < f0.fp {
		return NaN
	}

	return Decimal{fp: f.fp - f0.fp}
}

// Mul multiplies f by f0 returning a Decimal. If either operand is NaN, NaN is returned
func (f Decimal) Mul(f0 Decimal) Decimal {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}

	fp_a := f.fp / scale
	fp_b := f.fp % scale

	fp0_a := f0.fp / scale
	fp0_b := f0.fp % scale

	var result uint64

	if fp0_a != 0 {
		result = fp_a * fp0_a
		if float64(result) > MAX {
			return NaN
		}

		result = result*scale + fp_b*fp0_a
	}

	if fp0_b != 0 {
		result = result + (fp_a * fp0_b) + ((fp_b)*fp0_b)/scale
	}

	return Decimal{fp: result}
}

// Div divides f by f0 returning a Decimal. If either operand is NaN, NaN is returned
func (f Decimal) Div(f0 Decimal) Decimal {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}
	return NewF(f.Float() / f0.Float())
}

// Round returns a rounded (half-up, away from zero) to n decimal places
func (f Decimal) Round(n int) Decimal {
	if f.IsNaN() {
		return NaN
	}

	round := .5
	if f.fp < 0 {
		round = -0.5
	}

	f0 := f.Frac()
	f0 = f0*math.Pow10(n) + round
	f0 = float64(int(f0)) / math.Pow10(n)

	return NewF(float64(f.Int()) + f0)
}

// Equal returns true if the f == f0. If either operand is NaN, false is returned. Use IsNaN() to test for NaN
func (f Decimal) Equal(f0 Decimal) bool {
	if f.IsNaN() || f0.IsNaN() {
		return false
	}

	if f.fp == f0.fp {
		return true
	}
	return false
}

// GreaterThan tests Cmp() for 1
func (f Decimal) GreaterThan(f0 Decimal) bool {
	if f.IsNaN() || f0.IsNaN() {
		return false
	}

	if f.fp > f0.fp {
		return true
	}
	return false
}

// GreaterThaOrEqual tests Cmp() for 1 or 0
func (f Decimal) GreaterThanOrEqual(f0 Decimal) bool {
	if f.IsNaN() || f0.IsNaN() {
		return false
	}

	if f.fp >= f0.fp {
		return true
	}
	return false
}

// LessThan tests Cmp() for -1
func (f Decimal) LessThan(f0 Decimal) bool {
	if f.IsNaN() || f0.IsNaN() {
		return false
	}

	if f.fp < f0.fp {
		return true
	}
	return false
}

// LessThan tests Cmp() for -1 or 0
func (f Decimal) LessThanOrEqual(f0 Decimal) bool {
	if f.IsNaN() || f0.IsNaN() {
		return false
	}

	if f.fp <= f0.fp {
		return true
	}
	return false
}

// Cmp compares two Decimal. If f == f0, return 0. If f > f0, return 1. If f < f0, return -1. If both are NaN, return 0. If f is NaN, return 1. If f0 is NaN, return -1
func (f Decimal) Cmp(f0 Decimal) int {
	if f.IsNaN() && f0.IsNaN() {
		return 0
	}
	if f.IsNaN() {
		return 1
	}
	if f0.IsNaN() {
		return -1
	}

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
	if fp == nan {
		return "NaN", -1
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

// Int return the integer portion of the Decimal, or 0 if NaN
func (f Decimal) Int() uint64 {
	if f.IsNaN() {
		return 0
	}
	return f.fp / scale
}

// Frac return the fractional portion of the Decimal, or NaN if NaN
func (f Decimal) Frac() float64 {
	if f.IsNaN() {
		return math.NaN()
	}
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
		return NaN, err
	}
	return Decimal{fp: fp}, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (f *Decimal) UnmarshalJSON(bytes []byte) error {
	s := string(bytes)
	if s == "null" {
		return nil
	}
	if s == "\"NaN\"" {
		*f = NaN
		return nil
	}

	decimal, err := NewSErr(s)
	*f = decimal
	if err != nil {
		return fmt.Errorf("Error decoding string '%s': %s", s, err)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (f Decimal) MarshalJSON() ([]byte, error) {
	if f.IsNaN() {
		return []byte("\"NaN\""), nil
	}

	buffer := make([]byte, 24)
	return itoa(buffer, f.fp), nil
}
