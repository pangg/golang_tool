package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/kdar/factorlog"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	logFlag    = flag.String("log", "", "set log path")
	logger     *log.FactorLog
	weekValue  []string
	tasks      []*Params
	shopEnName string
)

const (
	_COOKIES_PATH_           = "./cookies/"
	_R_ITEM_DATA_            = "https://rdatatool.rms.rakuten.co.jp/access/?menu=pc&evt=RT_P03_01&stat=1&owin="
	_R_DOWNLOAD_ITEM_REPORT_ = "https://rdatatool.rms.rakuten.co.jp/access/"
	_SAVING_REPORTS_         = "http://rakuten.gaopan.xibao100.com/report/report/rakuten_product_info_saving"
)

type Params struct {
	Menu  string
	Type  string
	Year  int
	Month int
	Day   int
	Week  string
	Date  string
}

type Cookie struct {
	Name   string `json:name`
	Value  string `json:value`
	Domain string `json:domain`
}

func main() {
	var cookies []*Cookie
	logger = SetGlobalLogger(*logFlag)
	param := []string{"ZpJ9q7wYuj", "kb2018wine18u", "tech@tianxi100.com", "aH1GA+x"}

	//getCk := GetCookies(param, 3)
	getCk := true
	if getCk {
		cookieFile := _COOKIES_PATH_ + param[0] + "-" + param[2] + ".json"
		isExist, _ := PathExists(cookieFile)
		if !isExist {
			logger.Info("Cookie file not exist")
			return
		}
		content, _ := ioutil.ReadFile(cookieFile)
		json.Unmarshal(content, &cookies)

		isUse := CheckCookies(cookies, 3)

		_ = getParams()

		tasks = tasks[0:15]

		//isUse := true
		if isUse && len(tasks) > 0 {
			for _, task := range tasks {
				logger.Info("start: ", *task)
				isDown := RequestAccessItemCsv(&cookies, *task, 3)

				if isDown {
					reports := DownloadItemCsv(cookies, 3)
					reportJson, _ := json.Marshal(reports)

					SaveReports(string(reportJson), *task, 3)
				}
			}
		}
	} else {
		logger.Info("Can't get cookies, please checking")
	}

	return
}

/*请求接口， 保存报表数据*/
func SaveReports(reports string, task Params, retries int) bool {
	if retries > 0 {
		retries--
		paraTime, _ := time.Parse("2006-01-02", task.Date)
		date := paraTime.Format("2006/01/02")
		param := "shop_en_name=" + shopEnName + "&period=" + task.Type + "&date=" + date + "&device=" + task.Menu + "&datas=" + reports

		client := &http.Client{}
		req, err := http.NewRequest("POST", _SAVING_REPORTS_, strings.NewReader(param))
		if err != nil {
			logger.Error(err)
			SaveReports(reports, task, retries)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error(err)
			SaveReports(reports, task, retries)
		}

		logger.Info(task, string(body))
		return true
	}
	return false
}

/*下载商品csv*/
func DownloadItemCsv(cookies []*Cookie, retries int) (records [][]string) {
	//logger.Info("Download item CSV, times: ", retries)
	now := time.Now()
	uinxTime := now.Unix()
	downloadUrl := _R_DOWNLOAD_ITEM_REPORT_ + "?menu=pc&evt=RT_D01_02&category=RT_P03_01&chk=" + strconv.FormatInt(uinxTime, 10) + "&limit=15&limit_pt=&limit_mn=&limit_gn=&limit_ci=&owin="

	if retries > 0 {
		retries--
		reqDownload, _ := http.NewRequest("GET", downloadUrl, nil)
		for _, cookie := range cookies {
			reqDownload.AddCookie(&http.Cookie{
				Name:   cookie.Name,
				Value:  cookie.Value,
				Domain: cookie.Domain,
			})
		}
		reqDownload.Header.Add("Host", "rdatatool.rms.rakuten.co.jp")
		reqDownload.Header.Add("Connection", "keep-alive")
		reqDownload.Header.Add("Upgrade-Insecure-Requests", "1")
		reqDownload.Header.Add("Charset", "Shift_JIS")
		client := &http.Client{
			Timeout: 15 * time.Second,
		}
		resp, err := client.Do(reqDownload)
		if err != nil {
			DownloadItemCsv(cookies, retries)
		}
		defer resp.Body.Close()
		body := resp.Body
		buf := new(bytes.Buffer)
		buf.ReadFrom(body)
		str := buf.String()

		var reader *transform.Reader
		reader = transform.NewReader(bytes.NewReader([]byte(str)), japanese.ShiftJIS.NewDecoder())
		d, err := ioutil.ReadAll(reader)
		if err != nil {
			DownloadItemCsv(cookies, retries)
		}
		str = string(d)
		rows := strings.Split(str, "\n")
		if len(rows) > 0 {
			for _, row := range rows {
				lines := strings.Split(row, ",")
				if len(lines) > 0 {
					var t []string
					for _, v := range lines {
						t = append(t, strings.Trim(v, `"`))
					}
					records = append(records, t)
				}
			}
		}
		return
	}
	return
}

