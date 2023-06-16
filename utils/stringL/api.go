package stringL

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func GetMd5(str string) string {
	md5New := md5.New()
	data := []byte(str)
	md5New.Write(data)
	fmt.Println(str)
	v := hex.EncodeToString(md5New.Sum(nil))
	fmt.Println(v)
	return v
}
