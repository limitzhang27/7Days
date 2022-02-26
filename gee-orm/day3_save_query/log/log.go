package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

// SetLevel controls log level
// 通过设置 level 设置 Output 来控制是否打印
func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	// level 大于 ErrorLevel, 则不打印Error级别日志
	if ErrorLevel > level {
		errorLog.SetOutput(ioutil.Discard)
	}

	// level 大于 InfoLevel，则不打印Info级别日志
	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}
	// 0 打印 Info, Error
	// 1 打印 Error,
	// 2 不打印
}
