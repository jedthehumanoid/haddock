// Package mainasn
//
// asfd
package main

import (
	"encoding/json"
	"fmt"
	gitignore "github.com/sabhiram/go-gitignore"
	"io/ioutil"
	"os"
	"path/filepath"
	"os/exec"
	"strings"
)



const gray = "\033[38;2;100;100;100m"
const reset = "\033[0m"

func readgo(file string) string {
	docCmd := exec.Command("go", "doc", "-u", "-all", file)
	docOut, err := docCmd.Output()
	    if err != nil {
	        panic(err)
	    }

	return strings.TrimSpace(string(docOut))
}

func parseGo(file string, long bool) string {
	ret := ""
	contents := readgo(file)
	if !long {
		if strings.Contains(contents, "CONSTANTS") {
			ret = strings.TrimSpace(strings.Split(contents, "CONSTANTS")[0])
		} else 	if strings.Contains(contents, "VARIABLES") {
			ret = strings.TrimSpace(strings.Split(contents, "VARIABLES")[0])
		} else 	if strings.Contains(contents, "FUNCTIONS") {
			ret = strings.TrimSpace(strings.Split(contents, "FUNCTIONS")[0])
		}
		ret = strings.Split(ret, "\n\n")[0]
	}

	return ret
}

var filetypes = map[string]func(string, bool)string{
	".go": parseGo,
}

func main() {

	file := "."
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	s, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No such file or directory:", file)
			os.Exit(0)
		}
		panic(err)
	}
	
	if s.IsDir() {
		fmt.Println("Reading as directory...\n")
		fmt.Println(gray+" ---------------"+reset)
		for _, line := range strings.Split(dirLong(file), "\n") {
			fmt.Println(gray+"| "+reset+line)
		}
		fmt.Println(gray+" ---------------"+reset)
		fmt.Println()
		files := find(file)
		for _, file := range files {
			ext := filepath.Ext(file)
			fn, ok := filetypes[ext]
			if ok {
				fmt.Println()
				fmt.Println(gray +fn(file, false)+reset)
				fmt.Println(file)
				
			} else {
				fmt.Println(file)
			}
		}
	}

}

func dirLong(dir string) string {
	b, err := ioutil.ReadFile(dir+"/README.md")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func dirShort(dir string) {

}

func find(dir string) []string {
	ret := []string{}

	ignore := false
	var ignored *gitignore.GitIgnore
	if exists(".gitignore") {
		var err error
		ignored, err = gitignore.CompileIgnoreFile(".gitignore")
		if err != nil {
			panic(err)
		}
		ignore = true
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {

	}

	for _, file := range files {
		//fmt.Println(file.Name())
		if file.IsDir() && file.Name() != ".git" {
			if ignore {
				if !ignored.MatchesPath(file.Name()) {
					ret = append(ret, file.Name()+"/")
				}
			} else {
				ret = append(ret, file.Name()+"/")
			}
		}
	}

	for _, file := range files {
		if !file.IsDir() && file.Name() != ".gitignore" {
			if ignore {
				if !ignored.MatchesPath(file.Name()) {
					ret = append(ret, file.Name())
				}
			} else {
				ret = append(ret, file.Name())
			}
		}
	}
	return ret
}

func toJSONPretty(i interface{}) string {
	ret, err := json.MarshalIndent(i, "", "   ")
	if err != nil {
		panic(err)
	}
	return string(ret)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
