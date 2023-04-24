package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	protoTypeUrlStr string = "@type"
	protoDataKeyStr string = "data"
	defaultCode     int    = 0
	defaultMsg      string = "操作成功"
)

type defaultData struct{}

// ProtoResponse 接受 gRPC 成功时返回的数据结构体
type ProtoResponse struct {
	Code any `json:"code"`
	Msg  any `json:"msg"`
	Data any `json:"data"`
}

// Response Http 服务返回的结构体
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// CustomMarshaller Custom marshaller
type CustomMarshaller struct{ runtime.JSONPb }

// Marshal Custom marshal
func (m *CustomMarshaller) Marshal(v interface{}) ([]byte, error) { return nil, nil }

// HttpErrorHandler Http service error handler
func HttpErrorHandler(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, req *http.Request, err error) {

	r := &Response{
		Code: HTTPStatusFromCode(status.Convert(err).Code()),
		Msg:  status.Convert(err).Proto().GetMessage(),
	}
	// 判断返回的 details 是否为空，如果是，则设置成空的结构体对象
	if status.Convert(err).Proto().GetDetails() == nil {
		r.Data = &defaultData{}
	} else {
		r.Data = status.Convert(err).Proto().GetDetails()
	}
	// 转译成返回 http 请求的 json 数据格式
	jsonStr, err := json.Marshal(r)
	if err != nil {
		log.Println("系统错误，错误信息：" + err.Error())
		r.Response(http.StatusInternalServerError, w, "系统错误")
		return
	}
	r.Response(r.Code, w, string(jsonStr))

}

// HttpSuccessResponseModifier Request successful return data format
func HttpSuccessResponseModifier(ctx context.Context, w http.ResponseWriter, pbMsg proto.Message) error {

	r := &Response{}
	pr := &ProtoResponse{}

	b, err := protojson.Marshal(pbMsg)
	if err != nil {
		log.Println("系统错误，错误信息：" + err.Error())
		r.Response(http.StatusInternalServerError, w, "系统错误")
		return nil
	}
	// 序列化 json 数据到需要返回的结构体
	err = json.Unmarshal(b, pr)
	if err != nil {
		log.Println("系统错误，错误信息：" + err.Error())
		r.Response(http.StatusInternalServerError, w, "系统错误")
		return nil
	}
	// 如果 grpc 没有传递数据，则为空的结构体对象数据
	if pr.Data == nil {
		pr.Data = &defaultData{}
	} else {
		// 如果 grpc 有传递 Data 数据
		// 过滤 protojson 转译后数据自带的 "@type" 字段
		data, ok := pr.Data.(map[string]interface{})
		if ok {
			pr.Data = r.FilterDataKey(r.FilterTypeUrl(data))
		}
	}
	// 如果 grpc 没有传递状态码，则默认为 0
	if pr.Code == nil {
		pr.Code = defaultCode
	} else {
		// 如果有传递 code 状态码，转成数字类型
		// 只有数字类型的字符串才能成功转译，否则报错
		codeStr, ok := pr.Code.(string)
		if ok {
			if codeInt, err := strconv.Atoi(codeStr); err == nil {
				pr.Code = codeInt
			}
		}
	}

	// 如果 grpc 没有传递 msg 数据，则默认为空字符串
	if pr.Msg == nil {
		pr.Msg = defaultMsg
	}

	// 转译成返回 http 请求的 json 数据格式
	jsonStr, err := json.Marshal(pr)
	if err != nil {
		log.Println("系统错误，错误信息：" + err.Error())
		r.Response(http.StatusInternalServerError, w, "系统错误")
		return nil
	}
	r.Response(http.StatusOK, w, string(jsonStr))
	return nil

}

// FilterTypeUrl 过滤 proto message 的 @type 参数
func (r *Response) FilterTypeUrl(data map[string]any) map[string]any {
	delete(data, protoTypeUrlStr)
	if len(data) > 0 {
		for _, v := range data {
			mapData, ok := v.(map[string]any)
			if ok {
				r.FilterTypeUrl(mapData)
			}
		}
	}
	return data
}

// FilterDataKey 过滤 `data` 参数关键字
func (r *Response) FilterDataKey(data map[string]any) map[string]any {
	for k, v := range data {
		if vData, ok := v.(map[string]any); ok {
			if dataValue, exists := vData[protoDataKeyStr]; exists {
				data[k] = dataValue
			}
		}
	}
	return data
}

// Response Method for handling returned http api requests.
// In json format, and other errors are specified by the parameter Code.
func (r *Response) Response(httpStatus int, w http.ResponseWriter, response string) {
	w.WriteHeader(httpStatus)
	r.write(w.Write, []byte(response))
}

// write Logging write errors
func (*Response) write(write func([]byte) (int, error), body []byte) {

	_, err := write(body)
	if err != nil {
		log.Printf("http response write failed: %v", err)
	}

}

// HTTPStatusFromCode converts a gRPC error code into the corresponding HTTP response status.
// See: https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
func HTTPStatusFromCode(code codes.Code) int {

	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return 499
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		// Note, this deliberately doesn't translate to the similarly named '412 Precondition Failed' HTTP response status.
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	default:
		grpclog.Infof("Unknown gRPC error code: %v", code)
		return http.StatusInternalServerError
	}

}

// MarshalJSON 格式化当请求成功时，只有一个 response 结构体
func (r *Response) MarshalJSON() ([]byte, error) {

	mm := runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: true,
		},
	}
	buf := []byte("{")
	st := reflect.TypeOf(*r)
	vt := reflect.ValueOf(*r)
	for i := 0; i < st.NumField(); i++ {
		if i != 0 {
			buf = append(buf, ',')
		}
		field := st.Field(i)
		tag := field.Tag.Get("json")
		buf = append(buf, []byte("\""+tag+"\": ")...)
		value := vt.Field(i).Interface()
		var vBuf []byte
		if tag == "data" {
			vBuf, _ = mm.Marshal(r.Data)
		} else {
			vBuf, _ = json.Marshal(value)
		}
		buf = append(buf, vBuf...)
	}
	end := []byte{'}'}
	buf = append(buf, end...)
	return buf, nil

}
