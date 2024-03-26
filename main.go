package main

import (
	"bufio"
	"fmt"
	"log"
	"midnight/module/spammer"
	get_info "midnight/utilities"
	"os"
	"sort"
	"strconv"

	"github.com/fatih/color"
)

func main() {
	logo := `
                            ___  ____     _       _       _     _    ______      _     _           
            ..&@.           |  \/  (_)   | |     (_)     | |   | |   | ___ \    (_)   | |          
          ..@@@&.           | .  . |_  __| |_ __  _  __ _| |__ | |_  | |_/ /__ _ _  __| | ___ _ __ 
          .&&&&&,..         | |\/| | |/ _` + "`" + ` | '_ \| |/ _` + "`" + ` | '_ \| __| |    // _` + "`" + ` | |/ _` + "`" + ` |/ _ \ '__|
          ..&&&&&&&#.       | |  | | | (_| | | | | | (_| | | | | |_  | |\ \ (_| | | (_| |  __/ |   
            ..#&&&...       \_|  |_/_|\__,_|_| |_|_|\__, |_| |_|\__| \_| \_\__,_|_|\__,_|\___|_|   
                                                     __/ |                                         
                                                    |___/                                      
`
	info := fmt.Sprintf(`  HWID: [%s]                Version: [1.0.0]`, get_info.GetHwid())

	print_logo := color.New(color.FgWhite).Add(color.FgBlue)
	print_logo.Println(logo)
	fmt.Println(info)
	fmt.Println("\n------------------------------------------------------------------------------------------------------------------------\n")

	opt := map[int]string{
		1: "Coming Soon",
		2: "Spammer",
		3: "Coming Soon",
		4: "Coming Soon",
		5: "Coming Soon",
		6: "Coming Soon",
		7: "Coming Soon",
	}
	PrintMenu(opt)

	mode := Input("\nMode >> ")

	switch mode {
	case "1", "01":
		fmt.Println("Mass DM")
	case "2", "02":
		serverid := Input("Server ID: ")
		channelid := Input("Channel ID: ")
		msg := Input("Message: ")
		tokens_file := Input("Tokens File: ")
		proxie_file := Input("Proxie File: ")
		allchannel := InputBool("Allchannel")
		threads := InputInt("Threads")
		delay := InputInt("Delay")
		//allping := getInput("Allping: ")
		allping := InputBool("Allping")
		mentions_str := 0
		if allping {
			//mentions_str := getInput("Mentions: ")
			mentions_str = InputInt("How many Mentions?")
			//fmt.Println("Dm Spam")
		}
		spammer.Start(serverid, channelid, msg, tokens_file, proxie_file, threads, allchannel, delay, allping, mentions_str)
	case "3", "03":
		fmt.Println("React Verify")
	case "4", "04":
		fmt.Println("Joiner")
	case "5", "05":
		fmt.Println("Leaver")
	case "6", "06":
		fmt.Println("Accept Rules")
	case "7", "07":
		fmt.Println("Raid Channel")
	default:
		fmt.Println("Invalid mode")
	}
}

func Input(text string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(text)
	scanner.Scan()
	return scanner.Text()
}

func InputBool(text string) bool {
	if Input(text+" y/n: ") == "y" {
		return true
	}
	return false
}

// InputInt Takes user input from the Command line
// and returns it as an int
func InputInt(text string) int {
	var d int

	fmt.Print(text + ": ")
	_, err := fmt.Scanln(&d)
	if err != nil {
		log.Println(err)
	}
	return d
}
func PrintMenu(options map[int]string) {
	var maxLen int
	for _, value := range options {
		if len(value) > maxLen {
			maxLen = len(value)
		}
	}

	keys := make([]int, 0, len(options))
	for k := range options {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	numCols := 6
	numRows := (len(keys) + numCols - 1) / numCols

	for col := 0; col < numCols; col++ {
		for row := 0; row < numRows; row++ {
			index := col + row*numCols
			if index < len(keys) {
				optnum := strconv.Itoa(keys[index])
				if len(optnum) == 1 {
					optnum = "0" + optnum
				}
				fmt.Printf("  [%s]  %-20s", optnum, options[keys[index]])
			}
		}
		fmt.Println()
	}
}
