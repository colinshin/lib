package json

import (
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
	"github.com/bytedance/sonic/encoder"
	"io"
	"sync"
)

type jsonTool struct {
	configDefault sonic.API
	configStd     sonic.API
	Option        sonic.Config
	sync          sync.Once
}

var (
	tool jsonTool
)

func initTool() {
	tool.sync.Do(func() {
		tool.configDefault = sonic.ConfigDefault
		tool.configStd = sonic.Config{
			EscapeHTML:       true,
			SortMapKeys:      false,
			CompactMarshaler: false,
			CopyString:       false,
			ValidateString:   false,
		}.Froze()
	})
}

// 转为json
func Encode(val interface{}) ([]byte, error) {
	initTool()
	return tool.configDefault.Marshal(val)

}

// 如果做转义，性能损耗15%，没必要不要做转义
func EncodeEscape(val interface{}) ([]byte, error) {
	initTool()
	return tool.configStd.Marshal(val)
}

// 转为json
func Decode(data []byte, v interface{}) error {
	initTool()
	return tool.configDefault.Unmarshal(data, v)
}

// 流式encode
func StreamEncode(w io.Writer) *encoder.StreamEncoder {
	initTool()
	return encoder.NewStreamEncoder(w)
}
func StreamDecode(r io.Reader) *decoder.StreamDecoder {
	initTool()
	return decoder.NewStreamDecoder(r)

}
