package main

import (
	"log"
	"net/http"
	//"wxpay/pay"
	"github.com/willxm/WechatPayGo/pay"
)

func main() {
	http.HandleFunc("/downloadbill", pay.Bill)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
