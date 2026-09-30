package config

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gopkg.in/natefinch/lumberjack.v2"
)

func AccessLogger() logger.Config {
	logFile := &lumberjack.Logger{
		Filename:   "logs/access.log", // base name
		MaxSize:    100,               // MB (per file)
		MaxBackups: 30,                // keep 30 days
		MaxAge:     30,                // days
		Compress:   true,
	}

	return logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Output:     logFile,
	}
}

var ErrorLog = log.New(&lumberjack.Logger{
	Filename:   "logs/error.log",
	MaxSize:    50,
	MaxBackups: 30,
	MaxAge:     30,
	Compress:   true,
}, "", log.LstdFlags)

// func ErrorHandler(c *fiber.Ctx, err error) error {
// 	ErrorLog.Printf(
// 		"ERROR %d %s %s | %v",
// 		c.Response().StatusCode(),
// 		c.Method(),
// 		c.Path(),
// 		err,
// 	)

// 	b, err := json.Marshal(ApiResponse.Response{
// 		Code:    fiber.StatusInternalServerError,
// 		Message: err.Error(),
// 		Datas:   map[string]interface{}{},
// 	})

// 	resp := &commonpb.Response{
// 		Code:    int32(fiber.StatusInternalServerError), // Convert int → int32 for protobuf
// 		Message: "Internal Server Error",
// 		Datas:   b,
// 	}

// 	bytes, err := proto.Marshal(resp)
// 	if err != nil {
// 		return err
// 	}

// 	c.Set("Content-Type", "application/x-protobuf")
// 	return c.Status(fiber.StatusInternalServerError).Send(bytes)

// }

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it is a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	} else {
		message = err.Error()
	}

	ErrorLog.Printf(
		"ERROR %d %s %s | %v",
		code,
		c.Method(),
		c.Path(),
		err,
	)

	// Return JSON instead of letting nginx show HTML
	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"code":    code,
		"message": message,
	})
}

func Info(msg string) {
	//log to console
	log.Println(msg)

	//log to file
	ErrorLog.Printf(
		"info | %v", msg,
	)

}

func Error(msg string) {
	//log to console
	log.Println(msg)

	//log to file
	ErrorLog.Printf(
		"error | %v", msg,
	)
}

func Warning(msg string) {
	//log to console
	log.Println(msg)

	//log to file
	ErrorLog.Printf(
		"warning | %v", msg,
	)
}
