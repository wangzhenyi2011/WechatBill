package tool

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func Query(bill_time string) (float32, int) {
	db, err := sql.Open("mysql", "")

	if err != nil {
		fmt.Println("db QueryRow error", err)
	}
	defer db.Close()


	var order_amount float32
	var order_count int

	tm2, _ := time.Parse("20060102", bill_time)
	tm1 := time.Date(tm2.Year(), tm2.Month(), tm2.Day(), 0, 0, 0, 0, time.Local)
	tm3 := tm1.AddDate(0, 0, 1)

	//fmt.Println(tm3.Unix())

	err = db.QueryRow("", tm1.Unix(), tm3.Unix()).Scan(&order_amount, &order_count)

	if err != nil {
		if err == sql.ErrNoRows {
      fmt.Println("QueryRow  is null")
			// there were no rows, but otherwise no error occurred
		} else {
      fmt.Println(order_amount,order_count,"查询结果")
		}
	}

	//  defer rows.Close()

	return order_amount/100, order_count

	//普通demo
	//for rows.Next() {
	//	var userId int
	//	var userName string
	//	var userAge int
	//	var userSex int

	//	rows.Columns()
	//	err = rows.Scan(&userId, &userName, &userAge, &userSex)
	//	checkErr(err)

	//	fmt.Println(userId)
	//	fmt.Println(userName)
	//	fmt.Println(userAge)
	//	fmt.Println(userSex)
	//}

	//字典类型
	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	/*
		columns, _ := rows.Columns()
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		for rows.Next() {
			//将行数据保存到record字典
			err = rows.Scan(scanArgs...)
			record := make(map[string]string)
			for i, col := range values {
				if col != nil {
					record[columns[i]] = string(col.([]byte))
				}
			}
			fmt.Println(record)
		}
	*/
}
