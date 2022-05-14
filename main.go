package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"time"
)

var (
	urls   = make(chan string)
	limit  = make(chan int, 20)
	result = make(chan string)
	quit   = make(chan bool)
)

var (
	start = 0
	end   = 10000
)

var urlMap = make(map[string]bool)

func main() {
	now := time.Now()
	go generateUrl()
	go func() {
		for {
			select {
			case res := <-result:
				urlMap[res] = true
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
	file, _ := OpenFile("./urls-" + time.Now().Format("20060102") + ".txt")
	write(file)
	fmt.Println(time.Since(now) - 3*time.Second)
}

func doResult(s string, file *os.File) {
	write := bufio.NewWriter(file)
	write.WriteString(s + "\n")
	write.Flush()
}

func write(file *os.File) {
	var set []string
	for url := range urlMap {
		set = append(set, url)
	}
	sort.Strings(set)
	l := len(set)
	for i := range set {
		doResult(set[l-i-1], file)
	}
}

func generateUrl() {
	for i := start; i < end; i++ {
		url := fmt.Sprintf("http://mp4.yyhgxgy.cn:520/mfdsp/mf1.php?id=%d?_wv=xw.qq.com", i)
		urls <- url
	}
}

func do(url string) {
	client := http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("User-Agent", "QQ/114.514")
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	strs := regexpMatch(string(body), `"http([\s\S]*?)m3u8`)
	if len(strs) > 0 {
		s := strs[0]
		result <- s[1:]
	}
}

func regexpMatch(text, expr string) []string {
	reg, err := regexp.Compile(expr)
	if err != nil {
		return nil
	}
	return reg.FindStringSubmatch(text)
}
