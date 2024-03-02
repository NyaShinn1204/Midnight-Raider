package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var token_file string
var tokens []string
var serverid string
var invitelink string
var memberscreen string
var delay_str string
var bypasscaptcha string
var answers string
var apikey string
var deletejoinmsg string
var joinchannelid string

func main() {
	args := os.Args[1:]
	token_file = args[0]
	serverid = args[1]
	invitelink = args[2]
	memberscreen = args[3]
	delay_str = args[4]
	delay, err := strconv.Atoi(delay_str)
	bypasscaptcha = args[5]
	answers = args[6]
	apikey = args[7]
	deletejoinmsg = args[8]
	joinchannelid = args[9]

	if err != nil {
		fmt.Println("変換に失敗しました:", err)
		return
	}

	// トークンを格納するスライス
	// トークンファイルからトークンを読み込む
	if token_file != "" {
		file, err := os.Open(token_file)
		if err != nil {
			fmt.Println("Error opening token file:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			tokens = append(tokens, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading token file:", err)
			return
		}
	} else {
		fmt.Println("Token file path is required.")
		return
	}

	start(tokens, serverid, invitelink, memberscreen, delay, answers, apikey, bypasscaptcha, deletejoinmsg, joinchannelid)
}

func start(tokens []string, serverID, inviteLink string, memberScreen string, delay int, answers string, apis string, bypassCaptcha string, deleteJoinMs string, joinChannelID string) {
	for _, token := range tokens {
		go joinerThread(token, serverID, inviteLink, memberScreen, answers, apis, bypassCaptcha, deleteJoinMs, joinChannelID)
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

// Headers関係はここに置くかも しらんけど

func joinerThread(token, serverID, inviteLink string, memberScreen string, answers string, apis string, bypassCaptcha string, deleteJoinMs string, joinChannelID string) {
	// 必要な処理を実装
	fmt.Printf(token)
	fmt.Printf(serverID)
	fmt.Printf(inviteLink)
	fmt.Printf(memberScreen)
	fmt.Printf(answers)
	fmt.Printf(apis)
	fmt.Printf(bypassCaptcha)
	fmt.Printf(deleteJoinMs)
	fmt.Print(joinChannelID)
}

func deleteJoinMsg(token, joinChannelID string) {
	// 処理を記述
}

// 他の関数やモジュールの実装

func getRandomToken(filepath string) string {
	tokens := readTokensFromFile(filepath)
	if len(tokens) == 0 {
		fmt.Println("No tokens found in token.txt")
		return ""
	}

	rand.Seed(time.Now().UnixNano())
	return tokens[rand.Intn(len(tokens))]
}

func readTokensFromFile(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	return strings.Fields(string(content))
}
