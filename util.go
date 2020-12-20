package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func mergeFiles(files []string, target string) {
	fs, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("合并%d个文件到%s\n", len(files), target)
	init := false

	for _, file := range files {
		fi, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		br := bufio.NewReader(fi)
		for {

			b, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			line := string(b)
			if strings.HasPrefix(line, "syntax ") ||
				strings.HasPrefix(line, "import ") ||
				strings.HasPrefix(line, "option ") ||
				strings.HasPrefix(line, "package ") {
				continue
			}
			if !init {
				if target == "tmpmsg.proto" {
					fs.WriteString("syntax=\"proto3\";\n" +
						"option go_package=\"/proto/proto_csmsg\";\n" +
						"package proto_csmsg;\nimport \"tmpobj.proto\";\n")

				} else {
					fs.WriteString("syntax=\"proto3\";\n" +
						"option go_package=\"/proto/proto_object\";\n" +
						"package proto_object;\n")
				}
				init = true
			}
			fs.WriteString(line + "\n")
		}
		// fs.WriteString("\n")

	}
	fs.Close()
}

func getFiles(paths ...string) []string {
	files := make([]string, 0)
	for _, path := range paths {
		filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			files = append(files, path+"/"+info.Name())
			return nil
		})
	}
	return files
}
