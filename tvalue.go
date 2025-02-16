package thrift

type TValue interface {
	Equals(other *TValue) bool
}
