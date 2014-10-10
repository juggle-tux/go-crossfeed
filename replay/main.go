

package main


import (
	//"bytes"
	"fmt"
	"flag"
	"log"
	"net"
	"os"

)

type Connection struct {
	ClientAddr *net.UDPAddr // Address of the client
	ServerConn *net.UDPConn // UDP connection to server
}



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

	// Create UDP connection
	addr_str := fmt.Sprintf("127.0.0.1:%d", *iport)
	_, err_conn := net.Dial("udp", addr_str )
	if err_conn != nil {
		logger.Println("Fail Conn", err_conn)
		return
	}

	// Lets go
	logger.Println("Starting with port: ",  addr_str, "  Log: ", *ifile)

	// Open cf.log file
	file, err_open := os.Open(*ifile)
	if err_open != nil {
		logger.Println("Failed open cflog: ", err_open)
		return
	}
	defer file.Close()

	// setup buffers
	data := make([]byte, 256)
	counter := 0
	for {
		//data = data[:cap(data)]
		n, err := file.Read(data)
		if err != nil {
			//if err == io.EOF {
				logger.Println("EOF: ", err)
				return
			//}
		}
		counter += 1
		//data = data[:n]
		logger.Println(counter, n, len(data), cap(data) )
		for i := 4; i < n; i++ {
			if data[i + 0] == 'S' && data[i + 1] == 'F'  && data[i + 2] == 'G'  && data[i + 3] == 'S'  {
				fmt.Print("YESSSSSSSSSS")
			}
		}

	}

	fmt.Println("Done bye !")
}
