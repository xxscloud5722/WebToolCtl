package main

import (
	"fmt"
	"github.com/longyuan/domain.v3/client"
	"github.com/longyuan/lib.v3/times"
	"testing"
	"time"
)

func Test1(test *testing.T) {
	ssl, err := client.SSL("www.jd.com:6443")
	if err != nil {
		return
	}
	fmt.Println(ssl)
	info, err := client.IpInfo("124.222.4.134")
	if err != nil {
		return
	}
	fmt.Println(info)
}

func Test2(test *testing.T) {
	utcTime := time.Date(2023, 7, 26, 12, 34, 56, 0, time.UTC)

	shanghaiTimeStr := times.In(utcTime).Format("2006-01-02 15:04:05")
	fmt.Println("上海时区的datetime格式:", shanghaiTimeStr)

}

func Test3(test *testing.T) {

}
