package interceptors

import (
	"context"
	"log/slog"
	"time"

	"github.com/Ranik23/avito-tech-spring/api/proto/gen/pvz_v1"
	"github.com/Ranik23/avito-tech-spring/internal/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)


func LoggingUnaryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		defer metrics.RequestsTotal.Inc()
		startTime := time.Now()

		md, _ := metadata.FromIncomingContext(ctx)

		resp, err := handler(ctx, req)

		attrs := []slog.Attr{
			slog.String("method", info.FullMethod),
			slog.Duration("duration", time.Since(startTime)),
			slog.Any("headers", md),
		}

		st, _ := status.FromError(err)
		attrs = append(attrs,
			slog.String("status_code", st.Code().String()),
			slog.String("status", st.Code().String()),
		)
		if err != nil {
			attrs = append(attrs, slog.Any("error", err))
			logger.LogAttrs(ctx, slog.LevelError, "gRPC call failed", attrs...)
			return nil, err // нельзя привести же ошибку к dto, есть еще вариант сделать это на уровне выше, то есть на самом HTTP
			// сервере, но там нет интерсепторов 
		}

		logger.LogAttrs(ctx, slog.LevelInfo, "gRPC call completed", attrs...)

		return resp.(*pvz_v1.GetPVZListResponse), nil 
		
	}
}