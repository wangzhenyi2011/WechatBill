package pay

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/willxm/WechatPayGo/tool"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	//"os"
	//"encoding/csv"
)

type BillReq struct {
	Appid     string `xml:"appid"`
	Bill_date string `xml:"bill_date"`
	Bill_type string `xml:"bill_type"`
	Mch_id    string `xml:"mch_id"`
	Nonce_str string `xml:"nonce_str"`
	Sign      string `xml:"sign"`
}

type BillResp struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
}

type TotalBill struct {
	Total        int
	Total_amount float32
}

func Bill(w http.ResponseWriter, r *http.Request) {

	/*
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println("fail to read body: ", err)
				http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			fmt.Println("app request，HTTP Body: ", string(body))

		  /*
			err = json.Unmarshal(body, &clientReq)
			if err != nil {
				fmt.Println("decode http body error: ", err)
				http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
	*/

	value := r.FormValue("data")

	//fmt.Println(value)

	//fmt.Println("app request，HTTP Body: ", string(value))

	var clinetResp BillResp
	the_time, err := time.Parse("2006-01-02", value)
	if err != nil {
		fmt.Println("new http request fail: ", err)
		msg := BillResp{
			Return_code: "BadRequest",
			Return_msg:  "日期格式参数错误，格式为:2017-03-06",
		}
		jsonBytes, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("encode json error: ", err)
		}
		w.Write(jsonBytes)

		//  fmt.Println("decode http body error: ", _err)
		//    http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		//    		fmt.Println("new http request fail: ", err)
		return
	}

	var Bill_date string
	if value == "" {
		nTime := time.Now()
		nTime.Format("2006-01-02 15:04:05")
		yesTime := nTime.AddDate(0, 0, -1)
		Bill_date = yesTime.Format("20060102")
	} else {
		yesTime := the_time.Format("20060102")
		Bill_date = yesTime
	}

	fmt.Println(Bill_date)

	//test data
	var clientReq BillReq

	clientReq.Appid = ""
	clientReq.Bill_date = Bill_date
	// clientReq.Bill_date = "20170316"
	clientReq.Bill_type = ""
	clientReq.Mch_id = ""
	clientReq.Nonce_str = ""


	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appid"] = clientReq.Appid
	m["bill_date"] = clientReq.Bill_date
	m["bill_type"] = clientReq.Bill_type
	m["mch_id"] = clientReq.Mch_id
	m["nonce_str"] = clientReq.Nonce_str

	clientReq.Sign = tool.WxpayCalcSign(m, "")
	
	//xml encoding
	bytesReq, err := xml.Marshal(clientReq)
	if err != nil {
		fmt.Println("xml encoding fail: ", err)
		return
	}
	strReq := string(bytesReq)
	strReq = strings.Replace(strReq, "BillReq", "xml", -1)
	//fmt.Println(strReq)
	bytesReq = []byte(strReq)

	reqURL := "https://api.mch.weixin.qq.com/pay/downloadbill"
	req, err := http.NewRequest("POST", reqURL, bytes.NewReader(bytesReq))
	if err != nil {
		fmt.Println("new http request fail: ", err)
		return
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println("wxpay api send fail", err)
		return
	}
	bytesResp, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(bytesResp))

	_err := xml.Unmarshal(bytesResp, &clinetResp)
	if _err != nil {
		strResp := strings.Split(string(bytesResp), "\n")

		total := strings.Split(strings.Join(strings.Split(strResp[len(strResp)-2], "`"), ""), ",")

		a, error := strconv.Atoi(string(total[0]))
		if error != nil {
			fmt.Println("字符串转换成整数失败", error)
		}

		b, error := strconv.ParseFloat(string(total[1]), 32)
		if error != nil {
			fmt.Println("字符串转换成浮点失败", error)
		}

		bill := TotalBill{
			Total:        a,
			Total_amount: float32(b),
		}

		/*
		  //  strResp_total:= strings.Split(strResp[1],"`")

		*/

		//  total := map[string]string{"总交易单数":strResp_total[1], "总交易额":strResp_total[2], "总退款金额":strResp_total[3]}

		order_amount, order_count := tool.Query(Bill_date)

		var msg BillResp

		if m := float32(0); order_amount == m {
			msg.Return_code = "FAIL"
			msg.Return_msg = "数据库无结果"
		} else if bill.Total == order_count && bill.Total_amount == order_amount {
			msg.Return_code = "Success"
			msg.Return_msg = "ok"

		} else if bill.Total == order_count && bill.Total_amount != order_amount {
			msg.Return_code = "FAIL"
			msg.Return_msg = "账单金额不对"
		} else if bill.Total != order_count && bill.Total_amount == order_amount {
			msg.Return_code = "FAIL"
			msg.Return_msg = "账单数目不对"
		} else {
			msg.Return_code = "FAIL"
			msg.Return_msg = "数目和金额不对"
		}

		jsonBytes, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("encode json error: ", err)
		}
		w.Write(jsonBytes)
		return
	}
	fmt.Println("wxpay api return data: ", string(bytesResp))
	fmt.Println("encoding data: ", clinetResp)
	jsonBytes, err := json.Marshal(clinetResp)
	if err != nil {
		fmt.Println("encode json error: ", err)
	}
	w.Write(jsonBytes)

	/*




	    fmt.Println(strResp_total[2])
	    //str := strings.Split(strResp_total[1],",")


	    rating := map[string]string{"总交易单数":strResp_total[1], "总交易额":strResp_total[2], "总退款金额":strResp_total[3]}

	    wechatBill := tool.Query()
	    fmt.Println(wechatBill)

	    if rating["总交易单数"] == wechatBill {
	  		fmt.Sprintln("verify success")
	  	//	return true
	  	}
	  	fmt.Println("verify fail")


	*/

	/*
		    f, err := os.Create("test.csv")//创建文件
			  if err != nil {
				     panic(err)
			  }
		    defer f.Close()

			  f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

			  w := csv.NewWriter(f)//创建一个新的写入文件流

			  w.Write(strResp)//写入数据
		    w.Flush()

		    jsonBytes,err := tool.ReadFile("./test.csv")
		    if err != nil {
		         fmt.Println("wxpay api send fail", err)
		    }

		    fmt.Println(string(jsonBytes))
	*/

}
