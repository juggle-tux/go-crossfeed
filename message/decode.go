
package message

import (
	"fmt"
	"errors"

	"github.com/davecgh/go-xdr/xdr"
)

// Magic value for messages - currently FGFS
const MSG_MAGIC = 0x46474653  // "FGFS"
//const MSG_MAGIC = "SFGF"  // "FGFS"



// Protocol Version - currently 1.1
const PROTOCOL_VER = 0x00010001  // 1.1

// Message Types
const (
	TYPE_CHAT = 1 //= is this used ??
	TYPE_RESET = 6
	TYPE_POS = 7
)



/*
	XDR demands Id4 byte alignment, but some compilers use 8 byte alignment
	so it's safe to let the overall size of a network message be a
	multiple of 8!
*/
const (
	MAX_CALLSIGN_LEN	= 8
	MAX_CHAT_MSG_LEN   	= 256
	MAX_MODEL_NAME_LEN 	= 96
	MAX_PROPERTY_LEN   	= 52
)


type Header struct {

	// Magic Value
	Magic uint32 //xdr_data_t

	// Protocol version
	Version uint32 //xdr_data_t

	// Message identifier
	Type uint32 //xdr_data_t

	// Absolute length of message
	Len uint32 //xdr_data_t

	// DEPRECEATED: Player's receiver address
	ReplyAddress uint32 //xdr_data_t

	// DEPRECEATED: Player's receiver port
	ReplyPort uint32 //xdr_data_t

	// Callsign used by the player
	Callsign [MAX_CALLSIGN_LEN]byte //Callsign[MAX_CALLSIGN_LEN]
}

// returns Callsign as string
func (me *Header) CallsignString() string{
	return string(me.Callsign[:])
}


func Decode(bits []byte)(Header, error) {

	var header Header

	remainingBytes, err := xdr.Unmarshal(bits, &header)
	if err != nil{
		fmt.Println("XDR Decode Error", err)
		return header, nil
	}
	//fmt.Println("remain=", len(remainingBytes))
	//fmt.Println( header.Magic == MSG_MAGIC, header.Version ==  PROTOCOL_VER)
	fmt.Println ("Header=", len(remainingBytes), header.Type, header.Type == TYPE_POS, header.Version, header.CallsignString(), )

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
