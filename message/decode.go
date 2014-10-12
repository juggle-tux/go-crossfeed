
package message

import (
	"fmt"
	"errors"

	"github.com/davecgh/go-xdr/xdr"
)




func Decode(bits []byte)(Header, error) {

	var header Header

	remainingBytes, err := xdr.Unmarshal(bits, &header)
	if err != nil{
		fmt.Println("XDR Decode Error", err)
		return header, nil
	}
	//fmt.Println("remain=", len(remainingBytes))
	//fmt.Println( header.Magic == MSG_MAGIC, header.Version ==  PROTOCOL_VER)
	fmt.Println ("Header=", len(remainingBytes), header.Type, header.Type == TYPE_POS, header.Version, header.Callsign(), )

	if header.Version != PROTOCOL_VER {
		return header, errors.New("Invalid protocol version")
	}
	if header.Type != TYPE_POS {
		return header, errors.New("Not a position error")
	}
	var tt uint64
	rembits, err := xdr.Unmarshal(remainingBytes, &tt)
	fmt.Println(tt, err)
	if err != nil {
		fmt.Println(rembits)
	}

	return header, nil
}
