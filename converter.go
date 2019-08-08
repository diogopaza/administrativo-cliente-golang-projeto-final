package main

import (
	"fmt"
	"time"
)

func main() {
	/*
					s := "2147483647"
					i64, err := strconv.ParseInt(s, 10, 32)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(i64)
					fmt.Println(s)
					eval(s)

				hora_ini = "08:40:20"
				fmt.Println(hora_ini.String())
				fmt.Println(t.Format("15:04:05"))

			layout := "2014-09-12T11:45:26.371Z"
			str := "2015-11-12T11:45:26.371Z"
			t, err := time.Parse(layout, str)
			_ = err
			fmt.Println(t)

		yourString := "10:06:00"
		yourDate, err := time.Parse("03:04:00", yourString)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(yourDate)
	*/
	t, _ := time.Parse(time.RFC3339, "Jan 2, 2006 at 3:04pm (MST)")
	fmt.Println(t)
	/*
		str := "07:30:45"
		layout := "03:04:05"
		tm, err := time.Parse(layout, str)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(tm)
	*/
}
