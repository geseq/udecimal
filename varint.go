package decimal

import "io"

// The binary encoding package does not offer 'write' methods so the allocation costs in WriteTo are high
// so we duplicate the code here and implement them

// WriteUvarint encodes a uint64 onto w
func writeUvarint(w io.ByteWriter, x uint64) error {
	i := 0
	for x >= 0x80 {
		err := w.WriteByte(byte(x) | 0x80)
		if err != nil {
			return err
		}
		x >>= 7
		i++
	}
	return w.WriteByte(byte(x))
}
