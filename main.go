package main

import (
	"debug/elf"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/dustin/go-humanize"
)

func main() {
	outFile := ""
	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagset.StringVar(&outFile, "output", "debug.symbols", "The file path to put object file with extracted symbols.")

	err := flagset.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		fmt.Println("missing parameter, please provide  a file")
		return
	}
	path := os.Args[1]

	elfFile, err := elf.Open(path)
	if err != nil {
		log.Fatalf("error while opening ELF file %s: %s", path, err.Error())
	}
	defer elfFile.Close()

	type section struct {
		n string
		t string
		s uint64
	}
	var sections []section
	for _, s := range elfFile.Sections {
		sections = append(sections, section{n: s.Name, t: s.Type.String(), s: s.Size})
	}
	sort.Slice(sections, func(i, j int) bool {
		return sections[i].s > sections[j].s
	})

	var sum uint64
	for _, s := range sections {
		sum += s.s
		fmt.Println("Name:", s.n, "Size:", humanize.Bytes(s.s), "Type:", s.t)
	}
	fmt.Println("Total: ", humanize.Bytes(sum))
}
