package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"time"
)

var magicHashRegexp = regexp.MustCompile("^0+([eE]\\d+)?$")

func isMagicHash(b []byte) bool {
	if b[0]>>4 != 0 {
		// quick fail if first hex char is not 0
		return false
	}
	return magicHashRegexp.MatchString(hex.EncodeToString(b))
}

func doWork(data []byte, offset int, resCh chan<- []byte) {
	b := make([]byte, len(data)+offset+binary.MaxVarintLen64)
	copy(b, data)

	writePos := len(data) + offset
	var n uint64 = 1
	for ; ; n++ {
		if n == 0 {
			fmt.Printf("%d: overflow!", offset)
			return
		}
		binary.PutUvarint(b[writePos:], n)
		hashSum := sha1.Sum(b)
		if isMagicHash(hashSum[:]) {
			fmt.Printf("magicHash: %x [offset=%d, n=%d]\n", hashSum, offset, n)
			resCh <- b[len(data):]
			return
		}
	}
}

func main() {
	start := time.Now()
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %v image\n", os.Args[0])
		return
	}

	inFileName := os.Args[1]
	outFileName := inFileName + ".out"

	inFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(inFile); err != nil {
		log.Fatal(err)
	}

	// Spawn a couple of worker go routines that tries to find a
	// byte suffix that makes the total file bytes a magic hash.
	resCh := make(chan []byte)
	numWorkers := runtime.NumCPU()
	for i := 0; i < numWorkers; i++ {
		go doWork(buf.Bytes(), i, resCh)
	}

	// collect the first result
	res := <-resCh

	// Append to the buffer and write to output file
	if _, err := buf.Write(res); err != nil {
		log.Fatal(err)
	}
	hashSum := sha1.Sum(buf.Bytes())
	fmt.Printf("Found magic hash: %x\n", hashSum)
	fmt.Printf("Elapsed: %s\n", time.Since(start))
	if _, err := buf.WriteTo(outFile); err != nil {
		log.Fatal(err)
	}
}
