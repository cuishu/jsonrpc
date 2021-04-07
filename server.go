package jsonrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

const (
	// ParseError Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text.
	ParseError = -32700
	// InvalidRequest The JSON sent is not a valid Request object.
	InvalidRequest = -32600
	// MethodNotFound The method does not exist / is not available.
	MethodNotFound = -32601
	// InvalidParams Invalid method parameter(s).
	InvalidParams = -32602
	// InternalError Internal JSON-RPC error.
	InternalError = -32603
	// ServerError -32000 to -32099 Reserved for implementation-defined server-errors.
	ServerError = -32000

	// VERSION20 JSON RPC 2.0协议
	VERSION20 = "2.0"
)

// JRPCRequest JSON RPC 输入协议格式
type JRPCRequest struct {
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
}

// JRPCResponse JSON RPC 返回格式
type JRPCResponse struct {
	ID      int         `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
}

// JRPCErrMesg JSON RPC 错误消息
type JRPCErrMesg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// JRPCError JSON RPC 错误信息格式
type JRPCError struct {
	ID      *string     `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
	Error   JRPCErrMesg `json:"error"`
}

type Method struct {
	Param []reflect.Kind
	Func  reflect.Value
}

// Server 远程调用函数注册结构
type Server struct {
	Methods map[string]Method
}

// NewJRPCServer 创建JSON RPC服务实例
func NewJRPCServer() *Server {
	return &Server{Methods: make(map[string]Method)}
}

func (s *Server) Run(addr ...string) error {
	r := gin.Default()
	r.POST("/rpc", func(c *gin.Context) {
		s.HandleRequest(c.Writer, c.Request)
	})
	return r.Run(addr...)
}

func (s *Server) registObject(name string, rv reflect.Value) {
	rt := rv.Type()
	for i := 0; i < rv.NumMethod(); i++ {
		f := rv.Method(i)
		t := rt.Method(i)
		param := make([]reflect.Kind, f.Type().NumIn())
		for i := 0; i < f.Type().NumIn(); i++ {
			param[i] = f.Type().In(i).Kind()
		}
		reflect.Indirect(f).Type().Name()
		s.Methods[name+"."+t.Name] = Method{Param: param, Func: f}
	}
}

func (s *Server) RegistMethod(name string, o interface{}) {
	rv := reflect.ValueOf(o)
	rt := rv.Type()
	if rv.Kind() == reflect.Func {
		param := make([]reflect.Kind, rt.NumIn())
		for i := 0; i < rt.NumIn(); i++ {
			param[i] = rt.In(i).Kind()
		}
		s.Methods[name] = Method{Param: param, Func: rv}
	} else if rv.Kind() == reflect.Struct {
		s.registObject(name, rv)
	} else if rv.Kind() == reflect.Ptr {
		s.registObject(name, rv)
	} else {
		panic("")
	}
	fmt.Println(s.Methods)
}

// ExecJSONRPC 执行RPC函数
func (s *Server) ExecJSONRPC(js []byte) (interface{}, error) {
	var req JRPCRequest
	if err := json.Unmarshal(js, &req); err != nil {
		return &JRPCError{
			ID:      nil,
			JSONRPC: VERSION20,
			Error: JRPCErrMesg{
				Code:    ParseError,
				Message: err.Error(),
			},
		}, nil
	}
	if method, ok := s.Methods[req.Method]; ok {
		var params []reflect.Value = make([]reflect.Value, len(req.Params))
		function := method.Func
		if len(req.Params) != len(method.Param) {
			return nil, errors.New("Param mismatch")
		}
		for i, param := range req.Params {
			arg, err := convert(param, method.Param[i])
			if err != nil {
				return nil, err
			}
			params[i] = reflect.ValueOf(arg)
		}

		res := function.Call(params)
		var interfaceResult []interface{} = make([]interface{}, len(res))
		for i, v := range res {
			interfaceResult[i] = v.Interface()
		}
		return &JRPCResponse{
			ID:      req.ID,
			JSONRPC: VERSION20,
			Result:  interfaceResult,
		}, nil
	}
	return &JRPCError{
		ID:      nil,
		JSONRPC: VERSION20,
		Error: JRPCErrMesg{
			Code:    MethodNotFound,
			Message: "Method Not Found.",
		},
	}, nil
}

// HandleRequest 处理http request
func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		js, _ := json.Marshal(&JRPCError{
			ID:      nil,
			JSONRPC: VERSION20,
			Error: JRPCErrMesg{
				Code:    InvalidRequest,
				Message: err.Error(),
			},
		})
		w.Write(js)
		return
	}
	res, err := s.ExecJSONRPC(body)
	if err != nil {
		js, _ := json.Marshal(&JRPCError{
			ID:      nil,
			JSONRPC: VERSION20,
			Error: JRPCErrMesg{
				Code:    InvalidRequest,
				Message: err.Error(),
			},
		})
		w.Write(js)
		return
	}
	js, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	w.Write(js)
}
