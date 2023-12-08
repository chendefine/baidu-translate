package baidu_translate

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestBaiduTranslate(t *testing.T) {
	cfg := new(BaiduTranslateConfig)
	raw, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(raw, cfg)
	if err != nil {
		panic(err)
	}
	client := NewBaiduTranslateClient(cfg)
	result, err := client.Translate("你好", "zh", "en")
	if err != nil {
		fmt.Println(err)
	} else {
		for _, item := range result {
			fmt.Println(item.Src, item.Dst)
		}
	}
}
