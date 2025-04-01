package cache

type myCache interface {
	Set(string, interface{})
	Get(string) (interface{}, bool)
}
