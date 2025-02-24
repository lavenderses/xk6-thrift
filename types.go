package thrift

type TTypes struct {}

func (*TTypes) NewTString(v string) TString {
	return NewTstring(v)
}

func (*TTypes) NewTRequest(v *map[int16]TValue) *TRequest {
	return NewTRequestWithValue(v)
}
