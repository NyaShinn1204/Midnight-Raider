package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
	// token_file serverid invitelink memberscreen delay bypasscaptcha answers apikey deletejoinmsg joinchannelid
	args := os.Args[1:]
	token_file = args[0]
	serverid = args[1]
	invitelink = args[2]
	memberscreen = args[3]
	delay_str = args[4]
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

// ユーティリティー関係ここに置くかも 知らんけど
func extract(formatToken string) string {
	if match, _ := regexp.MatchString(`(.+):`, formatToken); match {
		return strings.Split(formatToken, ":")[1]
	}
	return formatToken
}

// Headers関係はここに置くかも しらんけど

func randomAgent() string {
	// ファイルを読み込み
	content, err := ioutil.ReadFile("user-agent.txt")
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

func getSession() *http.Client {
	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,                            // 最低限のTLSバージョン
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384}, // 楕円曲線の選択
		PreferServerCipherSuites: true,                                        // サーバーが使用する暗号スイートを優先する
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384, // 暗号スイートの指定
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		// 証明書の検証関連の設定
		//RootCAs:            certPool, // ルート証明書を検証するCAリスト
		//InsecureSkipVerify: false,    // サーバー証明書の検証をスキップするかどうか
		//// その他の設定
		//ClientAuth: tls.NoClientCert, // クライアント証明書の要求
		//ServerName: "example.com",    // サーバー名
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return &http.Client{
		Transport: transport,
	}
}

func solver(answers string, token string, url string, sitekey string, apikey string) string {
	fmt.Println(answers)
	fmt.Println(token)
	fmt.Println(url)
	fmt.Println(sitekey)
	fmt.Println(apikey)
	return "aiueo1425"
}

