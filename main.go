package main

import (
	"debug/elf"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/table"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing parameter, please provide  a file")
		return
	}
	path := os.Args[1]

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("error while opening file %s: %s", path, err.Error())
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		log.Fatalf("error while reading file stats %s: %s", path, err.Error())
	}

	elfFile, err := elf.NewFile(f)
	if err != nil {
		log.Fatalf("error while opening ELF file %s: %s", path, err.Error())
	}
	defer elfFile.Close()

	type sectionInfo struct {
		*elf.Section

		Name       string
		Type       string
		Offset     uint64
		Size       uint64
		FileSize   uint64
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
			Section: s,

			Name:       s.Name,
			Type:       s.Type.String(),
			Offset:     s.Offset,
			Size:       s.Size,
			FileSize:   s.FileSize,
			Flags:      s.Flags.String(),
			Link:       idxSection[int(s.Link)],
			Compressed: s.Flags&elf.SHF_COMPRESSED != 0,
		})
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetAutoIndex(true)
	t.AppendHeader(table.Row{"Name", "Type", "Offset", "Size", "File Size", "Flags", "Link", "Compressed"})
	var (
		sum   uint64
		fzSum uint64
	)
	for _, s := range sections {
		sum += s.Size
		size, err := fileSize(f, s.Section)
		if err != nil {
			log.Printf("failed to read the size of %s: %s\n", s.Name, err.Error())
		}
		fzSum += size
		t.AppendRow(table.Row{
			s.Name, s.Type, s.Offset, humanize.Bytes(s.Size), humanize.Bytes(s.FileSize), s.Flags, s.Link, s.Compressed,
		})
	}
	t.AppendFooter(
		table.Row{"", "", "Total", humanize.Bytes(sum), humanize.Bytes(fzSum), "Actual File Size", humanize.Bytes(uint64(stat.Size()))},
	)
	t.Render()
}

func fileSize(src io.ReaderAt, sec *elf.Section) (uint64, error) {
	if sec.Type == elf.SHT_NOBITS {
		return 0, nil
	}
	var rs io.ReadSeeker
	if sec.Flags&elf.SHF_COMPRESSED != 0 {
		// Compressed
		rs = io.NewSectionReader(src, int64(sec.Offset), int64(sec.FileSize))
	} else {
		rs = sec.Open()
	}
	b, err := io.ReadAll(rs)
	if err != nil {
		return 0, err
	}
	return uint64(len(b)), nil
}
