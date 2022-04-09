package rest

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type None struct{}

type Trim[T any] struct{ Value T }

func (t *Trim[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal([]byte(strings.Trim(string(data), `"`)), &t.Value)
}

type JsonRequest[B any, P any, Q any] struct {
	Body     B
	Params   P
	Query    Q
	Request  *http.Request
	Response http.ResponseWriter
}

func (r JsonRequest[B, P, Q]) URL() url.URL {
	u := *r.Request.URL
	u.Host = r.Request.Host
	if r.Request.TLS == nil {
		u.Scheme = "http"
	} else {
		u.Scheme = "https"
	}
	return u
}

type JsonHandler[B any, P any, Q any] func(r JsonRequest[B, P, Q]) (any, int)

func JsonHandle(handle func(w http.ResponseWriter, r *http.Request) (any, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, status := handle(w, r)

		if status != -1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)

			enc := json.NewEncoder(w)
			enc.Encode(result)
		}
	}
}

var (
	JsonParamsParseCode       = "ERR_PARSE_PARAMS"
	JsonParamsValidateCode    = "ERR_VALIDATE_PARAMS"
	JsonParamsValidateMessage = "Invalid params"
	JsonQueryParseCode        = "ERR_PARSE_QUERY"
	JsonQueryValidateCode     = "ERR_VALIDATE_QUERY"
	JsonQueryValidateMessage  = "Invalid query"
	JsonBodyParseCode         = "ERR_PARSE_BODY"
	JsonBodyValidateCode      = "ERR_VALIDATE_BODY"
	JsonBodyValidateMessage   = "Invalid request"
)

func JsonRoute[B any, P any, Q any](handler JsonHandler[B, P, Q]) http.HandlerFunc {
	return JsonHandle(func(w http.ResponseWriter, r *http.Request) (any, int) {
		ctx := ValidateContext{}

		request := JsonRequest[B, P, Q]{
			Request:  r,
			Response: w,
		}

		paramsMap := ParamsToMap(r)
		params, paramsError := ParseMap[P](paramsMap)
		if paramsError != nil {
			return JsonErrorFromParse(paramsError, JsonParamsParseCode), 400
		}
		paramsValidations := GetValidations(params, ctx)
		if len(paramsValidations) > 0 {
			return JsonErrorFromValidations(JsonParamsValidateMessage, JsonParamsValidateCode, paramsValidations), 400
		}
		request.Params = params

		queryMap := QueryToMap(r)
		query, queryError := ParseMap[Q](queryMap)
		if queryError != nil {
			return JsonErrorFromParse(queryError, JsonQueryParseCode), 400
		}
		queryValidations := GetValidations(query, ctx)
		if len(queryValidations) > 0 {
			return JsonErrorFromValidations(JsonQueryValidateMessage, JsonQueryValidateCode, queryValidations), 400
		}
		request.Query = query

		body, bodyError := ParseBody[B](r)
		if bodyError != nil {
			return JsonErrorFromParse(bodyError, JsonBodyParseCode), 400
		}
		bodyValidations := GetValidations(body, ctx)
		if len(bodyValidations) > 0 {
			return JsonErrorFromValidations(JsonBodyValidateMessage, JsonBodyValidateCode, bodyValidations), 400
		}
		request.Body = body

		return handler(request)
	})
}

type JsonError struct {
	Message     string       `json:"message"`
	Code        string       `json:"code"`
	Validations []Validation `json:"validations,omitempty"`
}

func JsonErrorFromParse(e error, code string) JsonError {
	return JsonError{
		Message: e.Error(),
		Code:    code,
	}
}

func JsonErrorFromValidations(message string, code string, validations []Validation) JsonError {
	return JsonError{
		Message:     message,
		Code:        code,
		Validations: validations,
	}
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
	if err.Error() == "EOF" {
		var p any = parsed
		if _, ok := p.(None); ok {
			return parsed, nil
		}
	}
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
			curr = curr.get(node)
		}
		curr.set(v[0])
	}

	return out.convert()
}

type queryNode struct {
	Obj   map[string]*queryNode
	Arr   []*queryNode
	Value any
	Kind  int
}

func (node *queryNode) get(x string) *queryNode {
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
func (node *queryNode) set(value any) {
	node.Value = value
	node.Kind = 3
}
func (node *queryNode) convert() any {
	switch node.Kind {
	case 1:
		c := make([]any, len(node.Arr))
		for i, item := range node.Arr {
			if item != nil {
				c[i] = item.convert()
			} else {
				c[i] = nil
			}
		}
		return c
	case 2:
		c := make(map[string]any)
		for key, value := range node.Obj {
			c[key] = value.convert()
		}
		return c
	}
	return node.Value
}