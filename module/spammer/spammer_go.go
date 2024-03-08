package main

import (
	"C"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var channelid string
var contents string
var token_file string
var proxie_file string
var threads_str string
var allchannel string
var delay_str string
var allping string
var mentions_str string

var sleepDuration time.Duration

func main() {
	args := os.Args[1:]

	serverid := args[0]
	channelid = args[1]
	contents = args[2]
	token_file = args[3]
	proxie_file = args[4]
	threads_str = args[5]
	threads, err := strconv.Atoi(threads_str)
	allchannel = args[6]
	delay_str = args[7]
	delay, err := strconv.ParseFloat(delay_str, 64)

	// delayが整数かどうかをチェックし、整数の場合は秒単位に変換
	if delay == float64(int(delay)) {
		sleepDuration = time.Duration(int(delay)) * time.Second
	} else {
		sleepDuration = time.Duration(delay * float64(time.Second))
	}

	// sleepDurationの間スリープ
	//time.Sleep(sleepDuration)
	//delay, err := strconv.Atoi(delay_str)
	allping = args[8]
	mentions_str = args[9]
	mentions, err := strconv.Atoi(mentions_str)
	users := args[10:]

	contents_tmp := ""

	fmt.Println(args[7:])

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
					channels, err := getChannels(getRandomToken(token_file), serverid)
					randomchannel, err := chooseRandomChannel(channels)
					if err != nil {
						fmt.Println("Error:", err)
					}
					channelid = randomchannel
				}
				if len(users) > 0 && allping == "True" {
					// ランダムに数個取り出す
					randomIDs := getRandomIDs(args[8:], mentions)
					formattedIDs := make([]string, len(randomIDs))
					for i, id := range randomIDs {
						formattedIDs[i] = formatID(id)
					}
					fmt.Printf("Random IDs: %s\n", strings.Join(formattedIDs, " | "))

					contents_tmp = contents + " " + strings.Join(formattedIDs, " ")
				}
				sendRequest(fmt.Sprintf("https://discord.com/api/v9/channels/%s/messages", channelid), contents_tmp, token_file, proxie_file)
				time.Sleep(sleepDuration)
			}
		}()
	}

	wg.Wait()
}

func sendRequest(url string, contents string, token_file string, proxie_file string) {
	proxy := getRandomProxy(proxie_file)

	headers := map[string]string{
		"Accept":             "*/*",
		"Accept-Encoding":    "gzip, deflate, br",
		"Accept-Language":    "en-US",
		"Authorization":      getRandomToken(token_file),
		"Connection":         "keep-alive",
		"Content-Type":       "application/json",
		"Host":               "discord.com",
		"Origin":             "https://discord.com",
		"Pragma":             "no-cache",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-origin",
		"sec-ch-ua-platform": "Windows",
		"sec-ch-ua-mobile":   "?0",
		"TE":                 "Trailers",
	}

	postData := map[string]interface{}{
		"content": contents,
	}

	postJSON, err := json.Marshal(postData)
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postJSON))
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("Success:", channelid, resp.StatusCode, proxy)
	} else {
		fmt.Println("Failed:", channelid, resp.StatusCode, proxy)
	}
}

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

func getRandomProxy(filepath string) *url.URL {
	proxies := readProxiesFromFile(filepath)
	if len(proxies) == 0 {
		fmt.Println("No proxies found in proxies.txt")
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	return proxies[rand.Intn(len(proxies))]
}

func readProxiesFromFile(filename string) []*url.URL {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	lines := strings.Split(string(content), "\n")
	var proxyList []*url.URL

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			proxyURL, err := url.Parse("http://" + line)
			if err != nil {
				fmt.Println("Error parsing proxy URL:", err)
			} else {
				proxyList = append(proxyList, proxyURL)
			}
		}
	}

	return proxyList
}

func getChannels(token string, guildID string) ([]string, error) {
	var channels []string

	for {
		url := fmt.Sprintf("https://discord.com/api/v9/guilds/%s/channels", guildID)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("authorization", token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == 200 {
			var data []map[string]interface{}
			err := json.Unmarshal(body, &data)
			if err != nil {
				return nil, err
			}

			for _, channel := range data {
				channelType, ok := channel["type"].(float64)
				if !ok {
					return nil, fmt.Errorf("Failed to parse channel type")
				}

				if channelType == 0 || channelType == 2 {
					channelID, ok := channel["id"].(string)
					if !ok {
						return nil, fmt.Errorf("Failed to parse channel id")
					}

					if !contains(channels, channelID) {
						channels = append(channels, channelID)
					}
				}
			}

			return channels, nil
		} else {
			fmt.Println(token)
			fmt.Println(resp.StatusCode)
			return nil, fmt.Errorf("Request failed with status code: %d", resp.StatusCode)
		}
	}
}

func contains(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

func chooseRandomChannel(channels []string) (string, error) {
	if len(channels) == 0 {
		return "", fmt.Errorf("No channels available")
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(channels))
	return channels[index], nil
}

func getRandomIDs(inputIDs []string, count int) []string {
	rand.Seed(time.Now().UnixNano())
	length := len(inputIDs)

	if count >= length {
		return inputIDs
	}

	result := make([]string, count)
	perm := rand.Perm(length)

	for i := 0; i < count; i++ {
		result[i] = inputIDs[perm[i]]
	}

	return result
}

func formatID(id string) string {
	id = strings.ReplaceAll(id, "'", "")
	id = strings.ReplaceAll(id, "\"", "")
	id = strings.ReplaceAll(id, "[", "")
	id = strings.ReplaceAll(id, "]", "")
	id = strings.ReplaceAll(id, ",", "")
	return fmt.Sprintf("<@%s>", id)
}
