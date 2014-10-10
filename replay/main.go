

package main


import (
	"bytes"
	"fmt"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"time"
)



func main() {

	//= Command Args"
	var ihelp *bool = flag.Bool("h", false, "Show Help")
	var iport *int = flag.Int("p", 3333, "UDP Transmit port")
	var ifile *string = flag.String("l", "./cf_raw.log", "Path to raw log file")

	flag.Parse()
	if *ihelp {
		flag.Usage()
		os.Exit(0)
	}

	// Setup logger
	logger := log.New(os.Stderr, "LOG: ", log.Lshortfile)

	// Check cf.log file exists
	if _, err_file := os.Stat(*ifile); os.IsNotExist(err_file) {
		logger.Printf("No such file or directory: %s", *ifile)
		return
	}

	// Open cf.log file
	file, err_open := os.Open(*ifile)
	if err_open != nil {
		logger.Println("Failed open cflog: ", err_open)
		return
	}
	defer file.Close()

	// Create UDP connection
	addr_str := fmt.Sprintf("127.0.0.1:%d", *iport)
	conn, err_conn := net.Dial("udp4", addr_str )
	if err_conn != nil {
		logger.Println("Fail UDP Connection", err_conn)
		return
	}
	logger.Println("UDP Conn", addr_str)
	// setup loop buffers _ vars
	buffer := make([]byte, 4096) // file data buffer
	data := make([]byte, 0) // working buffer
	read_counter := 0 // count no of reads for dev debug
	packets := 0

	Fb := byte('F') // F,G,S as a byte (later maybe we can compare all four chars
	Gb := byte('G')
	Sb := byte('S')
	magic := []byte("SFGF")

	// loop forever, currently till eof or intentional crash
	for {

		// read from cf.log file and check for errors
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				logger.Println("EOF: ", err)
				return
			}
		}
		read_counter += 1
		//fmt.Println("-----------------------", read_counter, time.Now().UTC())

		// append the buffer to existing data ([] at first loop)
		data = append(data, buffer[:n]...)

		// find first and last magic (first should always be zero
		first := bytes.Index(data, magic)
		last := bytes.LastIndex(data, magic)

		// hack off the `bits` within first to last
		bits := data[first : last - 1]

		// hackoff data, leaving remainder
		data = data[last:]

		packet_start := 0
		for i := 4; i < len(bits) - 4; i++ {
			if bits[i] == Sb && bits[i + 1] == Fb  && bits[i + 2] == Gb  && bits[i + 3] == Fb  {
				_, errw := conn.Write( bits[packet_start:i] )
				if errw != nil {
					logger.Println(errw)
				}
				packets += 1
				packet_start = i
			}
		}
		time.Sleep(55 * time.Millisecond)

		if read_counter % 2000 == 0{
			fmt.Print(read_counter, " ")
			//os.Exit(0)
		}
	}

	fmt.Println("Done bye !")
}
