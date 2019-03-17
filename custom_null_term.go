package struc

import (
	"encoding/binary"
	"io"
)

// Implements lunixbochs/struct/Custom field type

type NullTermBytes []byte

func (ntb *NullTermBytes) Pack(p []byte, opt *Options) (int, error) {
	for i, value := range *ntb {
		p[i] = byte(value)
	}

	return len(*ntb) + 1, nil
}

func (ntb *NullTermBytes) Unpack(r io.Reader, length int, opt *Options) error {
	for {
		var value byte
		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}
		if value == 0 {
			// The byte array does not include the NULL terminator.
			break
		}
		*ntb = append(*ntb, value)
	}
	return nil
}

func (ntb *NullTermBytes) Size(opt *Options) int {
	return len(*ntb) + 1
}

func (ntb *NullTermBytes) String() string {
	return string(*ntb)
}
