package json

import (
	json2 "CH123/Lib/utils/json"
	"encoding/json"
	"testing"
)

// 普通Encode测试
func TestEncode(t *testing.T) {
	data := map[string]string{"a": "<>asdf中国asdf%$", "b": "dddddddd<a href=\"\">asfdasdfasfd</a>"}
	b, e := json2.EncodeEscape(&data)
	b2, _ := json.Marshal(data)
	var data1 map[string]string
	//json2.Decode(b, &data1)
	t.Logf("decode %#v", data1)
	//json2.Encode(v, EscapeHTML)
	t.Logf("%#v=%#v", b, b2)
	if e != nil {
		t.Errorf("%#v", e.Error())
	}
	if string(b) != string(b2) {
		t.Errorf("产生了错误的值 %s = %s", string(b), string(b2))
	}
}
