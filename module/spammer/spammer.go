package spammer

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

var sleepDuration time.Duration

func Start(serverid string, channelid string, contents string, tokens_file string, proxie_file string, threads_str string, allchannel string, delay_str string, allping string, mentions_str string) {

	threads, err := strconv.Atoi(threads_str)
	delay, err := strconv.Atoi(delay_str)
	sleepDuration = time.Duration(delay) * time.Second

	mentions, err := strconv.Atoi(mentions_str)
	//token := getRandomToken(token_file)

	//members := getMembers(token, serverid, channelid)
	//userIDs := make([]string, len(members))
	//for i, member := range members {
	//	userIDs[i] = member
	//}

	contents_tmp := ""

	if err != nil {
		fmt.Println("変換に失敗しました:", err)
		return
	}

	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				contents_tmp = contents
				if allchannel == "True" {
					//channels, err := getChannels(getRandomToken(tokens_file), serverid)
					//randomchannel, err := chooseRandomChannel(channels)
					//if err != nil {
					//	fmt.Println("Error:", err)
					//}
					//channelid = randomchannel
					channelid = "LoL"
				}
				if allping == "True" {
					rand.Seed(time.Now().UnixNano())

					numElements := mentions

					selectedElements := make([]string, numElements)
					//if len(userIDs) > 0 {
					//	for i := 0; i < numElements; i++ {
					//		randomIndex := rand.Intn(len(userIDs))
					//		selectedElements[i] = "<@" + userIDs[randomIndex] + ">"
					//	}
					//} else {
					//	fmt.Println("元の配列に要素がありません")
					//}

					convert_mentions := strings.Join(selectedElements, ", ")

					contents_tmp = contents + " " + convert_mentions
				}
				sendRequest(fmt.Sprintf("https://discord.com/api/v9/channels/%s/messages", channelid), contents_tmp, tokens_file, proxie_file)
				fmt.Println()
				time.Sleep(sleepDuration)
			}
		}()
	}

	wg.Wait()
}

func sendRequest(url string, content string, tokens_file string, proxie_file string) {
	fmt.Println("Testing Threads.")
}
