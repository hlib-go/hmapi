package hmapi

import "testing"

func TestBizEncodeV1(t *testing.T) {
	v, err := BizEncode("12345678", "111")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(v)
}

func TestBizDecodeV1(t *testing.T) {
	v, err := BizDecode("12345678", "7e6de99a0cbf0edce7007eb0a0ac715a0Ye9VQqb91g=")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(v)
}

func TestBizEncodeV2(t *testing.T) {
	v, err := BizEncode("10000212", "111")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(v)
}

func TestBizDecodeV2(t *testing.T) {
	v, err := BizDecode("10000212", "=")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(v)
}

func TestBizEncodeV3(t *testing.T) {
	v, err := BizEncode("10000311", "111")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(v)
}

func TestBizEncodeTV3(t *testing.T) {
	biz := "`123"
	v, err := BizEncodeT("10000311", biz)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(v)
}

func TestBizDecodeV3(t *testing.T) {
	biz := "KpaB5dRtjyyJvnlMZmudug==:31dc75991163b53896f3b1caf43732737546f2a3:eFFOPF09UIMLAU7Onuy8Q/VHH8gAcnllAQvoLjHKOit0lHUlcLkr4Qmc9n+SLbWtjFReZ5eTlfuflGzfqoKACjEEqYM8upDxjmaZMl6xP7lOgPm9fBcXvKJASPOGbZt00hJEcgp42vtoBN9CggI1Gmv/nBOSsZM3D8gMn0axW2U="
	v, err := BizDecode("10000311", biz)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(v)
}
