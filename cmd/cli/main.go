package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("db> ")
		input, err := reader.ReadString(';')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)

		// 用户输入exit时退出
		if input == "exit" {
			fmt.Println("Exiting the program.")
			break
		}

		// 在这里处理其他命令
		fmt.Printf("You entered: %s\n", input)
	}
}
