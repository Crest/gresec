package main

import "net"
import "fmt"
import "os"
import "io"
import "bytes"
import "strings"

const (
	fieldSep   = " "
	fieldCount = 4
	maxNameLen = 255
	sizeIP6    = 16
	sizeIP4    = 4
	sizeAddrs  = 2*sizeIP6 + 1*sizeIP4
	maxSize    = 1 + maxNameLen + sizeAddrs
)

var (
	NameTooShort  = os.NewError("Name too short")
	NameTooLong   = os.NewError("Name too long")
	ExtIPInvalid  = os.NewError("External IP address is invalid")
	IntIP4Invalid = os.NewError("Internal IPv4 address is invalid")
	IntIP6Invalid = os.NewError("Internal IPv6 address is invalid")
)


type Node struct {
	Name                  string
	ExtIP, IntIP4, IntIP6 net.IP
}

func NewNode(name string, extIP, intIP4, intIP6 net.IP) *Node {
	return &Node{Name: name, ExtIP: extIP, IntIP4: intIP4, IntIP6: intIP6}
}

func NewNodeString(line string) (*Node, os.Error) {
	line = strings.TrimSpace(line)
	fields := strings.Split(line, fieldSep, fieldCount)
	if len(fields) != fieldCount {
		return nil, os.NewError("Expected four space seperated fields")
	}

	name := fields[0]

	extIP := net.ParseIP(fields[1])
	if extIP == nil {
		return nil, os.NewError("External IP address is invalid")
	}

	intIP4 := net.ParseIP(fields[2])
	if intIP4 == nil {
		return nil, os.NewError("Interal IPv4 address is invalid")
	}
	intIP4 = intIP4.To4()
	if intIP4 == nil {
		return nil, os.NewError("Interal IPv4 address is a non IPv4 IP address")
	}

	intIP6 := net.ParseIP(fields[3])
	if intIP6 == nil {
		return nil, os.NewError("Internal IPv6 address is invalid")
	}
	if intIP6.To4() != nil {
		return nil, os.NewError("Internal IPv6 address is an IPv4 address")
	}

	return NewNode(name, extIP, intIP4, intIP6), nil
}

func NewNodeReader(r io.Reader) (*Node, os.Error) {
	buf := make([]byte, maxSize)
	if _, err := r.Read(buf[0:1]); err != nil {
		return nil, err
	}
	nameLen := int(buf[0])
	if nameLen == 0 {
		return nil, NameTooShort
	}
	size := 1 + nameLen + sizeAddrs
	if _, err := r.Read(buf[1:size]); err != nil {
		return nil, err
	}
	pos := 1
	name := string(buf[pos : pos+nameLen])
	pos += nameLen
	extIP := net.IP(buf[pos : pos+sizeIP6])
	pos += sizeIP6
	intIP4 := net.IP(buf[pos : pos+sizeIP4])
	pos += sizeIP4
	intIP6 := net.IP(buf[pos : pos+sizeIP6])
	pos += sizeIP6

	return NewNode(name, extIP, intIP4, intIP6), nil
}

func (node *Node) ToBuffer() (*bytes.Buffer, os.Error) {
	name := node.Name
	nameLen := len(name)
	buf := bytes.NewBufferString("")
	if nameLen > maxNameLen {
		return nil, NameTooLong
	}
	if nameLen < 1 {
		return nil, NameTooShort
	}
	extIP := node.ExtIP.To16()
	intIP4 := node.IntIP4.To4()
	intIP6 := node.IntIP6.To16()
	if extIP == nil {
		return nil, ExtIPInvalid
	}
	if intIP4 == nil {
		return nil, IntIP4Invalid
	}
	if intIP6 == nil {
		return nil, IntIP6Invalid
	}

	if err := buf.WriteByte(byte(nameLen)); err != nil {
		return nil, err
	}
	if _, err := buf.Write([]byte(name)); err != nil {
		return nil, err
	}
	if _, err := buf.Write(extIP); err != nil {
		return nil, err
	}
	if _, err := buf.Write(intIP4); err != nil {
		return nil, err
	}
	if _, err := buf.Write(intIP6); err != nil {
		return nil, err
	}
	return buf, nil
}

func (node *Node) String() string {
	return node.Name + " " + node.ExtIP.String() + " " + node.IntIP4.String() + " " + node.IntIP6.String()
}

func scan(node *Node, name string, extIP, intIP4, intIP6 net.IP) os.Error {
	if len(name) < 1 {
		return NameTooShort
	}
	if len(name) > maxNameLen {
		return NameTooLong
	}

	fmt.Print("name = ", name)
	fmt.Print(", ext  = ", extIP)
	fmt.Println()

	if extIP == nil {
		return ExtIPInvalid
	}
	if intIP4 == nil {
		return IntIP4Invalid
	}
	intIP4 = intIP4.To4()
	if intIP4 == nil {
		return IntIP4Invalid
	}
	if intIP6 == nil {
		return IntIP6Invalid
	}
	intIP6 = intIP6.To16()
	if intIP6 == nil {
		return IntIP6Invalid
	}

	node.Name = name
	node.ExtIP = extIP
	node.IntIP4 = intIP4
	node.IntIP6 = intIP6

	return nil
}

func dup(bytes []byte) []byte {
	copied := make([]byte, len(bytes))
	copy(copied, bytes)
	return copied
}

func (node *Node) Scan(state fmt.ScanState, verb int) os.Error {
	nameBytes, err := state.Token(true, nil)
	if err != nil {
		return err
	}
	if len(nameBytes) == 0 {
		return os.EOF
	}
	nameBytes = dup(nameBytes)

	extIPBytes, err := state.Token(true, nil)
	if err != nil {
		return err
	}
	extIPBytes = dup(extIPBytes)

	intIP4Bytes, err := state.Token(true, nil)
	if err != nil {
		return err
	}
	intIP4Bytes = dup(intIP4Bytes)

	intIP6Bytes, err := state.Token(true, nil)
	if err != nil {
		return err
	}
	intIP6Bytes = dup(intIP6Bytes)

	name := string(nameBytes)
	if len(name) < 1 {
		return NameTooShort
	}
	if len(name) > maxNameLen {
		return NameTooLong
	}

	extIP := net.ParseIP(string(extIPBytes))
	if extIP == nil {
		return ExtIPInvalid
	}

	intIP4 := net.ParseIP(string(intIP4Bytes))
	if intIP4 == nil {
		return IntIP4Invalid
	}
	intIP4 = intIP4.To4()
	if intIP4 == nil {
		return IntIP6Invalid
	}

	intIP6 := net.ParseIP(string(intIP6Bytes))
	if intIP6 == nil {
		return IntIP6Invalid
	}
	if intIP6.To4() != nil {
		return IntIP6Invalid
	}

	node.Name = name
	node.ExtIP = extIP
	node.IntIP4 = intIP4
	node.IntIP6 = intIP6

	return nil
}
