package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("benchOut2.txt")
	if err != nil {
		log.Fatal(err)
	}
	outF, err := os.Create("benchOut2.csv")
	if err != nil {
		log.Fatal(err)
	}
	csvOut := csv.NewWriter(outF)
	bufIn := bufio.NewReader(f)
	for {
		line, err := bufIn.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		actualLine := []byte{}
			for _, b := range line {
				if b != 0 {
					actualLine = append(actualLine, b)
				}
			}
		strs := strings.Fields(line)
		if strs[0] == "PASS" {
			break
		}
		fmt.Println(strs[0])
		csvStrings := make([]string, 9)
		csvStrings[0] = strs[1]
		for i := 0; i < 2; i++ {
			line, err := bufIn.ReadBytes('\n')
			if err != nil {
				log.Fatal(err)
			}
			actualLine := []byte{}
			for _, b := range line {
				if b != 0 {
					actualLine = append(actualLine, b)
				}
			}
			strs := strings.Fields(string(actualLine))
			fmt.Println(len(strs), strs[len(strs)-2])
			csvStrings[i+1] = strs[len(strs)-2]
		}
		err = csvOut.Write(csvStrings)
		if err != nil {
			log.Fatal(err)
		}
	}
	csvOut.Flush()
}