/*请求获取需要下载的商品csv*/
func RequestAccessItemCsv(cookies *[]*Cookie, param Params, retries int) bool {
	//logger.Info("Request csv access, times: ", retries)
	if retries > 0 {
		retries--
		var uri string
		if param.Type == "day" {
			uri = "owin=&menu=" + param.Menu + "&evt=RT_P03_01&type=" + param.Type + "&y=" + strconv.Itoa(param.Year) + "&m=" + strconv.Itoa(param.Month) + "&d=" + strconv.Itoa(param.Day)
		} else if param.Type == "week" {
			uri = "owin=&menu=" + param.Menu + "&evt=RT_P03_01&type=" + param.Type + "&y=" + strconv.Itoa(param.Year) + "&w=" + param.Week
		} else if param.Type == "month" {
			uri = "owin=&menu=" + param.Menu + "&evt=RT_P03_01&type=" + param.Type + "&y=" + strconv.Itoa(param.Year) + "&m=" + strconv.Itoa(param.Month)
		}

		reqAccess, _ := http.NewRequest("POST", _R_DOWNLOAD_ITEM_REPORT_, strings.NewReader(fmt.Sprintf(uri)))
		for _, cookie := range *cookies {
			reqAccess.AddCookie(&http.Cookie{
				Name:   cookie.Name,
				Value:  cookie.Value,
				Domain: cookie.Domain,
			})
		}
		reqAccess.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
		//reqAccess.Header.Add("Accept-Encoding", "gzip, deflate, br")
		//reqAccess.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,und;q=0.8")
		reqAccess.Header.Add("Cache-Control", "no-cache")
		reqAccess.Header.Add("Connection", "keep-alive")
		reqAccess.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		reqAccess.Header.Add("Host", "rdatatool.rms.rakuten.co.jp")
		reqAccess.Header.Add("Origin", "https://rdatatool.rms.rakuten.co.jp")
		reqAccess.Header.Add("Pragma", "no-cache")
		reqAccess.Header.Add("Referer", "https://rdatatool.rms.rakuten.co.jp/access/?menu=pc&evt=RT_P03_01&stat=1")
		reqAccess.Header.Add("Upgrade-Insecure-Requests", "1")
		reqAccess.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36")
		client := &http.Client{
			Timeout: 15 * time.Second,
		}
		resp, err := client.Do(reqAccess)
		if err != nil {
			logger.Error(err)
			RequestAccessItemCsv(cookies, param, retries)
		} else {
			defer resp.Body.Close()
			pageBody, _ := ioutil.ReadAll(resp.Body)
			err = ioutil.WriteFile("./rakuten.html", pageBody, 0755)
			if resp.StatusCode == 200 && strings.Contains(string(pageBody), "商品ページランキング") {
				return true
			} else {
				RequestAccessItemCsv(cookies, param, retries)
			}
		}
	}
	return false
}

/*判断cookies是否可用*/
func CheckCookies(cookies []*Cookie, retries int) bool {
	//logger.Info("check cookies retries: ", retries)
	if retries > 0 {
		retries--
		req, _ := http.NewRequest("GET", _R_ITEM_DATA_, nil)
		for _, cookie := range cookies {
			req.AddCookie(&http.Cookie{
				Name:   cookie.Name,
				Value:  cookie.Value,
				Domain: cookie.Domain,
			})
		}
		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			logger.Error(err)
			CheckCookies(cookies, retries)
		} else {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Error(err)
				CheckCookies(cookies, retries)
			}
			htmlStr := string(body)
			//err = ioutil.WriteFile("./rakuten.html", body, 0755)
			if strings.Contains(htmlStr, "商品ページランキング") {
				//解析页面， 获取周次数据参数
				doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
				if err != nil {
					logger.Error(err)
					CheckCookies(cookies, retries)
				}
				getWeekValues(doc)

				return true
			}
			CheckCookies(cookies, retries)
		}
	}
	logger.Info("Cookie Unusable")
	return false
}

