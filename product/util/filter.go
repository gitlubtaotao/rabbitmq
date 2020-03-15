package util

import "net/http"

type FilterHandler func(rw http.ResponseWriter, r *http.Request) error
type WebHandler func(w http.ResponseWriter, r *http.Request)

type Filter struct {
	FilterMap map[string]FilterHandler
}

//初始化Filter
func NewFilter() (filter *Filter) {
	return &Filter{FilterMap: make(map[string]FilterHandler)}
}

func (filter *Filter) RegisterFilterUri(uri string, handler FilterHandler) {
	filter.FilterMap[uri] = handler
}

func (filter *Filter) GetFilterHandler(uri string) FilterHandler {
	return filter.FilterMap[uri]
}

//执行拦截器，返回函数类型
func (filter *Filter) Handler(webHandler WebHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for path, handler := range filter.FilterMap {
			if path == r.RequestURI {
				if err := handler(w, r); err != nil {
					_, _ = w.Write([]byte(err.Error()))
					return
				}
				break
			}
		}
		webHandler(w, r)
	}
}
