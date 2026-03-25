package assemblyai

import (
	"ai-tutor-backend/internal/config"
	"ai-tutor-backend/internal/log"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type AssemblyAiClient struct {
	config  config.AssemblyAiConfig
	httpClient *http.Client
	logger log.Logger
}

func NewAssemblyAiClient(config config.AssemblyAiConfig, httpClient *http.Client, logger log.Logger) *AssemblyAiClient {
	return &AssemblyAiClient{
		config: config,
		httpClient: httpClient,
		logger: logger,
	}
}

func (c *AssemblyAiClient) UploadAudio(ctx context.Context,fileBytes []byte) (string, error) {
	c.logger.Debug("entering assembly ai client upload audio")
	req, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseUrl + "/upload", bytes.NewReader(fileBytes))

	if err != nil {
		return "",fmt.Errorf("assembly ai client upload audio error: %w", err)
	}
	req.Header.Set("authorization", c.config.ApiKey)
	// application/octet-stream means The server treats the body as bytes only
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("assembly ai client upload audio do request error: %w", err)
	}

	defer resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
        b, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("assembly ai upload returned status %d: %s", resp.StatusCode, string(b))
    }

	var out struct {
		UploadUrl string `json:"upload_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("assembly ai client upload audio decode response error: %w", err)
	}

	c.logger.Debug("upload audio result", zap.Any("result", out))

	return out.UploadUrl, nil
}

type TranscribeAudioRequest struct {
	AudioUrl string `json:"audio_url"`
}

func (c *AssemblyAiClient) TranscribeAudio(ctx context.Context, uploadedFileUrl string) (string, error) {
	c.logger.Debug("entering assembly ai client transcribe audio")
	b, err := json.Marshal(TranscribeAudioRequest{
		AudioUrl: uploadedFileUrl,
	})
	if err != nil {
		return "", fmt.Errorf("assembly ai client transcribe audio marshal request error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseUrl + "/transcript", bytes.NewReader(b))

	if err != nil {
		return "",fmt.Errorf("assembly ai client transcribe audio error: %w", err)
	}
	req.Header.Set("authorization", c.config.ApiKey)
	// application/octet-stream means The server treats the body as bytes only
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("assembly ai client transcribe audio do request error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
        b, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("assembly ai transcribe returned status %d: %s", resp.StatusCode, string(b))
    }

	var out struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("assembly ai client transcribe audio decode response error: %w", err)
	}

	for {
		req, err := http.NewRequestWithContext(ctx, "GET", c.config.BaseUrl + "/transcript/" + out.ID, nil)
		if err != nil {
			return "", fmt.Errorf("assembly ai client transcribe audio decode response error: %w", err)
		}

		req.Header.Set("authorization", c.config.ApiKey)

		resp, err := c.httpClient.Do(req)

		if err != nil {
			return "", fmt.Errorf("assembly ai client transcribe audio decode response error: %w", err)
		}

		defer resp.Body.Close()

		
		var outResp struct {
			ID string `json:"id"`
			Text string `json:"text"`
			Status string `json:"status"`
			Error string `json:"error"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&outResp); err != nil {
			return "", fmt.Errorf("assembly ai client upload audio decode response error: %w", err)
		}

		c.logger.Debug("transcribe result", zap.Any("result", outResp))

		switch outResp.Status {
			case "completed":
				c.logger.Debug("assembly ai client transcribe audio result: " + outResp.Text)

				return outResp.Text, nil
			case "error":
				return "", fmt.Errorf("assembly ai client response error: %s", outResp.Error)
				
		}




	}

	
}


