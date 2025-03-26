package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"stockgame/internal/util"
)

func main() {
	// Open the constant.go file (from the internal/model package)
	file, err := os.Open(util.GetProjectRoot() + "/internal/model/constant.go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	contantsToAdd := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 5 && line[0:5] == "const" {
			contantsToAdd = append(contantsToAdd, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	// Save the content into dynamicConstants.ts file (in the ../cmd/frontend-server/src folder)
	f, err := os.Create(util.GetProjectRoot() + "/cmd/frontend-server/src/dynamicConstants.ts")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for _, line := range contantsToAdd {
		fmt.Fprintln(f, "export "+line+";")
	}
}
