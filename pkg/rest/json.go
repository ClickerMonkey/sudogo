package rest

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type None struct{}

type JsonHandler[B any, P any, Q any] func(body B, params P, query Q) (any, int)

func JsonHandle(handle func(w http.ResponseWriter, r *http.Request) (any, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, status := handle(w, r)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		json.NewEncoder(w).Encode(result)
	}
}

func JsonRoute[B any, P any, Q any](handler JsonHandler[B, P, Q]) http.HandlerFunc {
	return JsonHandle(func(w http.ResponseWriter, r *http.Request) (any, int) {
		ctx := ValidateContext{}

		body, bodyError := ParseBody[B](r)
		if bodyError != nil {
			return bodyError, 400
		}
		bodyValidations := GetValidations(body, ctx)
		if len(bodyValidations) > 0 {
			return bodyValidations, 400
		}

		paramsMap := ParamsToMap(r)
		params, paramsError := ParseMap[P](paramsMap)
		if paramsError != nil {
			return paramsError, 400
		}
		paramsValidations := GetValidations(params, ctx)
		if len(paramsValidations) > 0 {
			return paramsValidations, 400
		}

		queryMap := QueryToMap(r)
		query, queryError := ParseMap[Q](queryMap)
		if queryError != nil {
			return queryError, 400
		}
		queryValidations := GetValidations(query, ctx)
		if len(queryValidations) > 0 {
			return queryValidations, 400
		}

		return handler(body, params, query)
	})
}

func ParseMap[T any](m any) (T, error) {
	var parsed T
	s := strings.Builder{}
	err := json.NewEncoder(&s).Encode(m)
	if err == nil {
		reader := strings.NewReader(s.String())
		err = json.NewDecoder(reader).Decode(&parsed)
	}
	return parsed, err
}

func ParseBody[T any](r *http.Request) (T, error) {
	var parsed T
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&parsed)
	return parsed, err
}

func ParamsToMap(r *http.Request) any {
	out := make(map[string]string)
	ctx := chi.RouteContext(r.Context())
	if ctx != nil {
		for i, key := range ctx.URLParams.Keys {
			value := ctx.URLParams.Values[i]
			out[key] = value
		}
	}
	return out
}

func QueryToMap(r *http.Request) any {
	out := &queryNode{}
	pathRegex := regexp.MustCompile(`[\]\[]+`)
	queryValues := r.URL.Query()

	for k, v := range queryValues {
		if len(v) == 0 {
			continue
		}
		path := pathRegex.Split(strings.TrimRight(k, "]"), -1)
		curr := out
		for _, node := range path {
			curr = curr.Get(node)
		}
		curr.Set(v[0])
	}

	return out.Convert()
}

type queryNode struct {
	Obj   map[string]*queryNode
	Arr   []*queryNode
	Value any
	Kind  int
}

func (node *queryNode) Get(x string) *queryNode {
	if i, err := strconv.Atoi(x); err == nil {
		node.Kind = 1
		if len(node.Arr) <= i {
			arr := make([]*queryNode, i+1)
			copy(arr, node.Arr)
			node.Arr = arr
		}
		n := node.Arr[i]
		if n == nil {
			n = &queryNode{}
			node.Arr[i] = n
		}
		return n
	} else {
		node.Kind = 2
		if node.Obj == nil {
			node.Obj = map[string]*queryNode{}
		}
		n := node.Obj[x]
		if n == nil {
			n = &queryNode{}
			node.Obj[x] = n
		}
		return n
	}
}
func (node *queryNode) Set(value any) {
	node.Value = value
	node.Kind = 3
}
func (node *queryNode) Convert() any {
	switch node.Kind {
	case 1:
		c := make([]any, len(node.Arr))
		for i, item := range node.Arr {
			if item != nil {
				c[i] = item.Convert()
			} else {
				c[i] = nil
			}
		}
		return c
	case 2:
		c := make(map[string]any)
		for key, value := range node.Obj {
			c[key] = value.Convert()
		}
		return c
	}
	return node.Value
}
