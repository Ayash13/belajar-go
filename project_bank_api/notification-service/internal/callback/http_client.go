package callback

import (
	"bytes"
	"context"
	"net/http"
	"notification-service/pkg/telemetry"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type HTTPClient interface {
	SendCallback(ctx context.Context, payload []byte) (int, error)
}

type httpClient struct {
	client *http.Client
	url    string
}

func NewHTTPClient(url string) HTTPClient {
	return &httpClient{
		client: &http.Client{
			Timeout:   5 * time.Second,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
		url: url,
	}
}

func (c *httpClient) SendCallback(ctx context.Context, payload []byte) (int, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "callback.SendCallback")
	defer span.End()

	span.SetAttributes(
		attribute.String("http.url", c.url),
		attribute.String("http.method", http.MethodPost),
		attribute.Int("http.request_content_length", len(payload)),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewBuffer(payload))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create request")
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to send request")
		return 0, err
	}
	defer resp.Body.Close()

	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
	if resp.StatusCode >= 400 {
		span.SetStatus(codes.Error, "callback failed")
	} else {
		span.SetStatus(codes.Ok, "callback success")
	}

	return resp.StatusCode, nil
}
