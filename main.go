package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	doc    = "doc.go"
	readme = "README.md"
)

func main() {
	md, err := scanDoc()
	if err != nil {
		log.Fatal(err)
	}

	out := clean(md)

	if _, err := os.Stat(readme); err == nil {
		if !yesOrNo(fmt.Sprintf("%s already exists. Overwrite [Y/n]: ", readme)) {
			log.Fatal("Exiting without writing README.md")
		}
	}

	if err := ioutil.WriteFile(readme, []byte(out), 0755); err != nil {
		log.Fatal(err)
	}
}

func scanDoc() ([]string, error) {
	f, err := os.Open(doc)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	md := []string{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		t := s.Text()
		prefix := ""
		for _, p := range []string{"//", "/*", "*/", " */", "/**", "**/"} {
			if strings.HasPrefix(t, p) {
				prefix = p
			}
		}
		if len(prefix) > 0 {
			md = append(md, strings.TrimLeft(strings.TrimPrefix(t, prefix), " "))
		} else {
			md = append(md, strings.TrimLeft(t, " "))
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return md, nil
}

func clean(md []string) string {
	out := ""
	for i := 0; i < len(md); i++ {
		// remove empty first line
		if i == 0 && len(md[i]) == 0 {
			continue
		}
		// remove package statement
		if i == len(md)-1 {
			break
		}
		out += md[i] + "\n"
	}
	return out
}

func yesOrNo(text string) bool {
	fmt.Print(text)
	answer := ""
	fmt.Scanln(&answer)
	if len(answer) == 0 || strings.ToLower(answer) == "y" {
		return true
	}

	return false
}
