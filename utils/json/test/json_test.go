package json

import (
	"bytes"
	"encoding/json"
	"github.com/bytedance/sonic/decoder"
	json2 "github.com/flyerxp/lib/utils/json"
	"strings"
	"testing"
)

// 普通Encode测试
func TestEncode(t *testing.T) {
	data := map[string]string{"a": "<>asdf中国asdf%$", "b": "dddddddd<a href=\"\">asfdasdfasfd</a>"}
	b, e := json2.EncodeEscape(&data)
	_, _ = json2.Encode(&data)
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
func TestStreamEncode(t *testing.T) {
	var o1 = map[string]interface{}{
		"a": "b",
	}
	var o2 = 1
	var w = bytes.NewBuffer(nil)
	var enc = json2.StreamEncode(w)
	_ = enc.Encode(o1)
	_ = enc.Encode(o2)
	_ = json2.StreamDecode(w)
	t.Log(w.String())
}
func TestStreamDecode(t *testing.T) {
	var o = map[string]interface{}{}
	var r = strings.NewReader(`{"a":"b"}{"1":"2"}`)
	var dec = decoder.NewStreamDecoder(r)
	_ = dec.Decode(&o)
	_ = dec.Decode(&o)
	t.Logf("%+v", o)
}
