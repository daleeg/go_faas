package util

import (
	"fmt"
	"reflect"
)

type Result interface {
	GetCode() error
	GetData() interface{}
	ShowData()
}

type ErrorResult struct {
	err error
}

func (e ErrorResult) GetCode() error {
	return e.err
}
func (e ErrorResult) GetData() interface{} {
	return nil
}

func (e ErrorResult)ShowData() {}

type SuccessResult struct {
	data *[]reflect.Value
	retType *[]reflect.Type
}

func (s SuccessResult) GetCode() error {
	return nil
}
func (s SuccessResult) GetData() interface{} {
	return *s.data
}

func (s SuccessResult) ShowData() {
	ptr_data := s.data
	for i := 0; i < len(*ptr_data); i++ {
		switch (*ptr_data)[i].Type().Kind() {
		case reflect.Int:
			fmt.Println("result: ", i, ", ", (*ptr_data)[i].Interface().(int))
		case reflect.String:
			fmt.Println("result: ", i, ", ", (*ptr_data)[i].Interface().(string))
		default:
			fmt.Printf("type: %s[%s], value: %v \n",  (*ptr_data)[i].Type().Kind(),
				(*ptr_data)[i].Type().Name(),  (*ptr_data)[i].Interface())
		}
	}
}
