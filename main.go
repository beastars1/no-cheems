package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"
)

var (
	urls   = make(chan string)
	limit  = make(chan int, 20)
	result = make(chan string)
	quit   = make(chan bool)
	urlMap = make(map[string]bool)
)

var (
	start = 0
	end   = 1000
)

var (
	urlBodyReMap = make(map[string]func(string))
)

func init() {
	urlBodyReMap["http://mp4.yyhgxgy.cn:520/mfdsp/mf1.php?id=%d?_wv=xw.qq.com"] = Re1
	urlBodyReMap["http://1162as.oss-cn-beijing.aliyuncs.com/%03d.php"] = Re2
}

func main() {
	for format, re := range urlBodyReMap {
		do(format, re)
	}
}

func do(format string, re func(string)) {
	now := time.Now()
	go generateUrl(format)
	go func() {
		for {
			select {
			case res := <-result:
				urlMap[res] = true
			case url := <-urls:
				go func() {
					limit <- 1
					request(url, re)
					<-limit
				}()
			case <-time.After(3 * time.Second):
				quit <- true
			}
		}
	}()
	<-quit
	file, err := OpenFile("./" + time.Now().Format("20060102-150105") + ".txt")
	if err != nil {
		fmt.Println(err)
	}
	write(file)
	fmt.Println(format+":", time.Since(now)-3*time.Second)
}

func doResult(s string, file *os.File) {
	write := bufio.NewWriter(file)
	write.WriteString(s + "\n")
	write.Flush()
}

func write(file *os.File) {
	fmt.Println("write:" + file.Name())
	var set []string
	for url := range urlMap {
		set = append(set, url)
		delete(urlMap, url)
	}
	sort.Strings(set)
	l := len(set)
	for i := range set {
		doResult(set[l-i-1], file)
	}
}

func generateUrl(format string) {
	for i := start; i < end; i++ {
		url := fmt.Sprintf(format, i)
		urls <- url
	}
}

func request(url string, re func(string)) {
	client := http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("User-Agent", "QQ/114.514")
	resp, err := client.Do(request)
	if err != nil || resp == nil || resp.StatusCode != 200 {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	re(string(body))
}
