package hmapi

import "testing"

func TestDesMd5(t *testing.T) {
	appid := "12345678"
	value, err := EnDesMd5(appid, "111")
	if err != nil {
		return
	}
	t.Log(value)
	v, err := DeDesMd5(appid, value)
	if err != nil {
		return
	}
	t.Log(v)
}
