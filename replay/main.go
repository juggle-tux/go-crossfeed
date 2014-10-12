/*
This reads an UDP huge dump of flightear multiplayer data,
and replays it out to an UDP socket.
- A snipped (100k) test file is at stuff/cf_test.log
- Original (132mb) at  http://geoffair.org/tmp/cf_raw01.zip
The packets all start with the "magic"
`SFGF` ie FGFS (Flight Gear Flight Simulator).

The raw log was written packet
by packet, so the first packet in the
log it the first received, next received,
and so on... So it is -

<magic - packet-1>
<magic - packet-2>
<magic - packet-3>
... and so on...

The idea of this app, is to read the file,
detect the "packets", and then fire them off
into an UDP socket, like a replay.

TODO: Help wanted..
this code kinda works, but am a complete newbie
so maybe there isa  better way..

Needs to "transmit" at 10 - 25 hz per packet

*/
/* thanks to:
 geoffmcl@github - for being an expert in FlightGear
 juggle-tux@github - for golang help
 */
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/FreeFlightSim/go-crossfeed/message"
)

func main() {

	// Args"
	var ihelp *bool = flag.Bool("h", false, "Show Help")
	var iport *int = flag.Int("p", 3333, "UDP Transmit port")
	var ifile *string = flag.String("l", "../stuff/cf_test.log", "Path to raw log file")
	var ihz *int = flag.Int("z", 10, "Hz to transmit")
	flag.Parse()
	if *ihelp {
		flag.Usage()
		os.Exit(0)
	}

	// Setup logger
	logger := log.New(os.Stderr, "Replay: ", log.Lshortfile)

	// Check cf.log file exists
	if _, err := os.Stat(*ifile); os.IsNotExist(err) {
		logger.Printf("No such file or directory: %s", *ifile)
		return
	}

	// Open cf.log file
	file, err := os.Open(*ifile)
	if err != nil {
		logger.Println("Failed open cflog: ", err)
		return
	}
	defer file.Close()

	// Create UDP socket
	addr_str := fmt.Sprintf("127.0.0.1:%d", *iport)
	sock, err := net.Dial("udp4", addr_str)
	if err != nil {
		logger.Println("Fail UDP Connection", err)
		return
	}

	// setup buffers +  vars
	buffer := make([]byte, 4096) // file data buffer
	data := make([]byte, 0)      // temp buffer
	read_counter := 0            // count no of reads for dev debug
	packets := 0                 // packets processed

	// magic
	magic := []byte("SFGF")

	// timer for 10-25Hz
	var pulse time.Duration = time.Duration(int64(1000 / *ihz)) * time.Millisecond
	ticker := time.Tick(pulse)


	// loop forever, currently till eof or intentional crash
	for {

		// read from cf.log file into buffer and check for errors
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				logger.Println("EOF: ", err)
				return
			}
			logger.Println("error while reading ", err)
		}
		read_counter += 1

		// append the buffer to existing/remainder data ([] at first loop)
		data = append(data, buffer[:n]...)

		// find first and last magic  'SFGF'
		// (first should always be zero - so this is for testing)
		first := bytes.Index(data, magic)
		if first != 0 {
			fmt.Println("Strange, first magic not at zero")
		}
		last := bytes.LastIndex(data, magic)

		// hack off the `bits` within first to last
		bits := data[first : last-1]

		// hackoff data, leaving remainder
		data = data[last:]

		// loops the bits, finding magic
		packet_start := 0
		for i := 4; i < len(bits)-4; i++ {
			if bits[i] == magic[0] && bytes.Equal(bits[i:i+4], magic) {
				<-ticker
				packet, errp := message.Decode(bits[packet_start:i])
				if errp != nil {
					fmt.Println(errp)
				}else if (1 == 2) {
					fmt.Println("DECO+", packet, errp, packet.Type)
					//if packet.Id == message.POS_DATA_ID {

					//}
				}
				_, err := sock.Write(bits[packet_start:i])
				if err != nil {
					logger.Println(err)
				}
				packets += 1
				packet_start = i
				logger.Println(packets)
			}
		}

		if read_counter%2000 == 0 {
			fmt.Print(read_counter, " ")
			//os.Exit(0)
		}
	}

	// todo - restart again from top
	fmt.Println("Done bye !")
}
