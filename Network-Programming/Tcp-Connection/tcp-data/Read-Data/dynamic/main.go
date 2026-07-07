package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"testing"
)

const(
	BinaryType uint8 = iota + 1
	StringType
	MaxPayloadSize uint32 = 10 << 20 //10 mb payload for security purposes 
)
var ErrMaxPayloadSize = errors.New("maximum payload size exceeded")
type Payload interface{//describes methods each type must implement
	fmt.Stringer
	io.ReaderFrom
	io.WriterTo
	Bytes() []byte
}
type Binary []byte//byte slice
func(m Binary) Bytes() []byte{return m}//casts itself
func(m Binary) String() string{return string(m)}//casts itself as String

func (m Binary) WriteTo(w io.Writer)(int64, error){//returns number of bytes written to the writer and an error interface
	err := binary.Write(w, binary.BigEndian, BinaryType) //1-byte type
	if err != nil{
		return 0, err
	}
	var n int64 = 1
	err = binary.Write(w, binary.BigEndian, uint32(len(m)))//writes 4-byte length of the Binary
	if err != nil{
		return n, err
	}
	n += 4
	o, err := w.Write(m)//writes the Binary value itself
	return n + int64(o), err
}

func (m *Binary) ReadFrom(r io.Reader)(int64, error){
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)//reads 1-byte from the reader
	if err !=nil{
		return 0, err
	}
	var n int64 = 1
	if typ != BinaryType{//verifies the type is BinaryType
		return n, errors.New("invalid Binary")
	}
	var size uint32
	err = binary.Read(r, binary.BigEndian, &size)//read 4-byte size
	if err != nil{
		return n, err
	}
	n += 4
	if size > MaxPayloadSize{//enforce maximum paylord size
		return n, ErrMaxPayloadSize
	}
	*m = make([]byte, size)//sizes the new Binary byte slice
	o, err := r.Read(*m)//payload
	 return n + int64(o), err 
		
}
type String string

func(m String) Bytes()[]byte {return []byte(m)}//casts the string to a byte slice
func(m String) String() string {return string(m)}//casts String type to base type String

func (m String) WriteTo(w io.Writer) (int64, error){
	err := binary.Write(w, binary.BigEndian, StringType)//1-nyte type
	if err != nil{
		return 0, err
	}
	var n int64 = 1
	err = binary.Write(w, binary.BigEndian, uint32(len(m)))
	if err !=nil{
		return n, err
	}
	n += 4
	o, err := w.Write([]byte(m))//payload
	return n + int64(o), err
}
func (m *String) ReadFrom(r io.Reader)(int64, error){
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)//1-byte type
	if err != nil{
		return 0, err
	}
	var n int64 = 1
	if typ != StringType{
		return n, errors.New("invalid string")
		
	}

	var size uint32
	err = binary.Read(r, binary.BigEndian, &size)//4-byte size
	if err != nil{
		return n, err
	}
	n += 4
	buf := make([]byte, size)
    o, err := r.Read(buf)//payload
    if err != nil{
    return n, err
    }
    *m = String(buf)
    return n + int64(o), nil
}
func decode(r io.Reader) (Payload, error){//return payload interface and error interface
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil{
		return nil, err
	}
	var payload Payload
	switch typ{
		case BinaryType:
		payload = new(Binary)
		case StringType:
		payload = new(String)
		default:
		return nil, errors.New("unknowm type")		
	}
	_, err = payload.ReadFrom(
		io.MultiReader(bytes.NewReader([]byte{typ}),r))
	if err != nil{
		return nil, err
	}
	return payload, nil
}
func TestPayloads(t *testing.T){
	b1 := Binary("Clear is better than clever.")
	b2 := Binary("Don't panic.")
	s1 := String("Errors are values.")
	payloads := []Payload{&b1,&s1,&b2}
	 listener, err := net.Listen("tcp", "127.0.0.1:")
		if err != nil{
			t.Fatal(err)
		}

		go func(){
			conn, err := listener.Accept()
			if err != nil{
				t.Error(err)
				return
			}
			defer conn.Close()
			for _, p := range payloads{
				_, err = p.WriteTo(conn)
				if err !=nil{
					t.Error(err)
					break
				}
			}
		}()


	//client side
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err !=nil{
		t.Fatal(err)
	}
	defer conn.Close()
	for i := 0; i < len(payloads); i++{
		actual, err := decode(conn)
		if err != nil{
			t.Fatal(err)
		}

		if expected := payloads[i]; !reflect.DeepEqual(expected,actual){
			t.Errorf("value mismatch: %v != %v", expected, actual)
			continue
			
		}
		t.Logf("[%T] %[1]q", actual)
	}
}

func TestMaxPayloadSize(t *testing.T){
	buf := new(bytes.Buffer)
	err := buf.WriteByte(BinaryType)
	if err != nil{
		t.Fatal(err)
	}
	err = binary.Write(buf, binary.BigEndian, uint32(1<<30))//1 GB
	if err != nil{
		t.Fatal(err)
	}
	var b Binary
	_, err = b.ReadFrom(buf)
	if err != ErrMaxPayloadSize{
		t.Fatalf("expected ErrMaxPayloadize; actual: %v", err)
	}
}
