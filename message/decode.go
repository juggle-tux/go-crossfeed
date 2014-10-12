
package message

import (
	"fmt"
	"errors"

	"github.com/davecgh/go-xdr/xdr"
)




func Decode(bits []byte)(HeaderMsg, error) {

	var header HeaderMsg

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
	var position PositionMsg
	rembits, err := xdr.Unmarshal(remainingBytes, &position)

	if err != nil {
		fmt.Println(rembits)
	}
	fmt.Println(position.Model(), position.Time, err)
	return header, nil
}
