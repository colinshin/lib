package stringL

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMd5(str string) string {
	md5New := md5.New()
	data := []byte(str)
	md5New.Write(data)
	return hex.EncodeToString(md5New.Sum(nil))
}
