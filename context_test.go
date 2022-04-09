package boot4go

import (
	"fmt"
	"github.com/gohutool/log4go"
	"reflect"
	"testing"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : context_test.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/7 10:01
* 修改历史 : 1. [2022/4/7 10:01] 创建文件 by NST
*/

var logger log4go.Logger

func init() {
	log4go.LoggerManager.InitWithDefaultConfig()

	Context.RegistryBeanInstance("aaa", Hello{})

	Context.RegistryBeanInstance("boot4go.IHello", Hello{})
	/*
		h := &Hello{}
		Test{hello2: h.(IHello)}*/

	logger = log4go.LoggerManager.GetLogger("boot4go.context.test")
}

type Test struct {
	age     int16          `bootable:"${metadata.major}"`
	name    string         `bootable:"${metadata.name}"`
	version string         `bootable:"${metadata.version}"`
	hello   IHello         `bootable:"aaa"`
	hello2  IHello         `bootable`
	data    map[string]any `bootable:"${spec.runAsUser}"`
	list    []any          `bootable:"${spec.volumes}"`
}

type Hello struct {
}

func (h *Hello) sayHello(t Test) Test {
	return Test{}
}

type IHello interface {
	sayHello(t Test) Test
}

func TestContext(t *testing.T) {
	fmt.Println(log4go.LoggerManager)

	test := &Test{}
	typeof := reflect.TypeOf(test)

	fmt.Println(typeof.String())
	fmt.Println(typeof.Kind().String())

	fmt.Println(type2Str(reflect.TypeOf(test)))
	fmt.Println(type2Str(reflect.TypeOf(*test)))

	h, ok := interface{}(test).(IHello)

	fmt.Println(ok)

	var ih IHello = &Hello{}

	fmt.Println(type2Str(reflect.TypeOf(h)))
	fmt.Println(type2Str(reflect.TypeOf(ih)))

	fmt.Println(IHello.sayHello(ih, *test))

	h.sayHello(*test)
}

func TestContext2(t *testing.T) {
	ty := reflect.TypeOf(Test{})
	count := ty.NumField()
	fmt.Println("Count ", count)

	for idx := 0; idx < count; idx++ {
		a, ok := type2Str(ty.Field(idx).Type)
		fmt.Println(ty.Field(idx).Name, " ", a, " ", ok)
	}
}

func TestGetBean(t *testing.T) {
	bean, ok := Context.GetBean(Test{})

	t1 := bean.(*Test)
	logger.Info(&t1.hello2, "  ", &t1.hello)

	bean, _ = Context.getBeanByName("boot4go.Test")
	t1 = bean.(*Test)
	logger.Info(&t1.hello2, "  ", &t1.hello)

	logger.Info("%v %v", reflect.TypeOf(bean.(*Test)).String(), ok)
	logger.Info(bean)

	logger.Info("%v", &t1.data)
	logger.Info("%v", &t1.list)

	time.Sleep(10 * time.Second)
}

func TestContextConfiguration(t *testing.T) {

	bean, ok := Context.GetBean(&Test{})
	fmt.Println(reflect.TypeOf(bean.(*Test)).String(), bean, ok)

	logger := log4go.LoggerManager.GetLogger("test")

	logger.Info("YAML %v", ConfigurationContext.ToMap())

	logger.Info("YAML %v", ConfigurationContext.GetValue("${metadata.name}"))
	logger.Info("YAML %v", ConfigurationContext.GetValue("${spec.volumes[0]}"))

	time.Sleep(10 * time.Second)
}
