package TaibaiUtils

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"strconv"
	"sync"
	"time"
)

// 切勿直接使用这个
func generateMessageId() func() int {

	var GenerateMessageMutex sync.Mutex
	var LastMessageTimeS int
	var LastMessageIndex int

	return func() int {
		// ("2006-01-02 15:04:05.999999999 -0700 MST")
		// 年月日十分秒*1万 + 此秒内的index
		GenerateMessageMutex.Lock()
		defer GenerateMessageMutex.Unlock()

		loc, _ := time.LoadLocation("Asia/Chongqing")
		beijingTimeStr := time.Now().In(loc).Format("20060102150405")
		beijingTime, _ := strconv.Atoi(beijingTimeStr)
		if beijingTime != LastMessageTimeS {
			LastMessageIndex = 1
		} else {
			LastMessageIndex = LastMessageIndex + 1
		}

		return beijingTime*10000 + LastMessageIndex
	}

}

var GenerateMessageId func() int

func init() {
	GenerateMessageId = generateMessageId()
}

func SimpleJsonToStruct(from *simplejson.Json, to interface{}) error {
	str, _ := json.Marshal(from.Interface())
	err := json.Unmarshal(str, to)
	return err
}
