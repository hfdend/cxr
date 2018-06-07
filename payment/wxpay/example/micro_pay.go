package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"git.jiayougougou.com/snail/transaction/pay/wxpay"
)

func main() {
	var authCode string
	flag.StringVar(&authCode, "auth-code", "", "微信条码")
	bts, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalln(err)
	}
	flag.Parse()

	config := wxpay.PayConfig{}
	if err := json.Unmarshal(bts, &config); err != nil {
		log.Fatalln(err)
	}
	fmt.Println(config)

	api := wxpay.NewApi(config)
	args := wxpay.NewPayUnifiedOrder()
	args.SetBody("刷卡测试样例-支付")
	args.SetOutTradeNo(fmt.Sprintf("gotest%s", time.Now().Format("20060102150405")))
	args.SetTotalFee(1)
	args.SetAuthCode(authCode)
	result, err := api.MicroPay(args, 5*time.Second)
	fmt.Println(err)
	fmt.Println(result)
}
