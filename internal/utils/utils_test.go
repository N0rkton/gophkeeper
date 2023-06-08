package utils

import "fmt"

func ExampleGenerateRandomString() {
	str := GenerateRandomString(5)
	fmt.Println(len(str))
	//Output:
	//8
}
func ExampleGetMD5Hash() {
	hash := GetMD5Hash("test")
	fmt.Println(hash)
	//Output:
	//098f6bcd4621d373cade4e832627b4f6
}
