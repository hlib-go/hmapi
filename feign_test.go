package hmapi

import (
	"testing"
)

func TestPost(t *testing.T) {
	b, e := Request(nil, "POST", "http://himkt.cn", "")
	if e != nil {
		t.Error(e)
		return
	}
	t.Log(string(b))
}
