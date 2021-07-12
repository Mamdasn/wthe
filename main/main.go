package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/Mamdasn/wthe"

)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./wthe <input> [<output>]")
		os.Exit(1)
	}
	imagename := os.Args[1]
	var outputimagename string
	if len(os.Args) >= 3 {
		outputimagename = os.Args[2]
	}

	m_out := wthe.Wthe(imagename)

	if len(outputimagename) == 0 {
		folders := strings.Split(imagename, "/")
		imagename = folders[len(folders)-1]
		outputimagename = "output/Enhanced-" + imagename
	}
	fmt.Println(outputimagename)
	wthe.SaveImageToFilePath(outputimagename, m_out)
	fmt.Println("Done.")
}
