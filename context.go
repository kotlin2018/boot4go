package boot4go

import (
	"errors"
	"fmt"
	"github.com/gohutool/log4go"
	"reflect"
	"sync"
	"unsafe"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : context.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/7 09:52
* 修改历史 : 1. [2022/4/7 09:52] 创建文件 by NST
*/

var contextLogger = log4go.LoggerManager.GetLogger("boot4go.context")

func init() {
	configFile := "boot4go.yaml"

	if ConfigurationContext.IsConfigFileExist(configFile) {
		ConfigurationContext.LoadYaml(configFile)
	}

	contextLogger.Debug("Yaml %v", ConfigurationContext.ToMap())
}

var Context = context{processing: make(map[string]any), pooled: make(map[string]any)}

type context struct {
	lock       sync.RWMutex
	processing map[string]any
	pooled     map[string]any
}

func (c *context) RegistryBean(name string, t reflect.Type) {
	c.lock.RLock()
	defer c.lock.RUnlock()
}

func (c *context) registryBeanInstance(name string, i any) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	t := reflect.TypeOf(i)

	if t.Kind() == reflect.Struct || t.Kind() == reflect.Interface {
		v := reflect.New(t).Interface()
		c.pooled[name] = v
	} else if t.Kind() == reflect.Ptr {
		c.pooled[name] = i
	}

}

func (c *context) getBeanByName(name string) (any, bool) {
	v, ok := c.pooled[name]

	return v, ok
}

func (c *context) getBean(obj any) (any, error) {
	t := reflect.TypeOf(obj)
	beanName, _ := type2Str(t)
	//
	//c.lock.RLock()
	//defer c.lock.RUnlock()

	bean, ok := c.resolveBean(beanName)

	if ok == nil {
		return bean, nil
	}

	var newValue reflect.Value

	// get new Object pointer of source Type
	if t.Kind() == reflect.Struct {
		newValue = reflect.New(t)
	} else {
		newValue = reflect.New(t.Elem())
	}

	// get the Type of struct
	t = reflect.TypeOf(newValue.Elem().Interface())
	fmt.Println(t.NumField())
	count := t.NumField()
	for idx := 0; idx < count; idx++ {
		f := t.Field(idx)
		newFieldValue := newValue.Elem().FieldByName(f.Name)

		var v any
		k := f.Type.Kind()

		fmt.Println(k, "\t\t", f.Type.String())

		if k != reflect.Interface && k != reflect.Struct && k != reflect.Ptr {
			if tag := f.Tag.Get("bootable"); len(tag) > 0 {
				v = ConfigurationContext.GetValue(tag)
				if v != nil {
					s := fmt.Sprintf("%v", v)
					v, _ = str2Object(s, k)
				}
			}
		} else {
			bn := f.Tag.Get("bootable")
			if len(bn) == 0 {
				if k == reflect.Ptr {
					newFieldValue = reflect.New(f.Type.Elem())
					bn = newFieldValue.Elem().Type().String()
				} else {
					bn = newFieldValue.Type().String()
				}
			}

			if b, _ := c.pooled[bn]; b != nil {
				v = b
			}
		}

		if v != nil {
			if k == reflect.Ptr {
				reflect.NewAt(newFieldValue.Type(), unsafe.Pointer(newFieldValue.UnsafeAddr())).Elem().Set(reflect.ValueOf(v))
			} else {
				reflect.NewAt(newFieldValue.Type(), unsafe.Pointer(newFieldValue.UnsafeAddr())).Elem().Set(reflect.ValueOf(v))
			}
		}

	}

	//
	//v := reflect.ValueOf(bean).Elem()
	//
	//count := t.NumField()
	//for idx := 0; idx < count; idx++ {
	//	fv := v.Field(idx)
	//	ft := fv.Type()
	//	f := t.Field(idx)
	//	fmt.Println(f.Name, "==========", ft.Kind())
	//
	//	if f.IsExported() {
	//
	//		if fv.Kind() == reflect.Int {
	//			fv.SetInt(20)
	//		} else if fv.Kind() == reflect.String {
	//			fv.SetString("DavidLiu")
	//		}
	//	}
	//
	//	//fV := reflect.ValueOf(bean).Elem().Field(idx)
	//	//a, ok := type2Str(fv.Elem().Type())
	//	//fmt.Println(f.Name, " ", a, " ", ok)
	//}

	c.registryBeanInstance(beanName, newValue.Interface())

	return newValue.Interface(), nil
	/*
		if t.Kind() == reflect.Struct {
			return bean, nil
		} else {
			return bean, nil
		}
	*/
}

func (c *context) resolveBean(name string) (any, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	v, ok := c.pooled[name]
	if ok {
		return v, nil
	}

	return nil, errors.New("Not found \"" + name + "\" bean")
}

// utils
func type2Str(t reflect.Type) (string, error) {
	if t.Kind() == reflect.Struct || t.Kind() == reflect.Interface {
		return t.String(), nil
	} else if t.Kind() == reflect.Ptr {
		return t.Elem().String(), nil
	}

	return t.String(), errors.New(t.String() + " is not struct or interface")
}