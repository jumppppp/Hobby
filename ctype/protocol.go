package ctype

import "fmt"

var (
	DScoketPort = "1,7,9,13,19,21-23,25,37,42,49,53,69,79-81,85,105,109-111,113,123,135,137-139,143,161,179,222,264,384,389,402,407,443-446,465,500,502,512-515,523-524,540,548,554,587,617,623,689,705,771,783,873,888,902,910,912,921,993,995,998,1000,1024,1030,1035,1090,1098-1103,1128-1129,1158,1199,1211,1220,1234,1241,1300,1311,1352,1433-1435,1440,1494,1521,1530,1533,1581-1582,1604,1720,1723,1755,1811,1900,2000-2001,2049,2082,2083,2100,2103,2121,2199,2207,2222,2323,2362,2375,2380-2381,2525,2533,2598,2601,2604,2638,2809,2947,2967,3000,3037,3050,3057,3128,3200,3217,3273,3299,3306,3311,3312,3389,3460,3500,3628,3632,3690,3780,3790,3817,4000,4322,4433,4444-4445,4659,4679,4848,5000,5038,5040,5051,5060-5061,5093,5168,5247,5250,5351,5353,5355,5400,5405,5432-5433,5498,5520-5521,5554-5555,5560,5580,5601,5631-5632,5666,5800,5814,5900-5910,5920,5984-5986,6000,6050,6060,6070,6080,6082,6101,6106,6112,6262,6379,6405,6502-6504,6542,6660-6661,6667,6905,6988,7001,7021,7071,7080,7144,7181,7210,7443,7510,7579-7580,7700,7770,7777-7778,7787,7800-7801,7879,7902,8000-8001,8008,8014,8020,8023,8028,8030,8080-8082,8087,8090,8095,8161,8180,8205,8222,8300,8303,8333,8400,8443-8444,8503,8800,8812,8834,8880,8888-8890,8899,8901-8903,9000,9002,9060,9080-9081,9084,9090,9099-9100,9111,9152,9200,9390-9391,9443,9495,9809-9815,9855,9999-10001,10008,10050-10051,10080,10098,10162,10202-10203,10443,10616,10628,11000,11099,11211,11234,11333,12174,12203,12221,12345,12397,12401,13364,13500,13838,14330,15200,16102,17185,17200,18881,19300,19810,20010,20031,20034,20101,20111,20171,20222,22222,23472,23791,23943,25000,25025,26000,26122,27000,27017,27888,28222,28784,30000,30718,31001,31099,32764,32913,34205,34443,37718,38080,38292,40007,41025,41080,41523-41524,44334,44818,45230,46823-46824,47001-47002,48899,49152,50000-50004,50013,50500-50504,52302,55553,57772,62078,62514,65535"
)

func HTTP_(ip string, port string) []byte {
	return []byte(fmt.Sprintf("GET / HTTP/1.1\r\nHost: %v:%v\r\nConnection: close\r\n\r\n", ip, port))

}

func FTP_POP3_(name string) []byte {
	return []byte(fmt.Sprintf("USER %v", name))

}
func SSH_(name string) []byte {
	return []byte(fmt.Sprintf("SSH %v", name))

}
func SMTP_(domain string) []byte {
	return []byte(fmt.Sprintf("HELO %v", domain))

}
func DNS_() []byte {
	return []byte{
		0xaa, 0xaa, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77,
		0x07, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
		0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01,
	}
}
func NULL_() []byte {
	return []byte{}
}
