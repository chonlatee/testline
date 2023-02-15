package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type fileStat struct {
	testFile     float64
	workFile     float64
	testFileLine float64
	workFileLine float64

	suffixTestFile string
	suffixWorkFile string
}

func main() {

	var root, suffixTestFile, suffixWorkFile string
	flag.StringVar(&root, "root", "", "-root=./")
	flag.StringVar(&suffixTestFile, "suffixTestFile", "_test.go", "-suffixTestFile=_test.go")
	flag.StringVar(&suffixWorkFile, "suffixWorkFile", ".go", "-suffixWorkFile=.go")
	flag.Parse()

	if root == "" {
		log.Fatal("need folder to walk")
	}

	fst := &fileStat{
		suffixTestFile: suffixTestFile,
		suffixWorkFile: suffixWorkFile,
	}

	filepath.Walk(root, fst.visit)

	log.Printf("test file %.0f\n", fst.testFile)
	log.Printf("test line %.0f\n", fst.testFileLine)

	log.Printf("work file %.0f\n", fst.workFile)
	log.Printf("work file line %.0f\n", fst.workFileLine)

	log.Printf("test line percentage %.2f%%\n", fst.testFileLine/fst.workFileLine*100)
	log.Printf("test file percentage %.2f%%\n", fst.testFile/fst.workFile*100)
}

func (fst *fileStat) visit(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	line, err := fst.countFile(path, info)
	if err != nil {
		return err
	}
	if strings.HasSuffix(path, fst.suffixWorkFile) &&
		!strings.HasSuffix(path, fst.suffixTestFile) {
		fst.countWorkFile(line)
	} else if strings.HasSuffix(path, fst.suffixTestFile) {
		fst.countTestFile(line)
	}

	return nil
}

func (fst *fileStat) countTestFile(line float64) {
	fst.testFileLine += line
	fst.testFile += 1
}

func (fst *fileStat) countWorkFile(line float64) {
	fst.workFileLine += line
	fst.workFile += 1
}

func (fst *fileStat) countFile(path string, info os.FileInfo) (float64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, nil
	}
	c, err := lineCounter(f)
	if err != nil {
		return 0, err
	}

	return float64(c), nil
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}
