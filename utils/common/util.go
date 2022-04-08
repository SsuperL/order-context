package common

import (
	"fmt"
	"math"
	"math/rand"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const alphabet = "1234567890"

// GenerateNumber 根据时间生成编号
func GenerateNumber() string {
	year, month, day := time.Now().Year(), strconv.Itoa(int(time.Now().Month())), strconv.Itoa(time.Now().Day())
	if newMonth, _ := strconv.Atoi(month); newMonth < 10 {
		month = "0" + month
	}
	if newDay, _ := strconv.Atoi(day); newDay < 10 {
		day = "0" + day
	}
	stamp := time.Now().UnixNano()/int64(1000) - int64(math.Pow10(15))
	return fmt.Sprintf("%d%s%s%d", year, month, day, stamp)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// GetProjectAbPathByCaller get absolute path of file
func GetProjectAbPathByCaller() (abPath string) {
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		filePath := path.Dir(filename)
		abPath = path.Dir(filePath)
	}
	return
}
