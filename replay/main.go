

package main


import (
	//"bytes"
	"fmt"
	"flag"
	"io"
	"log"
	"net"
	"os"
)



func main() {

	//= Command Args"
	var ihelp *bool = flag.Bool("h", false, "Show Help")
	var iport *int = flag.Int("p", 4444, "UDP Transmit port")
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
	conn, err_conn := net.Dial("udp", addr_str )
	if err_conn != nil {
		logger.Println("Fail UDP Connection", err_conn)
		return
	}

	// setup loop buffers _ vars
	data := make([]byte, 1024) // file data buffer
	counter := 0 // count no of reads for dev debug

	Fb := byte('F') // F,G,S as a byte (later maybe we can copare all four chars
	Gb := byte('G')
	Sb := byte('S')

	// loop forever, currently till eof or intentional crash
	for {

		// read from cf.log file and check for errors
		n, err := file.Read(data)
		if err != nil {
			if err == io.EOF {
				logger.Println("EOF: ", err)
				return
			}
		}
		counter += 1

		// loop the buffer and detect the 'SFGF'
		last := n - 4
		for i := 0; i < last; i++ {

			if data[i + 0] == Sb && data[i + 1] == Fb  && data[i + 2] == Gb  && data[i + 3] == Fb  {
				fmt.Println("YES found FGFS", counter,  i)
				// TODO
				conn.Write(data) // Write to UDP
			}
		}
		if counter == 10 {
			fmt.Println("Killed after a few")
			os.Exit(0)
		}
	}

	fmt.Println("Done bye !")
}
