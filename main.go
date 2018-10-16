package main

import (
	"crypto"
	_ "crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/frostschutz/go-fibmap"
)

func getExtents(filename string) ([]fibmap.Extent, int64, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}
	defer fd.Close()

	fm := fibmap.NewFibmapFile(fd)

	bsz, errno := fm.Figetbsz()
	if errno != 0 {
		return nil, 0, fmt.Errorf("figetbsz: %v", errno)
	}

	stat, err := fd.Stat()
	if err != nil {
		return nil, 0, fmt.Errorf("fstat: %v", err)
	}
	size := stat.Size()

	blocks := uint32((size-1)/int64(bsz)) + 1

	extents, errno := fm.Fiemap(blocks)
	if errno != 0 {
		return nil, 0, fmt.Errorf("fiemap: %v", errno)
	}
	return extents, size, nil
}

func getFIENode(filename string) (string, error) {
	extents, size, err := getExtents(filename)
	if err != nil {
		return "", err
	}

	w := crypto.SHA1.New()
	fmt.Fprintf(w, "%d\n", size)
	for _, e := range extents {
		fmt.Fprintf(w, "%d + %d\n", e.Physical>>12, e.Length>>12)
	}

	return hex.EncodeToString(w.Sum(nil)), nil
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("fienode: ")

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: fienode <filename>...")
		os.Exit(1)
	}

	for _, filename := range os.Args[1:] {
		fienode, err := getFIENode(filename)
		if err != nil {
			log.Fatal(err)
		}
		_, err = fmt.Println(fienode, filename)
		if err != nil {
			log.Fatal(err)
		}
	}
}
