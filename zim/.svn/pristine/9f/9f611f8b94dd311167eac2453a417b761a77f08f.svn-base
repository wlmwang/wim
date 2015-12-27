package common

import (
	//"bytes"
	//"crypto/cipher"
	//"crypto/des"
	"crypto/md5"
	//"crypto/rand"
	"encoding/hex"
	"fmt"
	//"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var LogSvr *log.Logger

func GetLogger() *log.Logger {
	dir := "./log"
	t := time.Now()
	filepath := dir + "/log-" + strings.Replace(t.String()[:7], ":", "_", 3) + ".txt"
	logfile, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(-1)
	}
	logger := log.New(logfile, "", log.Ldate|log.Ltime|log.Llongfile)
	return logger
}

func GetCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	splitstring := strings.Split(path, "\\")
	size := len(splitstring)
	splitstring = strings.Split(path, splitstring[size-1])
	ret := strings.Replace(splitstring[0], "\\", "/", size-1)
	return ret
}

func GetDevice(r *http.Request) (device string) {
	UserAgent := r.Header.Get("User-Agent")
	if strings.Contains(UserAgent, "Mozilla") {
		device = "web"
	} else if strings.Contains(UserAgent, "Ios") {
		device = "ios"
	} else if strings.Contains(UserAgent, "Android") {
		device = "android"
	} else {
		device = "others"
	}
	return
}

func Md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func RandomStr() string {
	nano := time.Now().UnixNano()
	seed := rand.New(rand.NewSource(nano))
	salt := seed.Intn(1000)
	sig := fmt.Sprintf("%s%s", nano, salt)
	h := md5.New()
	h.Write([]byte(sig))
	return hex.EncodeToString(h.Sum(nil))
}

func ApiKeyGenerate(apiKey string, timeSlice int) (key string, err error) {
	base := time.Now().Unix() / int64(timeSlice)
	sig := strings.Join([]string{apiKey, strconv.FormatInt(base, 10)}, "")
	h := md5.New()
	h.Write([]byte(sig))
	key = hex.EncodeToString(h.Sum(nil))
	return
}

func ApiKeyCheck(crypted, apiKey string, timeSlice int) bool {
	key, _ := ApiKeyGenerate(apiKey, timeSlice)
	if strings.EqualFold(key, crypted) {
		return true
	}
	return false
}

func HandleError() {
	if x := recover(); x != nil {
		logger := GetLogger()
		logger.Fatal(x)
		fmt.Println(x)
		os.Exit(-1)
	}
}

/*
func DesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	//origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	//根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	//crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func DesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))
	//origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	//origData = ZeroUnPadding(origData)
	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	//去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}
*/
