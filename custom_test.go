package struc

import (
	"bytes"
	"encoding/binary"
	"io"
	"strconv"
	"testing"
)

type Int3 uint32

func (i *Int3) Pack(p []byte, opt *Options) (int, error) {
	var tmp [4]byte
	binary.BigEndian.PutUint32(tmp[:], uint32(*i))
	copy(p, tmp[1:])
	return 3, nil
}
func (i *Int3) Unpack(r io.Reader, length int, opt *Options) error {
	var tmp [4]byte
	if _, err := r.Read(tmp[1:]); err != nil {
		return err
	}
	*i = Int3(binary.BigEndian.Uint32(tmp[:]))
	return nil
}
func (i *Int3) Size(opt *Options) int {
	return 3
}
func (i *Int3) String() string {
	return strconv.FormatUint(uint64(*i), 10)
}

func TestCustom(t *testing.T) {
	var buf bytes.Buffer
	var i Int3 = 3
	if err := Pack(&buf, &i); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), []byte{0, 0, 3}) {
		t.Fatal("error packing custom int")
	}
	var i2 Int3
	if err := Unpack(&buf, &i2); err != nil {
		t.Fatal(err)
	}
	if i2 != 3 {
		t.Fatal("error unpacking custom int")
	}
}

type Int3Struct struct {
	I Int3
}

func TestCustomStruct(t *testing.T) {
	var buf bytes.Buffer
	i := Int3Struct{3}
	if err := Pack(&buf, &i); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), []byte{0, 0, 3}) {
		t.Fatal("error packing custom int struct")
	}
	var i2 Int3Struct
	if err := Unpack(&buf, &i2); err != nil {
		t.Fatal(err)
	}
	if i2.I != 3 {
		t.Fatal("error unpacking custom int struct")
	}
}

// TODO: slices of custom types don't work yet
type ArrayInt3Struct struct {
	I [2]Int3
}

func TestArrayOfCustomStructPanic(t *testing.T) {
	// Expect panic.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("did not panic")
		}
	}()
	var buf bytes.Buffer
	i := ArrayInt3Struct{[2]Int3{3, 4}}
	if err := Pack(&buf, &i); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), []byte{0, 0, 3}) {
		t.Fatal("error packing custom int struct")
	}
	var i2 ArrayInt3Struct
	if err := Unpack(&buf, &i2); err != nil {
		t.Fatal(err)
	}
	if i2.I[0] != 3 && i2.I[1] != 4 {
		t.Fatal("error unpacking custom int struct")
	}
}

// Slice of uint8, stored in a zero terminated list.
type IntSlice []uint8

func (ia *IntSlice) Pack(p []byte, opt *Options) (int, error) {
	for i, value := range *ia {
		p[i] = value
	}

	return len(*ia) + 1, nil
}

func (ia *IntSlice) Unpack(r io.Reader, length int, opt *Options) error {
	for {
		var value uint8
		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}
		*ia = append(*ia, value)
		if value == 0 {
			break
		}
	}
	return nil
}

func (ia *IntSlice) Size(opt *Options) int {
	return len(*ia) + 1
}

func (ia *IntSlice) String() string {
	panic("not implemented")
}

func TestCustomLength(t *testing.T) {
	var buf bytes.Buffer
	i := make(IntSlice, 0)
	i = append(i, 128)
	i = append(i, 64)
	i = append(i, 32)

	if err := Pack(&buf, &i); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), []byte{128, 64, 32, 0}) {
		t.Fatal("error packing custom int array")
	}
	var i2 IntSlice
	if err := Unpack(&buf, &i2); err != nil {
		t.Fatal(err)
	}
	if i2[0] != 128 {
		t.Fatal("error unpacking custom int array")
	}
}

type IntSliceStruct struct {
	I IntSlice
	N uint8 // A field after to ensure the length is correct.
}

func TestCustomLengthStruct(t *testing.T) {
	var buf bytes.Buffer
	i := IntSliceStruct{
		I: IntSlice{128, 64, 32},
		N: 192,
	}
	if err := Pack(&buf, &i); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), []byte{128, 64, 32, 0, 192}) {
		t.Fatal("error packing custom int array struct")
	}
	var i2 IntSliceStruct
	if err := Unpack(&buf, &i2); err != nil {
		t.Fatal(err)
	}
	if i2.I[0] != 128 || i2.N != 192 {
		t.Fatal("error unpacking custom int array struct")
	}
}
