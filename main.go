package main

import (
	"debug/elf"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/table"
)

func main() {
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

	type sectionInfo struct {
		Name       string
		Type       string
		Size       uint64
		Flags      string
		Link       string
		Compressed bool
	}

	idxSection := make(map[int]string)
	for i, s := range elfFile.Sections {
		idxSection[i] = s.Name
	}

	var sections []sectionInfo
	for _, s := range elfFile.Sections {
		sections = append(sections, sectionInfo{
			Name:       s.Name,
			Type:       s.Type.String(),
			Size:       s.Size,
			Flags:      s.Flags.String(),
			Link:       idxSection[int(s.Link)],
			Compressed: s.Flags&elf.SHF_COMPRESSED != 0,
		})
	}
	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Size > sections[j].Size
	})

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetAutoIndex(true)
	t.AppendHeader(table.Row{"Name", "Type", "Size", "Flags", "Link", "Compressed"})
	var sum uint64
	for _, s := range sections {
		sum += s.Size
		t.AppendRow(table.Row{
			s.Name, s.Type, humanize.Bytes(s.Size), s.Flags, s.Link, s.Compressed,
		})
	}
	t.AppendFooter(table.Row{"", "Total Size", humanize.Bytes(sum)})
	t.Render()
}
