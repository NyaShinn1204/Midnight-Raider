package main

import (
	"C"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
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
	delay, err := strconv.Atoi(delay_str)
	sleepDuration = time.Duration(delay) * time.Second

	// delayから変換するやつバグってるので一時的に殺します
	//delay, err := strconv.ParseFloat(delay_str, 64)

	//// delayが整数かどうかをチェックし、整数の場合は秒単位に変換
	//if delay == float64(int(delay)) {
	//	sleepDuration = time.Duration(int(delay)) * time.Second
	//} else {
	//	sleepDuration = time.Duration(delay * float64(time.Second))
	//}

	// sleepDurationの間スリープ
	//time.Sleep(sleepDuration)
	//delay, err := strconv.Atoi(delay_str)

	allping = args[8]
	mentions_str = args[9]
	mentions, err := strconv.Atoi(mentions_str)
	//	users := args[10:]
	token := getRandomToken(token_file)

	members := getMembers(token, serverid, channelid)
	userIDs := make([]string, len(members))
	for i, member := range members {
		userIDs[i] = member
	}

	contents_tmp := ""

	//fmt.Println(args[7:])

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
				if allping == "True" {
					//// ランダムに数個取り出す
					//randomIDs := getRandomIDs(args[10:], mentions)
					//formattedIDs := make([]string, len(randomIDs))
					//for i, id := range randomIDs {
					//	formattedIDs[i] = formatID(id)
					//}
					//fmt.Printf("Random IDs: %s\n", strings.Join(formattedIDs, " | "))
					// ランダムシードの初期化
					rand.Seed(time.Now().UnixNano())

					// 取り出す要素の数
					numElements := mentions

					selectedElements := make([]string, numElements)
					if len(userIDs) > 0 {
						for i := 0; i < numElements; i++ {
							randomIndex := rand.Intn(len(userIDs))
							selectedElements[i] = "<@" + userIDs[randomIndex] + ">"
						}
					} else {
						fmt.Println("元の配列に要素がありません")
					}

					convert_mentions := strings.Join(selectedElements, ", ")

					contents_tmp = contents + " " + convert_mentions
				}
				sendRequest(fmt.Sprintf("https://discord.com/api/v9/channels/%s/messages", channelid), contents_tmp, token_file, proxie_file)
				fmt.Println()
				time.Sleep(sleepDuration)
			}
		}()
	}

	wg.Wait()
}

func generateSuperProperties() string {
	buildNum, err := getBuildnum()
	if err != nil {
		fmt.Println("Error:", err)
	}
	agentString := randomAgent()
	browserData := strings.Split(agentString, " ")[len(strings.Split(agentString, " "))-1]
	var agentOS string
	if strings.Contains(agentString, "Windows") {
		agentOS = "Windows"
	} else if strings.Contains(agentString, "Macintosh") {
		agentOS = "Macintosh"
	}
	var osVersion string
	if agentOS == "Macintosh" {
		osVersion = fmt.Sprintf("Intel Mac OS X 10_15_%d", rand.Intn(3)+5)
	} else {
		osVersion = "10"
	}
	deviceInfo := map[string]interface{}{
		"os":                       agentOS,
		"browser":                  strings.Split(browserData, "/")[0],
		"device":                   "",
		"system_locale":            "ja-JP",
		"browser_user_agent":       agentString,
		"browser_version":          strings.Split(browserData, "/")[1],
		"os_version":               osVersion,
		"referrer":                 "",
		"referring_domain":         "",
		"referrer_current":         "",
		"referring_domain_current": "",
		"release_channel":          "stable",
		"client_build_number":      buildNum,
		"client_event_source":      nil,
	}
	jsonData, _ := json.Marshal(deviceInfo)
	return base64.StdEncoding.EncodeToString(jsonData)
}

func randomAgent() string {
	// ファイルを読み込み
	content, err := ioutil.ReadFile("../../data/user-agent.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return "" // エラーが発生した場合、空の文字列を返す
	}

	// ファイルの内容を改行で分割し、agentsスライスに追加
	lines := strings.Split(string(content), "\n")

	// キャリッジリターンをトリム
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}

	rand.Seed(time.Now().UnixNano())
	return lines[rand.Intn(len(lines))]
}

