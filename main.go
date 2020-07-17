package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Too few arguments, pass atleast 1 filename to be executed")
		panic(errors.New("No files passed"))
	}

	var files []byte
	for _, file := range os.Args[1:] {
		curr, err := ioutil.ReadFile(file)
		check(err)
		files = append(files, curr...)
	}

	data := string(files)

	prog := makeProgram(data)
	err := execute(prog)
	check(err)
}

type Program struct {
	size    int
	command []byte
	at      int
}

func makeProgram(data string) *Program {
	prog := new(Program)
	prog.command = []byte(data)
	prog.size = 1
	var index int = 1
	for _, char := range data {
		switch char {
		case '<':
			index--
		case '>':
			index++
		default:
			break
		}
		if prog.size < index {
			prog.size = index
		}
	}
	prog.at = 0
	return prog
}

func execute(prog *Program) error {
	storage := make([]byte, prog.size)
	read := bufio.NewReader(os.Stdin)
	nestLevel := 0
	for i := 0; i < len(prog.command); i++ {
		switch prog.command[i] {
		case '+':
			storage[prog.at]++
		case '-':
			storage[prog.at]--
		case '>':
			if prog.at+1 == prog.size {
				return errors.New("unable to access memory")
			} else {
				prog.at++
			}
		case '<':
			if prog.at == 0 {
				return errors.New("unable to access memory")
			} else {
				prog.at--
			}
		case '.':
			fmt.Printf("%c", storage[prog.at])
		case ',':
			fmt.Printf("input char: ")
			in, _, err := read.ReadRune()
			check(err)
			storage[prog.at] = byte(in)
		case '[':
			if storage[prog.at] == 0 {
				tempNestLevel := nestLevel
				i++
				for ; i < len(prog.command); i++ {
					if prog.command[i] == '[' {
						tempNestLevel++
					}
					if prog.command[i] == ']' {
						if tempNestLevel != nestLevel {
							tempNestLevel--
						} else {
							break
						}
					}
					if i == len(prog.command)-1 {
						return errors.New("Loop opened without closing, '[' without ']'")
					}
				}
			} else {
				nestLevel++
			}
		case ']':
			if storage[prog.at] != 0 {
				tempNestLevel := nestLevel
				i--
				for ; i >= 0; i-- {
					if prog.command[i] == ']' {
						tempNestLevel++
					}
					if prog.command[i] == '[' {
						if tempNestLevel != nestLevel {
							tempNestLevel--
						} else {
							break
						}
					}
					if i == 0 {
						return errors.New("Loop closed without opening, ']' without '['")
					}
				}
			} else {
				nestLevel--
			}
		default:
		}
	}
	return nil
}
