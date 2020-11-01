package hm

/*
根据结构体生成结构文档页面
*/

type ApiDoc struct {
	Pattern     string
	Description string
	ReqBody     interface{}
	ResBody     interface{}
}

// 接口名、输入参数结构体对象、响应参数结构体对象
func DefDoc(pattern string, i interface{}, o interface{}) {

}
