package helpers

import (
	"time"
	"fmt"
	"encoding/base64"
	"crypto/md5"
	"strings"
	"strconv"
	"asyncTask/config"
)

var h = md5.New()

func cipherEncode(sourceText string) string {
	h.Write([]byte(config.All().SecretKey))
	cipherHash := fmt.Sprintf("%x", h.Sum(nil))
	h.Reset()
	inputData := []byte(sourceText)
	loopCount := len(inputData)
	outData := make([]byte,loopCount)
	for i:= 0; i < loopCount ; i++ {
		outData[i] = inputData[i] ^ cipherHash[i%32]
	}
	return fmt.Sprintf("%s", outData)
}

func HlcEncode(sourceText string) string {
	h.Write([]byte(time.Now().Format("2006-01-02 15:04:05")))
	noise := fmt.Sprintf("%x", h.Sum(nil))
	h.Reset()
	inputData := []byte(sourceText+"-"+strconv.Itoa(int(time.Now().Unix())))
	loopCount := len(inputData)
	outData := make([]byte,loopCount*2)
	for i, j := 0,0; i < loopCount ; i,j = i+1,j+1 {
		outData[j] = noise[i%32]
		j++
		outData[j] = inputData[i] ^ noise[i%32]
	}
	return base64.StdEncoding.EncodeToString([]byte(cipherEncode(fmt.Sprintf("%s", outData))))
}

func HlcDecode(sourceText string) string {
	buf, err := base64.StdEncoding.DecodeString(sourceText)
	if err != nil {
		fmt.Println("Decode(%q) failed: %v", sourceText, err)
		return ""
	}
	inputData := []byte(cipherEncode(fmt.Sprintf("%s", buf)))
	loopCount := len(inputData)
	outData := make([]byte,loopCount)
	for i, j := 0,0; i < loopCount ; i,j = i+2,j+1 {
		outData[j] = inputData[i] ^ inputData[i+1]
	}
	str := fmt.Sprintf("%s", outData)
	arr := strings.Split(str,"-")
	if len(arr) < 2 {
		return ""
	}
	//cur := time.Now()
	//t,err := strconv.Atoi(fmt.Sprintf("%s",bytes.Trim([]byte(arr[len(arr)-1]), "\x00")))
	if err != nil{
		fmt.Println(err.Error())
		return ""
	}
	//if cur.Unix() - int64(t) > 5{
	//	fmt.Println(arr[1])
	//	return ""
	//}
	result := strings.Join(arr[:len(arr)-1],"")
	return fmt.Sprintf("%s", result)
}