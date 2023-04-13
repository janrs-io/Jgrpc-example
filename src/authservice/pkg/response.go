package pkg

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

const (
	ReturnSuccessCode      int32 = 200 // Operation success code
	ReturnErrorCode        int32 = 500 // System error identifier code appears for operation
	ReturnFailCode         int32 = 400 // Operation failure identifier
	ReturnUnauthorizedCode int32 = 401 // Operation Not Authenticated Error Identifier
	ReturnForbiddenCode    int32 = 403 // Operation unauthorized error identification code
	ReturnNotFoundCode     int32 = 404 // The request address does not have an identifier

	ReturnSuccessMessage      string = "操作成功"   // 200
	ReturnErrorMessage        string = "系统错误"   // 500
	ReturnFailMessage         string = "操作失败"   // 400
	ReturnUnauthorizedMessage string = "操作未授权"  // 401
	ReturnForbiddenMessage    string = "操作没有权限" // 403
	ReturnNotFoundMessage     string = "未知操作"   // 404
)

// Response Http error response
type Response struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// CustomMarshaler Custom marshaler
type CustomMarshaler struct{ runtime.JSONPb }

// Marshal Custom marshal
func (m *CustomMarshaler) Marshal(v interface{}) ([]byte, error) { return nil, nil }

// HttpErrorHandler Http service error handler
func HttpErrorHandler(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, req *http.Request, err error) {

	r := &Response{}
	pb := status.Convert(err).Proto()
	pbMessage := pb.GetMessage()
	pbCode := pb.GetCode()
	pbDetail := pb.GetDetails()

	switch pbCode {
	case int32(codes.OK):
		r.Success(w, ReturnSuccessMessage, nil)
	case int32(codes.NotFound):
		r.NotFound(w, pbMessage, pbDetail)
	case int32(codes.PermissionDenied):
		r.Forbidden(w, pbMessage, pbDetail)
	case int32(codes.Unauthenticated):
		r.Unauthorized(w, pbMessage, pbDetail)
	case int32(codes.FailedPrecondition):
		r.Fail(w, pbMessage, pbDetail)
	case int32(codes.InvalidArgument):
		r.Fail(w, pbMessage, pbDetail)
	default:
		r.Error(w, pbMessage, pbDetail)
	}

}

// HttpSuccessResponseModifier Request successful return data format
func HttpSuccessResponseModifier(ctx context.Context, w http.ResponseWriter, pbMsg proto.Message) error {

	r := &Response{}
	r.Success(w, ReturnSuccessMessage, pbMsg)
	return nil

}

// Response Uniform method for handling returned http api requests.
// The uniform http status returned is 200, in json format, and other errors are specified by the parameter Code.
func (r *Response) Response(w http.ResponseWriter, response *Response) {

	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(response)
	if err != nil {
		r.responseMarshalError(w, response)
		return
	}
	r.writeWithLog(w.Write, b)

}

// responseMarshalError Handling json escape failures
func (*Response) responseMarshalError(w http.ResponseWriter, r *Response) {

	w.WriteHeader(http.StatusOK)
	b := []byte("{\"code\":" + strconv.Itoa(int(ReturnErrorCode)) + ",\"message\":\"" + ReturnErrorMessage + "\",\"data\":{}}")
	r.writeWithLog(w.Write, b)

}

// writeWithLog Logging write errors
func (*Response) writeWithLog(write func([]byte) (int, error), body []byte) {

	_, err := write(body)
	if err != nil {
		log.Printf("http response write failed: %v", err)
	}

}

// Success Operation success method
func (r *Response) Success(w http.ResponseWriter, msg string, data any) {

	r.Response(w, &Response{
		Code:    ReturnSuccessCode,
		Message: msg,
		Data:    data,
	})

}

// NotFound 404 NotFound
func (r *Response) NotFound(w http.ResponseWriter, msg string, data any) {

	r.Response(w, &Response{
		Code:    ReturnNotFoundCode,
		Message: msg,
		Data:    data,
	})

}

// Forbidden Operation no permission error
func (r *Response) Forbidden(w http.ResponseWriter, msg string, data any) {

	r.Response(w, &Response{
		Code:    ReturnForbiddenCode,
		Message: msg,
		Data:    data,
	})

}

// Unauthorized Operation not authenticated error
func (r *Response) Unauthorized(w http.ResponseWriter, msg string, data any) {

	r.Response(w, &Response{
		Code:    ReturnUnauthorizedCode,
		Message: msg,
		Data:    data,
	})

}

// ParamError Parameter error method
func (r *Response) ParamError(w http.ResponseWriter, msg string, data any) {

	r.Response(w, &Response{
		Code:    ReturnFailCode,
		Message: "参数错误",
		Data:    data,
	})

}

// Error Unknown error method
func (r *Response) Error(w http.ResponseWriter, msg string, data any) {

	r.Response(w, &Response{
		Code:    ReturnErrorCode,
		Message: msg,
		Data:    data,
	})

}

// Fail Operation failure method
func (r *Response) Fail(w http.ResponseWriter, msg string, data any) {

	r.Response(w, &Response{
		Code:    ReturnFailCode,
		Message: msg,
		Data:    data,
	})

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

// MarshalJSON Marshal to json response format
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
