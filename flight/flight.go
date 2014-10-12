
package flight

type Flight struct {


	//Origin string
	Address *net.UDPAddr
	//Conn *net.UDPConn

	Callsign string // But this is also key so maybe unneeded ?
	Model string

	JoinTime int64 // epoch
	Timestamp int64 // epoch

	LastPos Point3D
	LastOrientation Point3D

	IsLocal bool

	Error string //;    // in case of errors
	HasErrors bool

	ClientID int
	LastRelayedToInactive int64

	// Packets recieved from client
	PktsReceivedFrom uint

	// Packets sent to client
	PktsSentTo uint

	// Packets from client sent to other players/relays
	PktsForwarded uint


}