func getWeekValues(doc *goquery.Document) {
	//周数据 参数
	doc.Find("#w option").Each(func(i int, s *goquery.Selection) {
		v, _ := s.Attr("value")
		if len(v) > 0 {
			weekValue = append(weekValue, v)
		}
	})
	//获取店铺英文名
	shopEnName, _ = doc.Find("#ratShopUrl").Attr("value")
	return
}

func GetCookies(param []string, retries int) bool {
	var res []byte
	//var err error
	var cmd *exec.Cmd
	//ck = false //初始 cookies 不可用
	logger.Info("Get cookies times: ", retries)

	if retries > 0 {
		retries--
		cmd = exec.Command("phantomjs", "p2.js", param[0], param[1], param[2], param[3])
		res, _ = cmd.Output()
		/*if res, err = cmd.Output(); err != nil {
			logger.Error(err)
			return ck
		}*/
		resInfo := string(res)
		logger.Info(resInfo)
		if strings.Contains(resInfo, "get_cookies_success") {
			return true
		} else {
			GetCookies(param, retries)
		}
	}
	return false
}

//获取参数
func getParams() bool {
	devices := []string{"pc", "mobile", "smp"}
	getTypesParams(devices, "day")
	getTypesParams(devices, "week")
	getTypesParams(devices, "month")

	return true
}

func getTypesParams(dev []string, types string) {
	now := time.Now()
	switch types {
	case "day":
		for _, p := range dev {
			for i := 1; i <= 60; i++ {
				preTime := now.AddDate(0, 0, -i)
				tmp := &Params{
					Menu:  p,
					Type:  types,
					Year:  preTime.Year(),
					Month: int(preTime.Month()),
					Day:   preTime.Day(),
					Week:  "",
					Date:  preTime.Format("2006-01-02"),
				}
				tasks = append(tasks, tmp)
			}
		}
		break
	case "week":
		preTime := now.AddDate(0, -6, 0)
		var years []int
		if now.Year() == preTime.Year() {
			years = []int{now.Year()}
		} else {
			years = []int{now.Year() - 1, now.Year()}
		}
		for _, p := range dev {
			for _, year := range years {
				for _, week := range weekValue {
					weekTime, _ := time.Parse("20060102", strconv.Itoa(year)+week)
					if weekTime.Unix() >= preTime.Unix() && weekTime.Unix() < now.Unix() {
						tmp := &Params{
							Menu:  p,
							Type:  types,
							Year:  year,
							Month: 0,
							Day:   0,
							Week:  week,
							Date:  weekTime.Format("2006-01-02"),
						}
						tasks = append(tasks, tmp)
					}
				}
			}
		}
		break
	case "month":
		preTime := now.AddDate(-2, 0, 0)
		var years []int
		for y := preTime.Year(); y <= now.Year(); y++ {
			years = append(years, y)
		}
		for _, p := range dev {
			for _, year := range years {
				for m := 1; m <= 12; m++ {
					var monthStr string
					if m < 10 {
						monthStr = "0" + strconv.Itoa(m)
					} else {
						monthStr = strconv.Itoa(m)
					}
					monthTime, _ := time.Parse("20060102", strconv.Itoa(year)+monthStr+"01")
					if monthTime.Unix() >= preTime.Unix() && monthTime.Unix() < now.Unix() {
						tmp := &Params{
							Menu:  p,
							Type:  types,
							Year:  year,
							Month: m,
							Day:   0,
							Week:  "",
							Date:  monthTime.Format("2006-01-02"),
						}
						tasks = append(tasks, tmp)
					}

				}
			}
		}

	}

	return
}

func SetGlobalLogger(logPath string) *log.FactorLog {
	sfmt := `%{Color "red:white" "CRITICAL"}%{Color "red" "ERROR"}%{Color "yellow" "WARN"}%{Color "green" "INFO"}%{Color "cyan" "DEBUG"}%{Color "blue" "TRACE"}[%{Date} %{Time}] [%{SEVERITY}:%{ShortFile}:%{Line}] %{Message}%{Color "reset"}`
	logger := log.New(os.Stdout, log.NewStdFormatter(sfmt))
	if len(logPath) > 0 {
		logf, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
		if err != nil {
			return logger
		}
		logger = log.New(logf, log.NewStdFormatter(sfmt))
	}
	logger.SetSeverities(log.INFO | log.WARN | log.ERROR | log.FATAL | log.CRITICAL)
	return logger
}

/*判断文件是否存在*/
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
