

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
	logi := log.New(os.Stderr, "LOG: ", log.Lshortfile)

	// Check cf log file exists
	if _, err_file := os.Stat(*ifile); os.IsNotExist(err_file) {
		logi.Printf("No such file or directory: %s", *ifile)
		return
	}

	// Create UDP connection
	con, err_conn := net.Dial("udp", fmt.Sprintf("127.0.0.1:%d", *iport))
	if err_conn != nil {
		logi.Println("Fail Conn", err_conn)
		return
	}

	// Lets go
	fmt.Println("Starting: with port: ",  fmt.Sprintf("127.0.0.1:%d", *iport), "  Log: ")
	logi.Println("OK", con)



	//fmt.Println("Done bye !")
}
