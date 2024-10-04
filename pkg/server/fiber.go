package server

type FiberConfig struct {
	AppName      string
	HttpPort     string
	AllowOrigins string
	AllowMethods string
	AllowHeaders string
}

//type FiberServer struct {
//	conf *FiberConfig
//	App  *fiber.App
//}
//
//func NewFiberServer(
//	conf *FiberConfig,
//) *FiberServer {
//	// Init fiber app
//	if conf.AppName == "" {
//		conf.AppName = "goFiber"
//	}
//	app := fiber.New(fiber.Config{
//		ServerHeader: conf.AppName,
//		// ErrorHandler:      ErrorHandler,
//		AppName:           conf.AppName,
//		EnablePrintRoutes: false,
//	})
//
//	// Register global middleware
//	app.Use(recover.New(recover.Config{
//		EnableStackTrace: true,
//	}))
//	app.Use(cors.New(cors.Config{
//		AllowOrigins: conf.AllowOrigins,
//		AllowMethods: conf.AllowMethods,
//		AllowHeaders: conf.AllowHeaders,
//	}))
//	app.Use(slogfiber.NewWithConfig(slog.Default(), slogfiber.Config{
//		DefaultLevel:     slog.LevelDebug,
//		ClientErrorLevel: slog.LevelInfo,
//		ServerErrorLevel: slog.LevelError,
//		WithUserAgent:    true,
//		WithRequestID:    true,
//	}))
//	// Monitoring endpoint
//	app.Get("/metrics", monitor.New(monitor.Config{
//		Title: conf.AppName,
//	}))
//	// Health check
//	app.Get("/ping", func(c *fiber.Ctx) error {
//		return c.Status(200).SendString("pong")
//	})
//
//	return &FiberServer{
//		conf: conf,
//		App:  app,
//	}
//}
//
//func (s *FiberServer) Start() error {
//	slog.Info("starting http server on port :" + s.conf.HttpPort)
//	return s.App.Listen(":" + s.conf.HttpPort)
//}
//
//func (s *FiberServer) Shutdown() {
//	slog.Info("shutting down http server")
//	if err := s.App.Shutdown(); err != nil {
//		slog.Error("error shutting down http server", "error", err.Error())
//	}
//}

//
// func ErrorHandler(c *fiber.Ctx, err error) error {
// 	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
// 	var e *http2.ErrorResponse
// 	switch {
// 	case errors.As(err, &e):
// 		if e.StatusCode == http.StatusInternalServerError {
// 			slog.Error("internal server error", "error", err.Error())
// 			return c.Status(e.StatusCode).JSON(http2.ErrorResponse{
// 				StatusCode:  e.StatusCode,
// 				Message:     e.Message,
// 				DetailError: e.DetailError,
// 			})
// 		} else {
// 			return c.Status(e.StatusCode).JSON(ExceptionResponse{
// 				StatusCode: e.StatusCode,
// 				Message:    e.Message,
// 			})
// 		}
// 	default:
// 		fiberErr := new(fiber.Error)
// 		errors.As(err, &fiberErr)
// 		switch fiberErr.Code {
// 		case http.StatusInternalServerError, 0:
// 			slog.Error("internal server error", "error", err.Error())
// 			return c.Status(fiberErr.Code).JSON(http2.ErrorResponse{
// 				StatusCode:  http.StatusInternalServerError,
// 				Message:     "internal server error",
// 				DetailError: err.Error(),
// 			})
// 		default:
// 			return c.Status(fiberErr.Code).JSON(ExceptionResponse{
// 				StatusCode: fiberErr.Code,
// 				Message:    fiberErr.Message,
// 			})
// 		}
// 	}
// }
