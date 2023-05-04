package util

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// 从字符串读取
func ReadLines(content string) (lines []string, err error) {
	br := bufio.NewReader(strings.NewReader(content))

	lines = make([]string, 0)
	for lineEnd := true; ; {
		lineBytes, isPrefix, err1 := br.ReadLine()
		if err1 != nil {
			if err1 != io.EOF {
				err = err1
			}
			break
		}

		line := string(lineBytes)
		if lineEnd == false {
			lines[len(lines)-1] += line

		} else {
			lines = append(lines, line)
			lineEnd = !isPrefix
		}
	}

	return
}

// 从文件读取
func ReadLinesFormFile(path string) (lines []string, err error) {
	file, _ := os.OpenFile(path, os.O_RDONLY, 0666)
	defer file.Close()

	br := bufio.NewReader(file)

	lines = make([]string, 0)
	for lineEnd := true; ; {
		lineBytes, isPrefix, err1 := br.ReadLine()
		if err1 != nil {
			if err1 != io.EOF {
				err = err1
			}
			break
		}

		line := string(lineBytes)
		if lineEnd == false {
			lines[len(lines)-1] += line

		} else {
			lines = append(lines, line)
			lineEnd = !isPrefix
		}
	}

	return
}

func ReadLine(fileName string) ([]string, error) {
	var readSource io.Reader
	if len(fileName) > 0 {
		f, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		readSource = f
	} else {
		readSource = os.Stdin
	}

	fileScanner := bufio.NewScanner(readSource)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string
	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}
	return fileTextLines, nil
}