func getBuildnum() (int, error) {
	// Discordのログインページからテキストを取得
	resp, err := http.Get("https://discord.com/login")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	text := string(body)

	// スクリプトのURLを抽出
	re := regexp.MustCompile(`\d+\.\w+\.js|sentry\.\w+\.js`)
	matches := re.FindAllString(text, -1)
	if len(matches) == 0 {
		return 0, fmt.Errorf("script URL not found")
	}
	scriptURL := "https://discord.com/assets/" + matches[len(matches)-1]

	// スクリプトのテキストを取得
	resp, err = http.Get(scriptURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	text = string(body)

	// buildNumberを抽出
	index := 0
	for {
		index = regexp.MustCompile("buildNumber").FindStringIndex(text)[0]
		if index != 0 {
			break
		}
		text = text[index+1:]
	}
	index += 26
	buildNum := 0
	for _, c := range text[index : index+6] {
		if c >= '0' && c <= '9' {
			buildNum = buildNum*10 + int(c-'0')
		} else {
			break
		}
	}

	return buildNum, nil
}

func requestHeader(get_token string) map[string]string {
	// ランダムなユーザーエージェントを生成
	agentString := randomAgent()

	// システムのビルド番号を取得
	buildNum, err := getBuildnum()
	if err != nil {
		fmt.Println("Failed Get BuildNum:", err)
	}

	// デバイス情報の作成
	deviceInfo := map[string]interface{}{
		"os":                       "Windows",
		"browser":                  "Chrome",
		"device":                   "",
		"system_locale":            "ja-JP",
		"browser_user_agent":       agentString,
		"browser_version":          "95.0.4638.54",
		"os_version":               "10",
		"referrer":                 "",
		"referring_domain":         "",
		"referrer_current":         "",
		"referring_domain_current": "",
		"release_channel":          "stable",
		"client_build_number":      buildNum,
		"client_event_source":      nil,
	}

	// デバイス情報をBase64エンコード
	deviceInfoJSON, _ := json.Marshal(deviceInfo)
	deviceInfoBase64 := base64.StdEncoding.EncodeToString(deviceInfoJSON)

	// リクエストヘッダーの作成
	headers := map[string]string{
		"Accept":             "*/*",
		"Accept-Encoding":    "gzip, deflate, br",
		"Accept-Language":    "en-US",
		"Authorization":      get_token,
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
		"User-Agent":         agentString,
		"X-Debug-Options":    "bugReporterEnabled",
		"X-Discord-Locale":   "ja",
		"X-Discord-Timezone": "Asia/Tokyo",
		"X-Super-Properties": deviceInfoBase64,
	}

	return headers
}

func getSession(useproxies bool, proxyurl *url.URL) *http.Client {
	var transport *http.Transport
	if useproxies {
		transport = &http.Transport{
			//TLSClientConfig: tlsConfig,
			Proxy: http.ProxyURL(proxyurl),
		}
	} else {
		transport = &http.Transport{
			//	TLSClientConfig: tlsConfig,
		}
	}

	return &http.Client{
		Transport: transport,
	}
}

func sendRequest(url string, contents string, token_file string, proxie_file string) {
	proxy := getRandomProxy(proxie_file)
	session := getSession(true, proxy)

	reqHeader := requestHeader(getRandomToken(token_file))
	headers := reqHeader

	payload, err := json.Marshal(map[string]interface{}{"content": contents})
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}

	req.Close = true

	// リクエストヘッダー設定
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// リクエスト送信
	requests, err := session.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v", err)
		return
		//return
	}

	defer requests.Body.Close()

	//fmt.Println(requests.StatusCode)

	if requests.StatusCode == 200 {
		fmt.Println("[+] Succes to Sent ChannelID:", channelid, requests.StatusCode, proxy)
	} else {
		fmt.Println("[-] Failed to Sent ChannelID:", channelid, requests.StatusCode, proxy)
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

func getChannels(get_token string, guildID string) ([]string, error) {
	var channels []string

	for {
		url := fmt.Sprintf("https://discord.com/api/v9/guilds/%s/channels", guildID)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		reqHeader := requestHeader(get_token)
		headers := reqHeader
		// リクエストヘッダー設定
		for key, value := range headers {
			req.Header.Set(key, value)
		}
		//req.Header.Set("authorization", token)

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
			fmt.Println(get_token)
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

func getMembers(get_token, server, channel string) []string {
	conn, _, err := websocket.DefaultDialer.Dial("wss://gateway.discord.gg/?v=10&encoding=json", nil)
	if err != nil {
		fmt.Println("dial:", err)
	}
	defer conn.Close()

	users := []string{}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			if err.Error() == "websocket: close 4000: Unknown error." {
				fmt.Println("トークンがサーバーに入っていません")
				os.Exit(1)
			}
			return nil
		}

		var response map[string]interface{}
		if err := json.Unmarshal(message, &response); err != nil {
			fmt.Println("unmarshal:", err)
			return nil
		}

		if response["t"] == nil {
			sendData := map[string]interface{}{
				"op": 2,
				"d": map[string]interface{}{
					"token":        get_token,
					"capabilities": 16381,
					"properties": map[string]interface{}{
						"os":                       "Android",
						"browser":                  "Discord Android",
						"device":                   "Android",
						"system_locale":            "ja-JP",
						"browser_user_agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
						"browser_version":          "122.0.0.0",
						"os_version":               "",
						"referrer":                 "",
						"referring_domain":         "",
						"referrer_current":         "",
						"referring_domain_current": "",
						"release_channel":          "stable",
						"client_build_number":      263582,
						"client_event_source":      nil,
					},
					"presence": map[string]interface{}{
						"status":     "invisible",
						"since":      0,
						"activities": []interface{}{},
						"afk":        false,
					},
					"compress": false,
					"client_state": map[string]interface{}{
						"guild_versions":              map[string]interface{}{},
						"highest_last_message_id":     "0",
						"read_state_version":          0,
						"user_guild_settings_version": -1,
						"private_channels_version":    "0",
						"api_code_version":            0,
					},
				},
			}
			if err := conn.WriteJSON(sendData); err != nil {
				fmt.Println("write:", err)
				return nil
			}
		} else if response["t"].(string) == "READY_SUPPLEMENTAL" {
			sendData := map[string]interface{}{
				"op": 14,
				"d": map[string]interface{}{
					"guild_id":   server,
					"typing":     true,
					"activities": true,
					"threads":    true,
					"channels": map[string]interface{}{
						channel: [][]int{{0, 99}, {100, 199}, {200, 299}},
					},
				},
			}
			if err := conn.WriteJSON(sendData); err != nil {
				fmt.Println("write:", err)
				return nil
			}
		} else if response["t"].(string) == "GUILD_MEMBER_LIST_UPDATE" {
			items := response["d"].(map[string]interface{})["ops"].([]interface{})
			for _, item := range items {
				for _, member := range item.(map[string]interface{})["items"].([]interface{}) {
					if member.(map[string]interface{})["member"] != nil {
						users = append(users, member.(map[string]interface{})["member"].(map[string]interface{})["user"].(map[string]interface{})["id"].(string))
					}
				}
			}
			return users
		}
	}
}