//solver(answers, token, "https://discord.com", joinreq.JSON().(map[string]interface{})["captcha_sitekey"].(string), apis

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
	if serverid != "None" {
		// HTTP GETリクエストを送信してレスポンスを取得
		resp, err := http.Get(fmt.Sprintf("https://discord.com/api/v9/invites/%s?with_counts=true&with_expiration=true", inviteLink))
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()

		// レスポンスのステータスコードを確認
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Server returned non-OK status: %d", resp.StatusCode)
		}

		// JSONデコードしてマップに変換
		var data map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&data)
		//if err != nil {
		//	return nil, err
		//}
		// HTTPリクエストを送信してレスポンスを取得
		//data, err := sendRequest("your_invite_code_here")
		if err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}

		// "guild_id"キーが存在するか確認し、存在する場合はその値を取得
		serverID, ok := data["guild_id"].(string)
		if !ok {
			log.Fatal("Server ID not found or not a string")
		}

		fmt.Println("Server ID:", serverID)
	}
	// JSON形式の文字列に変換
	// お試しjson show
	//jsonData, err := json.MarshalIndent(requestHeader(token, false, false), "", "    ")
	//if err != nil {
	//	fmt.Println("Error marshalling JSON:", err)
	//	return
	//}
	//fmt.Println(string(jsonData))

	//fmt.Println(requestHeader("Token", false, false))

	extractToken := fmt.Sprintf("%s.%s", strings.Split(extract(token+"]"), ".")[0], strings.Split(extract(token+"]"), ".")[1])

	session := getSession()
	reqHeader := requestHeader(token, true, true)
	headers := reqHeader

	//client := session

	//joinreq, err := client.Post(fmt.Sprintf("https://discord.com/api/v9/invites/%s", inviteLink), "", nil)
	//if err != nil {
	//		log.Fatalf("Failed to send join request: %v", err)
	//}

	// HTTPリクエスト作成
	req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v9/invites/%s", inviteLink), nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// リクエストヘッダー設定
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// リクエスト送信
	joinreq, err := session.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	//defer joinreq.Body.Close()

	defer joinreq.Body.Close()

	// レスポンスボディをバイト配列に読み込む
	body, err := ioutil.ReadAll(joinreq.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// レスポンスボディをJSONとしてパース
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		log.Fatalf("Failed to parse response body: %v", err)
	}

	fmt.Println(joinreq.StatusCode)
	if joinreq.StatusCode == 400 {
		if bypassCaptcha == "True" {
			fmt.Printf("Solving Captcha | %s\n", extractToken)
			// jsonResponse["captcha_sitekey"]が文字列型であることを確認し、型アサーションを行う
			captchaSiteKey, ok := jsonResponse["captcha_sitekey"].(string)
			if !ok {
				log.Fatalf("Failed to assert captcha_sitekey as string")
			}
			payload := map[string]interface{}{
				"captcha_key": solver(answers, token, "https://discord.com", captchaSiteKey, apis),
			}
			encode_payload, err := json.Marshal(payload)
			if err != nil {
				log.Fatalf("Failed to encode payload: %v", err)
			}
			// HTTPリクエスト作成
			newreq, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v9/invites/%s", inviteLink), bytes.NewBuffer(encode_payload))
			if err != nil {
				log.Fatalf("Failed to create request: %v", err)
			}

			// リクエストヘッダー設定
			for key, value := range headers {
				newreq.Header.Set(key, value)
			}

			// リクエスト送信
			newresponse, err := session.Do(newreq)
			if err != nil {
				log.Fatalf("Failed to send request: %v", err)
			}
			fmt.Println(newresponse.StatusCode)
			defer newresponse.Body.Close()

			// レスポンスボディをバイト配列に読み込む
			body, err := ioutil.ReadAll(newresponse.Body)
			if err != nil {
				log.Fatalf("Failed to read response body: %v", err)
			}

			// レスポンスボディをJSONとしてパース
			var jsonResponse map[string]interface{}
			if err := json.Unmarshal(body, &jsonResponse); err != nil {
				log.Fatalf("Failed to parse response body: %v", err)
			}
			//defer joinreq.Body.Close()
			//joinreq, err = client.R().
			//	SetHeaders(headers).
			//	SetJSON(payload).
			//	Post(fmt.Sprintf("https://discord.com/api/v9/invites/%s", inviteLink))
			//if err != nil {
			//	log.Fatalf("Failed to send join request: %v", err)
			//}
		} else {
			payload := map[string]interface{}{
				"captcha_key": "",
			}
			encode_payload, err := json.Marshal(payload)
			if err != nil {
				log.Fatalf("Failed to encode payload: %v", err)
			}
			// HTTPリクエスト作成
			newreq, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v9/invites/%s", inviteLink), bytes.NewBuffer(encode_payload))
			if err != nil {
				log.Fatalf("Failed to create request: %v", err)
			}

			// リクエストヘッダー設定
			for key, value := range headers {
				newreq.Header.Set(key, value)
			}

			// リクエスト送信
			newresponse, err := session.Do(newreq)
			if err != nil {
				log.Fatalf("Failed to send request: %v", err)
			}
			fmt.Println(newresponse.StatusCode)
			defer newresponse.Body.Close()

			// レスポンスボディをバイト配列に読み込む
			body, err := ioutil.ReadAll(newresponse.Body)
			if err != nil {
				log.Fatalf("Failed to read response body: %v", err)
			}

			// レスポンスボディをJSONとしてパース
			var jsonResponse map[string]interface{}
			if err := json.Unmarshal(body, &jsonResponse); err != nil {
				log.Fatalf("Failed to parse response body: %v", err)
			}
			//defer joinreq.Body.Close()
			//joinreq, err = client.R().
			//	SetHeaders(headers).
			//	SetJSON(payload).
			//	Post(fmt.Sprintf("https://discord.com/api/v9/invites/%s", inviteLink))
			if err != nil {
				log.Fatalf("Failed to send join request: %v", err)
			}
		}
		if joinreq.StatusCode == 200 {
			if _, ok := jsonResponse["captcha_key"]; !ok {
				if _, ok := jsonResponse["You need to verify your account in order to perform this action."]; ok {
					fmt.Printf("認証が必要です | %s\n", extractToken)
					fmt.Println("失敗しました")
				}
				fmt.Printf("Successfully Token Join | %s\n", extractToken)
				if deleteJoinMs == "True" {
					fmt.Printf("Deleting Join Message | %s\n", extractToken)
					deleteJoinMsg(token, joinChannelID)
				}
				fmt.Println("成功しました")
			}
			//if memberScreen {
			//	acceptRulesBypass(token, joinreq.JSON(), serverID, inviteLink)
			//}
			//if changeNick {
			//	changeNicker(token, serverID, nickname)
			//}
		} else {
			if _, ok := jsonResponse["captcha_key"]; ok {
				fmt.Printf("Failed Token Join (Captcha Wrong) | %s\n", extractToken)
				fmt.Println(jsonResponse)
				fmt.Println("失敗しました")
			} else {
				fmt.Printf("Failed Captcha Bypass | %s | %s\n", extractToken, strings.ReplaceAll(string(body), "\n", ""))
			}
		}
	} else if joinreq.StatusCode == 200 {
		if _, ok := jsonResponse["captcha_key"]; !ok {
			if _, ok := jsonResponse["You need to verify your account in order to perform this action."]; ok {
				fmt.Printf("認証が必要です | %s\n", extractToken)
				fmt.Println("失敗しました")
			}
			fmt.Printf("Successfully Token Join | %s\n", extractToken)
			if deleteJoinMs == "True" {
				fmt.Printf("Deleting Join Message | %s\n", extractToken)
				deleteJoinMsg(token, joinChannelID)
			}
			fmt.Println("成功しました")
		}
		//if memberScreen {
		//	acceptRulesBypass(token, joinreq.JSON(), serverID, inviteLink)
		//}
		//if changeNick {
		//	changeNicker(token, serverID, nickname)
		//}
	} else if joinreq.StatusCode == 403 {
		if strings.Contains(string(body), "You need to verify your account in order to perform this action.") ||
			strings.Contains(string(body), "このユーザーは、このギルドからBANされています。") ||
			strings.Contains(string(body), "The user is banned from this guild.") {
			fmt.Printf("Banned from Server | %s\n", extractToken)
			//fmt.Println("失敗しました")
		}
	}
}

