package middleware

import (
	"context"
	"io"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// StreamingLoggerInterceptor создает интерцептор для логирования стриминговых RPC
func StreamingLoggerInterceptor(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		wrappedStream := &loggingServerStream{
			ServerStream: ss,
			logger:       logger,
			method:       info.FullMethod,
		}

		return handler(srv, wrappedStream)
	}
}

// loggingServerStream обертка для ServerStream с логированием
type loggingServerStream struct {
	grpc.ServerStream
	logger *slog.Logger
	method string
}

// RecvMsg переопределяет RecvMsg для логирования входящих сообщений
func (l *loggingServerStream) RecvMsg(m interface{}) error {
	err := l.ServerStream.RecvMsg(m)
	if err != nil && err != io.EOF {
		l.logger.Error("stream recv error",
			"method", l.method,
			"error", err,
		)
	} else if err == nil {
		l.logger.Info("stream message received",
			"method", l.method,
			"message_type", getMessageType(m),
		)
	}
	return err
}

// SendMsg переопределяет SendMsg для логирования исходящих сообщений
func (l *loggingServerStream) SendMsg(m interface{}) error {
	err := l.ServerStream.SendMsg(m)
	if err != nil {
		l.logger.Error("stream send error",
			"method", l.method,
			"error", err,
		)
	} else {
		l.logger.Info("stream message sent",
			"method", l.method,
			"message_type", getMessageType(m),
		)
	}
	return err
}

// Context возвращает контекст стрима
func (l *loggingServerStream) Context() context.Context {
	return l.ServerStream.Context()
}

// SetHeader устанавливает заголовки
func (l *loggingServerStream) SetHeader(md metadata.MD) error {
	return l.ServerStream.SetHeader(md)
}

// SendHeader отправляет заголовки
func (l *loggingServerStream) SendHeader(md metadata.MD) error {
	return l.ServerStream.SendHeader(md)
}

// SetTrailer устанавливает трейлеры
func (l *loggingServerStream) SetTrailer(md metadata.MD) {
	l.ServerStream.SetTrailer(md)
}

// getMessageType извлекает тип сообщения для логирования
func getMessageType(m interface{}) string {
	if m == nil {
		return "nil"
	}
	return "unknown"
}
