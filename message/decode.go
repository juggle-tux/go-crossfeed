
package message

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/davecgh/go-xdr/xdr"
)


func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, msInt*int64(time.Millisecond)), nil
}


// Decode the XDR packet
func Decode(xdr_enc []byte)(HeaderMsg, error) {

	var header HeaderMsg

	remainingBytes, err := xdr.Unmarshal(xdr_enc, &header)
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
	t := time.Unix(0, int64(position.Time) * int64(time.Nanosecond))
	//t2 := time.Unix(int64(position.Time), 0)
	fmt.Println(position.Model(), position.Time, ">>", t, err)
	return header, nil
}
