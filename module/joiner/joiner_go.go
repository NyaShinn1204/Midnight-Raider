package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
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

var sleepDuration time.Duration

func main() {
	args := os.Args[1:]
	token_file = args[0]
	serverid = args[1]
	invitelink = args[2]
	memberscreen = args[3]
	delay_str = args[4]
	delay, err := strconv.ParseFloat(delay_str, 64)
	if err != nil {
		// エラー処理
	}

	// delayが整数かどうかをチェックし、整数の場合は秒単位に変換
	if delay == float64(int(delay)) {
		sleepDuration = time.Duration(int(delay)) * time.Second
	} else {
		sleepDuration = time.Duration(delay * float64(time.Second))
	}

	// sleepDurationの間スリープ
	//time.Sleep(sleepDuration)
	//delay, err := strconv.Atoi(delay_str)
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

	start(tokens, serverid, invitelink, memberscreen, answers, apikey, bypasscaptcha, deletejoinmsg, joinchannelid)
}

func start(tokens []string, serverID, inviteLink string, memberScreen string, answers string, apis string, bypassCaptcha string, deleteJoinMs string, joinChannelID string) {
	for _, token := range tokens {
		go joinerThread(token, serverID, inviteLink, memberScreen, answers, apis, bypassCaptcha, deleteJoinMs, joinChannelID)
		time.Sleep(sleepDuration)
	}
}

// Headers関係はここに置くかも しらんけど

func randomAgent() string {
	agents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15_7; rv:92.0) Gecko/20100101 Firefox/92.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36 Edg/93.0.961.52",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/37.0.2062.94 Chrome/37.0.2062.94 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/600.8.9 (KHTML, like Gecko) Version/8.0.8 Safari/600.8.9",
	}
	rand.Seed(time.Now().UnixNano())
	return agents[rand.Intn(len(agents))]
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

func getFingerprint() (string, error) {
	// リクエスト用のヘッダーを定義
	headers := map[string]string{
		"Accept":          "*/*",
		"Accept-Language": "en-US,en;q=0.9",
		"Connection":      "keep-alive",
		"Referer":         "https://discord.com/",
		"Sec-Fetch-Dest":  "empty",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Site":  "same-origin",
		"Sec-GPC":         "1",
		"User-Agent":      "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36 Edg/114.0.1823.51",
		"X-Track":         "eyJvcyI6IklPUyIsImJyb3dzZXIiOiJTYWZlIiwic3lzdGVtX2xvY2FsZSI6ImVuLUdCIiwiYnJvd3Nlcl91c2VyX2FnZW50IjoiTW96aWxsYS81LjAgKElQaG9uZTsgQ1BVIEludGVybmFsIFByb2R1Y3RzIFN0b3JlLCBhcHBsaWNhdGlvbi8yMDUuMS4xNSAoS0hUTUwpIFZlcnNpb24vMTUuMCBNb2JpbGUvMTVFMjQ4IFNhZmFyaS82MDQuMSIsImJyb3dzZXJfdmVyc2lvbiI6IjE1LjAiLCJvc192IjoiIiwicmVmZXJyZXIiOiIiLCJyZWZlcnJpbmdfZG9tYWluIjoiIiwicmVmZXJyZXJfZG9tYWluX2Nvb2tpZSI6InN0YWJsZSIsImNsaWVudF9idWlsZF9udW1iZXIiOjk5OTksImNsaWVudF9ldmVudF9zb3VyY2UiOiJzdGFibGUiLCJjbGllbnRfZXZlbnRfc291cmNlIjoic3RhYmxlIn0",
	}

	// リクエストの作成
	req, err := http.NewRequest("GET", "https://discord.com/api/v9/experiments", nil)
	if err != nil {
		return "", err
	}

	// ヘッダーを追加
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// HTTPクライアントを作成
	client := &http.Client{}

	// リクエストを送信し、レスポンスを取得
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// レスポンスの内容を読み取る
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// レスポンスの内容をJSONとして解析
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	// "fingerprint"フィールドの値を取得
	fingerprint, ok := data["fingerprint"].(string)
	if !ok {
		return "", errors.New("fingerprint not found in response")
	}

	// 取得したFingerprintを返す
	return fingerprint, nil
}

func getCookie() string {
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://discord.com/api/v9/experiments", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return ""
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()
	cookieStrings := make([]string, len(cookies))
	for i, cookie := range cookies {
		cookieStrings[i] = fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
	}
	return strings.Join(cookieStrings, "; ")
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

func requestHeader(token string, includeFingerprint, includeCookie bool) map[string]string {
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

	// fingerprintを含める場合
	if includeFingerprint {
		fingerprint, err := getFingerprint()
		if err != nil {
			fmt.Println("Failed to get fingerprint:", err)
		}
		headers["X-Fingerprint"] = fingerprint
	}

	// cookieを含める場合
	if includeCookie {
		headers["Cookie"] = getCookie()
	}

	return headers
}

func joinerThread(token, serverID, inviteLink string, memberScreen string, answers string, apis string, bypassCaptcha string, deleteJoinMs string, joinChannelID string) {
	// 必要な処理を実装
	fmt.Println(token)
	fmt.Println(serverID)
	fmt.Println(inviteLink)
	fmt.Println(memberScreen)
	fmt.Println(answers)
	fmt.Println(apis)
	fmt.Println(bypassCaptcha)
	fmt.Println(deleteJoinMs)
	fmt.Println(joinChannelID)
	fmt.Println(requestHeader("Token", false, false))
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
