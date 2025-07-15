package functioncall

import (
	"fmt"
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	now := time.Now().Local()
	s := fmt.Sprintf("当前时间是 %d年%d月%d日 %d时%d分", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
	fmt.Println(s)
}