func deleteJoinMsg(token, joinChannelID string) {
	extractToken := fmt.Sprintf("%s.%s", strings.Split(extract(token+"]"), ".")[0], strings.Split(extract(token+"]"), ".")[1])
	reqHeader := requestHeader(token, false, false)
	headers := reqHeader

	client := http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://discord.com/api/v9/channels/%s/messages?limit=100", joinChannelID), nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to fetch messages: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var messages []map[string]interface{}
	if err := json.Unmarshal(body, &messages); err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	for _, message := range messages {
		botTokenID, err := base64.StdEncoding.DecodeString(strings.Split(token, ".")[0] + "==")
		if err != nil {
			log.Fatalf("Failed to decode bot token ID: %v", err)
		}
		if message["content"] == "" && string(botTokenID) == message["author"].(map[string]interface{})["id"].(string) {
			req, err := http.NewRequest("DELETE", fmt.Sprintf("https://discord.com/api/v9/channels/%s/messages/%s", joinChannelID, message["id"].(string)), nil)
			if err != nil {
				log.Fatalf("Failed to create request: %v", err)
			}
			for key, value := range headers {
				req.Header.Set(key, value)
			}

			resp, err := client.Do(req)
			if err != nil {
				log.Fatalf("Failed to delete join message: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusNoContent {
				fmt.Printf("Success Delete Join Message: %s\n", extractToken)
			} else {
				fmt.Printf("Failed Delete Join Message: %s\n", extractToken)
				body, _ := ioutil.ReadAll(resp.Body)
				fmt.Println(string(body))
			}
			break
		}
	}
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

// Check Status
