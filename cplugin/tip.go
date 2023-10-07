package cplugin

import "time"

func CSleep(num int) {
	time.Sleep(time.Duration(num) * time.Second)
}
