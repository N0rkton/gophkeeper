package utils

import "fmt"

var mysecret = "qwertyuiopmmasdf"

func ExampleEncryptDecrypt() {
	str := Encrypt("example string", []byte(mysecret))
	str = Decrypt(str, []byte(mysecret))
	fmt.Println(str)
	// Output:
	//example string

}
