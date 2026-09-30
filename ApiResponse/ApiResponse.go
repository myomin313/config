package ApiResponse

import (
	"config-service/pb/commonpb"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/proto"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Datas   interface{} `json:"datas"`
}

type ValidationResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type TypeResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Datas   struct{} `json:"datas"`
}

type PageResult struct {
	Count int         `json:"count"`
	Data  interface{} `json:"data"`
}

func Success(c *fiber.Ctx, code int, data interface{}, message string) error {

	// log.Println("original data", data)

	b, err := json.Marshal(Response{
		Code:    code,
		Message: message,
		Datas:   data,
	})

	if err != nil {
		return err
	}

	resp := &commonpb.Response{
		Code:    int32(code), // Convert int → int32 for protobuf
		Message: message,
		Datas:   b,
	}

	bytes, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	// log.Println("byte data", bytes)

	c.Set("Content-Type", "application/x-protobuf")

	return c.Status(code).Send(bytes)
}

// new Pagination function will use onward from march 16
func Pagination(
	c *fiber.Ctx,
	code int,
	data interface{},
	message string,
) error {

	b, err := json.Marshal(Response{
		Code:    code,
		Message: message,
		Datas:   data,
	})

	if err != nil {
		return err
	}

	resp := &commonpb.Response{
		Code:    int32(code), // Convert int → int32 for protobuf
		Message: message,
		Datas:   b,
	}

	bytes, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	// log.Println("byte data", bytes)

	c.Set("Content-Type", "application/x-protobuf")

	return c.Status(code).Send(bytes)
}

// Success response (WITH pagination)
func SuccessWithPagination(
	c *fiber.Ctx,
	code int,
	count int,
	data interface{},
	message string,
) error {

	b, err := json.Marshal(Response{
		Code:    code,
		Message: message,
		Datas: map[string]interface{}{
			"pageResult": PageResult{
				Count: count,
				Data:  data,
			},
		},
	})

	if err != nil {
		return err
	}

	resp := &commonpb.Response{
		Code:    int32(code), // Convert int → int32 for protobuf
		Message: message,
		Datas:   b,
	}

	bytes, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	// log.Println("byte data", bytes)

	c.Set("Content-Type", "application/x-protobuf")

	return c.Status(code).Send(bytes)
}

// validation message
func ValidationError(c *fiber.Ctx, code int, msg string, errMsg ...string) error {

	insideMsg := msg
	if len(errMsg) > 0 && errMsg[0] != "" {
		insideMsg = errMsg[0]
	}

	b, err := json.Marshal(ValidationResponse{
		Code:    code,
		Message: insideMsg,
	})

	resp := &commonpb.Response{
		Code:    int32(code), // Convert int → int32 for protobuf
		Message: msg,
		Datas:   b,
	}

	bytes, err := proto.Marshal(resp)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "application/x-protobuf")
	return c.Status(code).Send(bytes)
}

// Error response (common)
func Error(c *fiber.Ctx, code int, msg string) error {
	b, err := json.Marshal(Response{
		Code:    code,
		Message: msg,
		Datas:   map[string]interface{}{},
	})

	resp := &commonpb.Response{
		Code:    int32(code), // Convert int → int32 for protobuf
		Message: msg,
		Datas:   b,
	}

	bytes, err := proto.Marshal(resp)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "application/x-protobuf")
	return c.Status(code).Send(bytes)
}

func ListSuccess(c *fiber.Ctx, code int, data interface{}, message string) error {
	return c.JSON(Response{
		Code:    code,
		Message: message,
		Datas:   data, // just the slice of games
	})
}

func JsonSuccess(c *fiber.Ctx, code int, data interface{}, message string) error {
	return c.Status(code).JSON(Response{
		Code:    code,
		Message: message,
		Datas:   data,
	})
}

// Error response (common)
func JsonError(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(Response{
		Code:    code,
		Message: msg,
		Datas:   map[string]interface{}{},
	})
}

func JsonValidationError(c *fiber.Ctx, code int, msg string, errMsg ...string) error {

	insideMsg := msg
	if len(errMsg) > 0 && errMsg[0] != "" {
		insideMsg = errMsg[0]
	}

	return c.Status(code).JSON(ValidationResponse{
		Code:    code,
		Message: insideMsg,
	})
}
