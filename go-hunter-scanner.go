//Akshy Palanisamy
//Lexical Analyzer for Hunter-Power-Gramar

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	//Checking for test file input
	if len(os.Args) <= 1 {
		fmt.Println("Please input a test file.")
		return
	}
	//retriving test file name
	readFileName := os.Args[1]

	//Alerting the user that the file is being processed
	fmt.Println("Processing Input File:", readFileName)

	//Opening the file for reading and outputting an error if the file does not open
	//Additionally the closing of the file is defered to the end of reading
	rFile, rErr := os.Open(readFileName)
	if rErr != nil {
		fmt.Println(rErr)
	}
	defer rFile.Close()

	//using the same filename given as the filename for the output
	writeFileName := strings.TrimSuffix(readFileName, ".txt")
	writeFileName = writeFileName + ".out"

	//Creating a file to write to and outpputting an error if the file is not created
	//Additionally the closing of the file is defered to the end of writing
	wFile, wErr := os.Create(writeFileName)
	if wErr != nil {
		fmt.Println(wErr)
	}
	defer wFile.Close()

	//Creating a buffered reader and writer to read from the input file and write to the output file
	//using bufio provides buffereing and help for textual I/O
	f := bufio.NewReader(rFile)
	w := bufio.NewWriter(wFile)

	//Reading the first rune in the input file
	r, _, fErr := f.ReadRune()

	//initializing the number of tokens to 0
	tokens := 0

	//This loop continues until the rune reaches EOF
	for fErr != io.EOF {

		//Switch statment for identifying tokens
		switch string(r) {

		//for every except ID case the tag is identified and then the variable by using the " "
		//deliminator, then the id is written to the output file using the required format
		//After each write the buffer is flushed to prevent errors
		case "$":
			obj, _ := f.ReadString(' ')
			w.WriteString("ID[STRING]: ")
			w.WriteString(obj)
			w.WriteString("\n")
			w.Flush()

		//For letters, the quotation mark is used to identify the string
		//Then a rune is read to keep the iterator at a " " to easily read the next token
		//The string is trimmed to have just the letters, which is then checked for uppercase
		case "\"":
			obj, _ := f.ReadString('"')
			f.ReadRune()
			obj = strings.TrimSuffix(obj, "\"")

			//When iterating through the string, if an upper case letter is found then the
			//val is incremented
			val := 0
			for i, _ := range obj {
				if unicode.IsUpper(int32(obj[i])) {
					val++
					break
				}
			}

			//If val was incremented then we know it doesn't match with grammar so the error
			//is thrown
			if val == 0 {
				w.WriteString("STRING: ")
				w.WriteString(obj)
				w.WriteString("\n")
				w.Flush()
			} else {
				w.WriteString("Lexical Error, unrecognized symbol ")
				w.WriteString("\"" + obj + "\"")
				w.WriteString("\n")
				w.Flush()
				tokens--

			}
		case "#":
			obj, _ := f.ReadString(' ')
			w.WriteString("ID[INT]: ")
			w.WriteString(obj)
			w.WriteString("\n")
			w.Flush()
		case "%":
			obj, _ := f.ReadString(' ')
			w.WriteString("ID[REAL]: ")
			w.WriteString(obj)
			w.WriteString("\n")
			w.Flush()

		//For tokens that require a string to identify, the string is read and checked before
		//writing to file, else it throws an error
		case "<":
			obj, _ := f.ReadString(' ')
			if 0 == strings.Compare(obj, "= ") {
				w.WriteString("ASSIGN\n")
				w.Flush()
			} else {
				w.WriteString("Lexical Error, unrecognized symbol ")
				w.WriteString("\"" + obj + "\"")
				w.WriteString("\n")
				w.Flush()
				tokens--
			}
		case ":":
			_, _ = f.ReadString('\n')
			w.WriteString("COLON\n")
			w.Flush()

		case "W":
			obj, _ := f.ReadString(' ')
			if 0 == strings.Compare(obj, "RITE ") {
				w.WriteString("WRITE\n")
				w.Flush()
			} else {
				w.WriteString("Lexical Error, unrecognized symbol ")
				w.WriteString("\"" + obj + "\"")
				w.WriteString("\n")
				w.Flush()
				tokens--

			}
		case "(":
			_, _ = f.ReadString(' ')
			w.WriteString("LPAREN\n")
			w.Flush()
		case "+":
			_, _ = f.ReadString(' ')
			w.WriteString("PLUS\n")
			w.Flush()
		case "/":
			_, _ = f.ReadString(' ')
			w.WriteString("OVER\n")
			w.Flush()
		case "*":
			_, _ = f.ReadString(' ')
			w.WriteString("TIMES\n")
			w.Flush()
		case "^":
			_, _ = f.ReadString(' ')
			w.WriteString("POWER\n")
			w.Flush()
		case "-":
			_, _ = f.ReadString(' ')
			w.WriteString("MINUS\n")
			w.Flush()
		case ")":
			_, _ = f.ReadString(' ')
			w.WriteString("RPAREN\n")
			w.Flush()
		case "B":
			obj, _ := f.ReadString('\n')
			if 0 == strings.Compare(obj, "EGIN\n") {
				w.WriteString("BEGIN\n")
				w.Flush()
			} else {
				w.WriteString("Lexical Error, unrecognized symbol ")
				w.WriteString("\"" + obj + "\"")
				w.WriteString("\n")
				w.Flush()
				tokens--
			}
		case "E":
			obj, _ := f.ReadString(' ')
			if 0 == strings.Compare(obj, "ND ") {
				w.WriteString("END\n")
				w.Flush()
			} else {
				w.WriteString("Lexical Error, unrecognized symbol ")
				w.WriteString("\"" + obj + "\"")
				w.WriteString("\n")
				w.Flush()
				tokens--
			}
		case ".":
			_, _ = f.ReadString(' ')
			w.WriteString("POINT\n")
			w.Flush()
		//In case of tabs being used, the tokens are reduced
		case "	":
			tokens--

		default:
			//if no case is matched then the string could be a number
			//So the string is trimmed to have no " " character in the end
			obj, _ := f.ReadString(' ')
			obj = strings.TrimSuffix(obj, " ")

			//first the string is checked for a "." to indicate that it is a float
			//if so then the string is converted to a float
			//else if the string is converted to an int
			//if the string is not a number then the default error is thrown
			//and the tokens are reduced by one
			if strings.Contains(obj, ".") {
				_, err := strconv.ParseFloat(obj, 64)
				if err == nil {
					w.WriteString("REAL_CONST: ")
					w.WriteString(string(r) + obj)
					w.WriteString("\n")
					w.Flush()
				}
			} else if _, err := strconv.ParseInt(obj, 10, 32); err == nil {
				w.WriteString("INT_CONST: ")
				w.WriteString(string(r) + obj)
				w.WriteString("\n")
				w.Flush()
			} else {
				w.WriteString("Lexical Error, unrecognized symbol ")
				w.WriteString(string(r) + obj)
				w.WriteString("\n")
				w.Flush()
				tokens--
			}

		}

		//reading rune by rune after each token
		r, _, fErr = f.ReadRune()

		//incrementing token after it has been read and written to a file
		tokens++
	}

	//Printing out the number of tokens and notifying the user that the output file has been created
	fmt.Println(tokens, "Tokens produced")
	fmt.Println("Processing Output File:", writeFileName)

}
