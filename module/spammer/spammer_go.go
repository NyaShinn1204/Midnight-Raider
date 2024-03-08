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
import (
	"encoding/base64"
	"regexp"
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

func requestHeader(token string) map[string]string {
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
		"Authorization":      token,
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

	//tlsConfig := &tls.Config{
	//	MinVersion:               tls.VersionTLS12,                            // 最低限のTLSバージョン
	//	CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384}, // 楕円曲線の選択
	//	PreferServerCipherSuites: true,                                        // サーバーが使用する暗号スイートを優先する
	//	CipherSuites: []uint16{
	//		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384, // 暗号スイートの指定
	//		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	//	},
	//	// 証明書の検証関連の設定
	//	//RootCAs:            certPool, // ルート証明書を検証するCAリスト
	//	//InsecureSkipVerify: false,    // サーバー証明書の検証をスキップするかどうか
	//	//// その他の設定
	//	//ClientAuth: tls.NoClientCert, // クライアント証明書の要求
	//	//ServerName: "example.com",    // サーバー名
	//}
	//proxyURL, err := url.Parse("http://tbkzktta:de8si82ghq2y@154.95.36.199:6893")
	//if err != nil {
	//	fmt.Println(err)
	//}
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

	// HTTPリクエスト作成
	//req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v9/invites/%s", inviteLink), nil)
	//if err != nil {
	//	log.Fatalf("Failed to create request: %v", err)
	//}

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

	// リクエストヘッダー設定
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// リクエスト送信
	requests, err := session.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v", err)
	}
	//defer requests.Body.Close()

	defer requests.Body.Close()

	// レスポンスボディをバイト配列に読み込む
	body, err := ioutil.ReadAll(requests.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v", err)
	}

	//fmt.Println(requests.Body)

	// レスポンスボディをJSONとしてパース
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		fmt.Printf("Failed to parse response body: %v", err)
	}

	fmt.Println(requests.StatusCode)

	//reqHeader := requestHeader(getRandomToken(token_file), true, true)
	//headers := reqHeader
	//
	//payloaddata := map[string]interface{}{
	//	"content": contents,
	//}
	//
	//payload, err := json.Marshal(payloaddata)
	//if err != nil {
	//	fmt.Println("JSON marshal error:", err)
	//	return
	//}
	//
	//client := &http.Client{
	//	Transport: &http.Transport{
	//		Proxy: http.ProxyURL(proxy),
	//	},
	//}
	//
	//req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	//if err != nil {
	//	fmt.Println("Request error:", err)
	//	return
	//}
	//
	//for key, value := range headers {
	//	req.Header.Set(key, value)
	//}
	//
	//resp, err := client.Do(req)
	//if err != nil {
	//	return
	//}
	//defer resp.Body.Close()
	//
	if requests.StatusCode == 200 {
		fmt.Println("Success:", channelid, requests.StatusCode, proxy)
	} else {
		fmt.Println("Failed:", channelid, requests.StatusCode, proxy)
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
