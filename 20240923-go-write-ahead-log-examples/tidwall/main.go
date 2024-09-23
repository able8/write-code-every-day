package main

import (
	"fmt"

	"github.com/tidwall/wal"
)

func main() {
	// open a new log file
	log, err := wal.Open("mylog", nil)
	fmt.Println(err)

	// write some entries
	err = log.Write(1, []byte("first entry"))
	err = log.Write(1, []byte("first entry"))
	err = log.Write(2, []byte("second entry"))
	err = log.Write(3, []byte("third entry"))

	// read an entry
	data, err := log.Read(1)
	println(string(data)) // output: first entry

	// close the log
	err = log.Close()
}
