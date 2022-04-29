package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

var (
	urls   = make(chan string)
	limit  = make(chan int, 20)
	result = make(chan string)
	quit   = make(chan bool)
)

var (
	num = 10000
)

func main() {
	start := time.Now()
	go generateUrl()
	go func() {
		file, _ := OpenFile("./urls.txt")
		for {
			select {
			case res := <-result:
				doResult(res, file)
			case url := <-urls:
				go func() {
					limit <- 1
					do(url)
					<-limit
				}()
			case <-time.After(3 * time.Second):
				quit <- true
			}
		}
	}()
	<-quit
	fmt.Println(time.Since(start))
}

func doResult(s string, file *os.File) {
	write := bufio.NewWriter(file)
	write.WriteString(s + "\n")
	write.Flush()
}

func generateUrl() {
	for i := 0; i < num; i++ {
		url := fmt.Sprintf("http://mp4.yyhgxgy.cn:520/mfdsp/mf1.php?id=%d?_wv=xw.qq.com", i)
		urls <- url
	}
}

func do(url string) {
	client := http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 12; Redmi K30 Build/SKQ1.210908.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/89.0.4389.72 MQQBrowser/6.2 TBS/046011 Mobile Safari/537.36 V1_AND_SQ_8.8.85_2712_YYB_D A_8088500 QQ/8.8.85.7685 NetType/WIFI WebP/0.3.0 Pixel/1080 StatusBarHeight/96 SimpleUISwitch/0 QQTheme/1000 InMagicWin/0 StudyMode/0 CurrentMode/0 CurrentFontScale/1.0 GlobalDensityScale/0.9818182 AppId/1")
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	strs := regexpMatch(string(body), `http:([\s\S]*?)m3u8`)
	if len(strs) > 0 {
		result <- strs[0]
	}
}

func regexpMatch(text, expr string) []string {
	reg, err := regexp.Compile(expr)
	if err != nil {
		fmt.Println("regexp cant compile", err)
		return nil
	}
	return reg.FindStringSubmatch(text)
}
