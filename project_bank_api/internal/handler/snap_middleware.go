package handler

import (
	"belajar-go/project_bank_api/config"
	"belajar-go/project_bank_api/internal/dto"
	"belajar-go/project_bank_api/pkg"
	"belajar-go/project_bank_api/pkg/telemetry"
	"encoding/json"
	"net/http"
	"os"

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
		attribute.String("snap.external_id", externalID),
		attribute.String("snap.channel_id", channelID),
		attribute.Bool("snap.has_authorization", authorization != ""),
		attribute.Bool("snap.has_signature", signature != ""),
		attribute.String("snap.ip_address", r.Header.Get("X-IP-ADDRESS")),
		attribute.String("snap.device_id", r.Header.Get("X-DEVICE-ID")),
	)

	if timestamp == "" {
		span.SetStatus(codes.Error, "missing X-TIMESTAMP")
		config.Log.Warn("Missing X-TIMESTAMP header")
		snapErr := dto.SnapInvalidMandatoryField.WithField("X-TIMESTAMP").ToResponse("00")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	if partnerID == "" {
		span.SetStatus(codes.Error, "missing X-PARTNER-ID")
		config.Log.Warn("Missing X-PARTNER-ID header")
		snapErr := dto.SnapInvalidMandatoryField.WithField("X-PARTNER-ID").ToResponse("00")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	if externalID == "" {
		span.SetStatus(codes.Error, "missing X-EXTERNAL-ID")
		config.Log.Warn("Missing X-EXTERNAL-ID header")
		snapErr := dto.SnapInvalidMandatoryField.WithField("X-EXTERNAL-ID").ToResponse("00")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	if channelID == "" {
		span.SetStatus(codes.Error, "missing CHANNEL-ID")
		config.Log.Warn("Missing CHANNEL-ID header")
		snapErr := dto.SnapInvalidMandatoryField.WithField("CHANNEL-ID").ToResponse("00")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	if authorization == "" {
		span.SetStatus(codes.Error, "missing Authorization")
		config.Log.Warn("Missing Authorization header")
		snapErr := dto.SnapInvalidTokenB2B.ToResponse("00")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	if signature != "" {
		secretKey := os.Getenv("SNAP_SECRET_KEY")
		if !pkg.VerifySignature(secretKey, timestamp+"|"+partnerID, signature) {
			span.SetStatus(codes.Error, "invalid signature")
			config.Log.Warn("Invalid X-SIGNATURE")
			snapErr := dto.SnapUnauthorized.WithReason("Invalid signature").ToResponse("00")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(snapErr)
			return
		}
		span.SetAttributes(attribute.Bool("snap.signature_valid", true))
	}

	span.SetStatus(codes.Ok, "headers validated")

	config.Log.Info("SNAP request",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("partner_id", partnerID),
		zap.String("external_id", externalID),
	)

	r = r.WithContext(ctx)
	m.next.ServeHTTP(w, r)
}
