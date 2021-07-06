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

func TestBizDecodeV3(t *testing.T) {
	v, err := BizDecode("12345678", "000")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(v)
}
