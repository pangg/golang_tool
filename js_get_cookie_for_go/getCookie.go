package main

import (
	"os"
	//"fmt"
	//"strings"
	"os/exec"
	log "github.com/kdar/factorlog"
	"flag"
)

var (
	logFlag        = flag.String("log", "", "set log path")
	logger *log.FactorLog
)

const (
	_COOKIES_PATH_ = "./cookies/"
)

func main () {
	logger = SetGlobalLogger(*logFlag)
	var res []byte
	var err error
	var cmd *exec.Cmd

	param := []string{"ashiya77", "aaaa5555", "tech@tianxi100.com", "aH1GA+x"}

	// 执行单个shell命令时, 直接运行即可
	cmd = exec.Command("phantomjs", "p2.js", param[0], param[1], param[2], param[3])
	if res, err = cmd.Output(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}


	// 默认输出有一个换行
	logger.Info(string(res))


	cookieFile := _COOKIES_PATH_+param[0]+"-"+param[2]
	isExist, _ := PathExists(cookieFile)
	if isExist {

	}else{

	}

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