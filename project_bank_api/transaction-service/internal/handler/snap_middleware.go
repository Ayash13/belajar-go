package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"transaction-service/config"
	"transaction-service/internal/dto"
	"transaction-service/pkg"
	"transaction-service/pkg/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type SNAPMiddleware struct {
	next http.Handler
}

func NewSNAPMiddleware(next http.Handler) *SNAPMiddleware {
	return &SNAPMiddleware{next: next}
}

func (m *SNAPMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.URL.Path == "/health" || r.URL.Path == "/metrics" {
		m.next.ServeHTTP(w, r)
		return
	}

	ctx, span := telemetry.Tracer.Start(r.Context(), "middleware.SNAPHeaderValidation")
	defer span.End()

	timestamp := r.Header.Get("X-TIMESTAMP")
	partnerID := r.Header.Get("X-PARTNER-ID")
	externalID := r.Header.Get("X-EXTERNAL-ID")
	channelID := r.Header.Get("CHANNEL-ID")
	authorization := r.Header.Get("Authorization")
	signature := r.Header.Get("X-SIGNATURE")

	span.SetAttributes(
		attribute.String("snap.timestamp", timestamp),
		attribute.String("snap.partner_id", partnerID),
		attribute.Bool("snap.has_authorization", authorization != ""),
	)

	if timestamp == "" {
		span.SetStatus(codes.Error, "missing X-TIMESTAMP")
		snapErr := dto.SnapInvalidMandatoryField.WithField("X-TIMESTAMP").ToResponse("00")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	if partnerID == "" {
		span.SetStatus(codes.Error, "missing X-PARTNER-ID")
		snapErr := dto.SnapInvalidMandatoryField.WithField("X-PARTNER-ID").ToResponse("00")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	if externalID == "" {
		span.SetStatus(codes.Error, "missing X-EXTERNAL-ID")
		snapErr := dto.SnapInvalidMandatoryField.WithField("X-EXTERNAL-ID").ToResponse("00")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	if channelID == "" {
		span.SetStatus(codes.Error, "missing CHANNEL-ID")
		snapErr := dto.SnapInvalidMandatoryField.WithField("CHANNEL-ID").ToResponse("00")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	if authorization == "" {
		snapErr := dto.SnapInvalidTokenB2B.ToResponse("00")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	if signature != "" {
		secretKey := os.Getenv("SNAP_SECRET_KEY")
		if !pkg.VerifySignature(secretKey, timestamp+"|"+partnerID, signature) {
			snapErr := dto.SnapUnauthorized.WithReason("Invalid signature").ToResponse("00")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(snapErr)
			return
		}
	}

	span.SetStatus(codes.Ok, "headers validated")
	config.Log.Info("SNAP request",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("partner_id", partnerID),
	)

	r = r.WithContext(ctx)
	m.next.ServeHTTP(w, r)
}
