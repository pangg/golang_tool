package main

import (
    "fmt"
    "io/ioutil"
    "encoding/json"
    "strings"
    "net/http"
    "bytes"
    "golang.org/x/text/encoding/japanese"
    "golang.org/x/text/transform"
)

type Cookie struct {
    Name string     `json:name`
    Value string    `json:value`
    Domain string   `json:domain`
}

func main() {
    var cookies []*Cookie
    content, _ := ioutil.ReadFile("./cookies/ashiya77-tech@tianxi100.com.json")
    json.Unmarshal(content, &cookies)

    req, _ := http.NewRequest("GET", "https://rdatatool.rms.rakuten.co.jp/access/?menu=pc&evt=RT_P03_01&stat=1&owin=", nil)
    for _, cookie := range cookies {
        req.AddCookie(&http.Cookie{
            Name: cookie.Name,
            Value: cookie.Value,
            Domain: cookie.Domain,
        })
    }
    client := &http.Client{}
    res, _ := client.Do(req)

    reqAccess, _ := http.NewRequest("POST", "https://rdatatool.rms.rakuten.co.jp/access/", strings.NewReader(fmt.Sprintf("owin=&menu=pc&evt=RT_P03_01&type=week&y=2018&w=1001")))
    for _, cookie := range cookies {
        reqAccess.AddCookie(&http.Cookie{
            Name: cookie.Name,
            Value: cookie.Value,
            Domain: cookie.Domain,
        })
    }
    client.Do(reqAccess)
    _, _ = ioutil.ReadAll(res.Body)
    defer res.Body.Close()

    reqDownload, _ := http.NewRequest("GET", "https://rdatatool.rms.rakuten.co.jp/access/?menu=pc&evt=RT_D01_02&category=RT_P03_01&chk=1539243783&owin=", nil)
    for _, cookie := range cookies {
        reqDownload.AddCookie(&http.Cookie{
            Name: cookie.Name,
            Value: cookie.Value,
            Domain: cookie.Domain,
        })
    }

    /*fileName := "./report.csv"
    output, err := os.Create(fileName)
    defer output.Close()*/

    reqDownload.Header.Add("Host", "rdatatool.rms.rakuten.co.jp")
    reqDownload.Header.Add("Connection", "keep-alive")
    reqDownload.Header.Add("Upgrade-Insecure-Requests", "1")
    reqDownload.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36")
    //reqDownload.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
    // reqDownload.Header.Add("Accept-Encoding", "gzip, deflate, br")
    //reqDownload.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7,ja;q=0.6")
    //reqDownload.Header.Add("Pragma", "no-cache")
    //reqDownload.Header.Add("Cache-Control", "no-cache")
    reqDownload.Header.Add("Charset", "Shift_JIS")
    res, _ = client.Do(reqDownload)
    defer res.Body.Close()

    body := res.Body

    buf := new(bytes.Buffer)
    buf.ReadFrom(body)
    str := buf.String()

    var reader *transform.Reader
    reader = transform.NewReader(bytes.NewReader([]byte(str)), japanese.ShiftJIS.NewDecoder())
    d, err := ioutil.ReadAll(reader)
    if err != nil {
        panic(err)
    }
    str = string(d)

    var records [][]string
    rows := strings.Split(str, "\n")
    if len(rows) > 0 {
        for _, row := range rows {
            lines := strings.Split(row, ",")
            records = append(records, lines)
        }
    }

    fmt.Println(records, "==============")

    /*_, err = io.Copy(output, body)
    if err != nil {
        panic(err)
    }*/




}

func Json(obj interface{}) string {
    var buf bytes.Buffer
    enc := json.NewEncoder(&buf)
    enc.Encode(obj)

    return buf.String()
}
