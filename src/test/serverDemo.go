package main

import (
	"fmt"
	"reflect"
)

// 函数名称小写，表示不允许外部包访问
// 函数定义
func add(a int, b int) (int, error){
	return a+b, nil
}

func main() {
	// 通过 map 注册函数，key 为函数名称，value 为函数实体
	funcs := make(map[string]reflect.Value)
	funcs["add"] = reflect.ValueOf(add)

	// 模拟客户端的访问请求（即参数列表）
	request := []reflect.Value{reflect.ValueOf(1), reflect.ValueOf(2)}

	// 通过Call调用，传入参数数组request。计算求得结果
	vals := funcs["add"].Call(request)
	//fmt.Println(vals) // [<int Value> <error Value>]，即此时的vals中的数据类型是Value，还有进一步处理

	// 存放返回结果
	var response []interface{}
	for _, val := range vals {
		response = append(response, val.Interface())
	}

	fmt.Println(response) // [3 <nil>]
}
