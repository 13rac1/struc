package struc

import (
	"bytes"
	"testing"
)

type nullTermBytesStruct struct {
	S NullTermBytes
	N uint8 // A field after to ensure the length is correct.
}

func TestNullTermStringStruct(t *testing.T) {
	var buf bytes.Buffer
	number8 := uint8(128)
	helloWorld := "Hello World!"
	byteSlice := []byte(helloWorld)
	i := nullTermBytesStruct{
		S: byteSlice,
		N: number8,
	}
	if err := Pack(&buf, &i); err != nil {
		t.Fatal(err)
	}

	nullTermByteSlice := append(byteSlice, 0)
	structBytes := append(nullTermByteSlice, number8)
	if !bytes.Equal(buf.Bytes(), structBytes) {
		t.Fatal("error packing custom null terminated string struct")
	}

	var i2 nullTermBytesStruct
	if err := Unpack(&buf, &i2); err != nil {
		t.Fatal(err)
	}

	if i2.S.String() != helloWorld || i2.N != number8 {
		t.Fatal("error unpacking custom null terminated string struct")
	}
}
