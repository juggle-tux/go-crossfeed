
package message

import (
	"fmt"

	"github.com/davecgh/go-xdr/xdr"
)

// Magic value for messages - currently FGFS
const MSG_MAGIC = 0x46474653  // "FGFS"
//const MSG_MAGIC = "SFGF"  // "FGFS"



// Protocol Version - currently 1.1
const PROTO_VER = 0x00010001  // 1.1

// Message Types
const (
	CHAT_MSG_ID = 1 //= is this used ??
	RESET_DATA_ID = 6
	POS_DATA_ID = 7
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
	Id uint32 //xdr_data_t

	// Absolute length of message
	Len uint32 //xdr_data_t

	// DEPRECEATED: Player's receiver address
	ReplyAddress uint32 //xdr_data_t

	// DEPRECEATED: Player's receiver port
	ReplyPort uint32 //xdr_data_t

	// Callsign used by the player
	Callsign [MAX_CALLSIGN_LEN]byte //Callsign[MAX_CALLSIGN_LEN]
}

// return callsign as string
// TODO There's got to be a better way
func (me *Header) CallsignString() string{
	s := ""
	for _, ele := range me.Callsign {
		if ele == 0 {
			return s
		}
		s += string(ele)
	}
	return s
}


func Decode(bits []byte)(Header, error) {

	var header Header

	remainingBytes, err := xdr.Unmarshal(bits, &header)
	if err != nil{
		fmt.Println("XDR Decode Error", err)
		return header, nil
	}
	fmt.Println("remain=", len(remainingBytes))
	fmt.Println ("Header=", header.Id, header.Version, header.CallsignString(), header.Magic, MSG_MAGIC)

	return header, nil
}
