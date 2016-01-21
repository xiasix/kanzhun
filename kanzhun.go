package main
import (
	"fmt"
	"net/http"
	"io/ioutil"
	"net/url"
	"flag"
	"strings"
	"encoding/json"
)
//登陆url
var loginurl string = "http://www.kanzhun.com/login.json"
//签到url
var checkurl string = "http://www.kanzhun.com/integral/userSign.json"


func main() {
	//用户名
	username := flag.String("username", "342450@qq.com", "请输入邮箱/手机号")
	//密码
	passwd := flag.String("passwd", "", "请输入密码")
	flag.Parse()

	//登陆获取cookie
	cookie := login(*username, *passwd)

	//签到
	check(cookie)


}
/*
	username string 用户名
	passwd string 密码
	return *cookie
 */
func login(username, passwd string) []*http.Cookie {
	//post的字段
	v := url.Values{}
	v.Set("redirect", "http://www.kanzhun.com/")
	v.Set("account", username)
	v.Set("password", passwd)
	v.Set("remember", "true")
	body := ioutil.NopCloser(strings.NewReader(v.Encode()))

	client := &http.Client{}
	request, err := http.NewRequest("POST", loginurl, body)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value") //必须加
	request.Header.Set("Host", "www.kanzhun.com")
	request.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	request.Header.Set("Referer", "http://www.kanzhun.com/login/?ka=head-signin")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36")
	request.Header.Set("X-Requested-With", "XMLHttpRequest")


	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(request)

	defer resp.Body.Close()

	return resp.Cookies()


}

func check(cookies []*http.Cookie) {

	client := &http.Client{}
	request, err := http.NewRequest("GET", checkurl, nil)
	if err != nil {
		fmt.Println(err)
	}
	//设置cookie
	for _, cookie := range cookies {
		request.AddCookie(cookie)

	}
	request.Header.Set("Host", "www.kanzhun.com")
	request.Header.Set("Proxy-Connection", "keep-alive")
	request.Header.Set("Pragma", "no-cache")
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	request.Header.Set("X-Requested-With", "XMLHttpRequest")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36")
	request.Header.Set("Referer", "http://www.kanzhun.com/timeline/?ka=top-menu1")
	request.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
	request.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		data, _ := ioutil.ReadAll(response.Body)
		var jsondata  map[string]interface{}

		jsonerr := json.Unmarshal([]byte(data), &jsondata)

		if jsonerr != nil {
			fmt.Println("Unmarshal: ", jsonerr.Error())
		}
		//登陆成功
		if (int(jsondata["rescode"].(float64)) == 1) {
			fmt.Println("签到成功")
		}else {
			fmt.Println("签到失败,失败原因:", string(data))
		}


	}

}

