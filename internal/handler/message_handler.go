package handler

import (
	"ai-tutor-backend/internal/dto"
	"ai-tutor-backend/internal/infrastructure/sse"
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/usecase"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MessageHandler struct {
	sseBroker          *sse.SseBroker
	aiChatUseCase      *usecase.AiChatUseCase
	sessionRoomUseCase *usecase.SessionRoomUseCase
	transcribeUseCase  *usecase.TranscribeUseCase
	log                log.Logger
}

func NewMessageHandler(sseBroker *sse.SseBroker, aiChatUseCase *usecase.AiChatUseCase, sessionRoomUseCase *usecase.SessionRoomUseCase, transcribeUseCase *usecase.TranscribeUseCase, logger log.Logger) *MessageHandler {
	return &MessageHandler{
		sseBroker:          sseBroker,
		aiChatUseCase:      aiChatUseCase,
		sessionRoomUseCase: sessionRoomUseCase,
		transcribeUseCase:  transcribeUseCase,
		log:                logger,
	}
}

func (h *MessageHandler) FindRecentMessages(c *gin.Context) {
	n := c.Query("n")
	sessionId := c.Param("sessionId")

	limit, err := strconv.ParseInt(n, 10, 32)

	if err != nil {
		h.log.Error("message handler find recent messages error :", zap.Error(err))
		writeError(c, err)

		return
	}

	if limit <= 0 {
		h.log.Warn("message handler find recent messages error : limit is less than zero")
		writeError(c, err)
		return
	}

	result, err := h.aiChatUseCase.FindRecentMessages(c.Request.Context(), sessionId, int32(limit))

	if err != nil {
		h.log.Error("message handler find recent messages error :", zap.Error(err))

		writeError(c, err)

		return
	}

	c.JSON(200, result)

}

func (h *MessageHandler) RegisterSse(c *gin.Context) {

	rw := c.Writer
	req := c.Request

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := rw.(http.Flusher)

	if !ok {
	
        h.log.Error("register sse handler error", zap.Error(fmt.Errorf("cannot flush http")))
        return
    
	}

	sessionId, exist := c.GetQuery("sessionId")
	

	if !exist {
		h.sendSseEvent(rw, flusher, "error", dto.SseEvent{
			MessageId: time.Now().String(),
			Message: "failed to initiate session",
			Role: "",
		})
	}


	messageChan := make(chan dto.SseChatMessageRequest)

	// sessionId, err := h.sessionRoomUseCase.CreateSessionRoom(c.Request.Context(), userId.(string), roomType)

	// if err != nil {
	// 	h.log.Error("create session room sse handler error", zap.Error(err))
	// 	h.sendSseEvent(rw, flusher, "error", dto.SseEvent{
	// 		MessageId: time.Now().String(),
	// 		Message: "failed to initiate session",
	// 		Role: "",
	// 	})
	// 	return
	// }

	defer func() {
		h.sseBroker.UnregisterClient(sessionId)
		h.log.Info("sse client unregistered", zap.String("sessionId", sessionId))
	}()

	err := h.sseBroker.RegisterClient(sessionId, messageChan)
	
	if err != nil {
		h.log.Error("register sse error: can't register sse client", zap.Error(err))
		return 
	}

	h.log.Info("client is registered successfully")

	// TODO: send initial event here

	err = h.sendSseEvent(rw, flusher, "session", dto.SseEvent{
		MessageId: "",
		Message: sessionId,
		Role: "",
	})

	if err != nil {
		h.log.Error("register sse error: can't send sse event", zap.Error(err))
		return 
	}

	h.log.Info("session event is sent successfully")

	// Heartbeat ticker to keep connection alive
	heartbeatTicker := time.NewTicker(30 *time.Second)
	defer heartbeatTicker.Stop()

	chats, err := h.aiChatUseCase.FindRecentMessages(c.Request.Context(), sessionId, 1)

	if err != nil {
		h.log.Error("ai chat usecase find recent message error", zap.Error(err))
		h.sendSseEvent(rw, flusher, "error", dto.SseEvent{
			MessageId: time.Now().String(),
			Message: "there is something wrong in the system",
			Role: "",
		})
		return 
	}

	if len(chats) == 0 {
		room, err := h.sessionRoomUseCase.FindById(c.Request.Context(), sessionId)

		if err != nil {
			h.log.Error("session room usecase find by id error", zap.Error(err))
			h.sendSseEvent(rw, flusher, "error", dto.SseEvent{
				MessageId: time.Now().String(),
				Message: "there is something wrong in the system",
				Role: "",
			})
			return 

		}

		err = h.askAi(c.Request.Context(), rw, flusher, sessionId, "the topic is " + room.Topic + ", give me first question, that can help the user practice english", false)

		if err != nil {
			h.log.Error("message handler ask ai method error", zap.Error(err))
			h.sendErrorSse(rw, flusher, "there is something wrong", err)
		}
	}

	
	for {
		select {
		case msg := <-messageChan:
			// result, err := h.geminiUseCase.Chat(c.Request.Context(), string(msg.Message))

			// if err != nil {
			// 	return
			// }

			audioBytes, err := base64.StdEncoding.DecodeString(msg.Message)

			if err != nil {
				h.log.Error("register sse handler error", zap.Error(err))
				h.sendSseEvent(rw, flusher, "error", dto.SseEvent{
					MessageId: time.Now().String(),
					Message: "failed to decode voice",
					Role: "",
				})
				return
			}

			if len(audioBytes) == 0 {
				h.log.Warn("received empty audio data")
				h.sendSseEvent(rw, flusher, "error", dto.SseEvent{
					MessageId: time.Now().String(),
					Message: "audio data is empty",
					Role: "",
				})
				return
			}
			
			// This is used for debugging purposes
			timestamp := time.Now().Format("20060102_150405")
            filename := fmt.Sprintf("audio_%s.webm", timestamp)
            if err := os.WriteFile(filename, audioBytes, 0644); err != nil {
                h.log.Error("failed to save audio file", zap.Error(err), zap.String("filename", filename))
            } else {
                h.log.Info("audio file saved", zap.String("filename", filename), zap.Int("size", len(audioBytes)))
            }

			uploadedFileUrl, err := h.transcribeUseCase.UploadAudioFile(c.Request.Context(), audioBytes)
			
			if err != nil {
				h.log.Error("register sse handler error on uploading audio file", zap.Error(err))
				h.sendSseEvent(rw, flusher, "error", dto.SseEvent{
					MessageId: time.Now().String(),
					Message: "failed to transcribe voice",
					Role: "",
				})
				return
			}

			transcribedText, err := h.transcribeUseCase.TranscribeAudio(c.Request.Context(), uploadedFileUrl)

			if err != nil {
				h.log.Error("register sse handler error on transcribing audio file", zap.Error(err))
				h.sendSseEvent(rw, flusher, "error", dto.SseEvent{
					MessageId: time.Now().String(),
					Message: "failed to transcribe voice",
					Role: "",
				})
				return
			}

			if transcribedText == "" {
				h.log.Error("register sse handler error on transcribing audio file: transcribed text is empty")
				h.sendSseEvent(rw, flusher, "error", dto.SseEvent{
					MessageId: time.Now().String(),
					Message: "transcribed text is empty",
					Role: "",
				})
				return
			}

			h.log.Debug("transcribed text :", zap.String("text", transcribedText))


			if err := h.sendSseEvent(rw, flusher, "data", dto.SseEvent{
				MessageId: time.Now().String(),
				Message:   transcribedText,
				Role:    "user",
			}); err != nil {
				h.log.Error("register sse handler error", zap.Error(err))
				return 
			}

			err = h.askAi(c.Request.Context(), rw, flusher, sessionId, transcribedText, true)

			if err != nil {
				h.log.Error("register sse handler error", zap.Error(err))
				return
			}
		case <-heartbeatTicker.C:

		// if connection is closed
		case <-req.Context().Done():
			h.log.Info("client connection closed")
			return
		}
	}
}

func (h *MessageHandler) sendErrorSse(rw http.ResponseWriter, flusher http.Flusher, errorMessage string, err error){
	h.log.Error("register sse handler error", zap.Error(err))
	
	h.sendSseEvent(rw, flusher, "error", dto.SseEvent{
		MessageId: time.Now().String(),
		Message: errorMessage,
		Role: "",
	})
}

func (h *MessageHandler) askAi(ctx context.Context, rw http.ResponseWriter, flusher http.Flusher, sessionId string, question string, saveRequestMessage bool) error{
	err := h.aiChatUseCase.ChatStream(ctx, sessionId, question, func(id string, chunk string) error {
				

		if err := h.sendSseEvent(rw, flusher, "data", dto.SseEvent{
			MessageId: id,
			Message:   chunk,
			Role:    "assistant",
		}); err != nil {
			
			h.sendErrorSse(rw, flusher, "there is something wrong", err)
			return err
		}

		return nil
	}, func(fullSentences string) error {
		

		if err := h.sendSseEvent(rw, flusher, "finished", dto.SseEvent{
			MessageId: "",
			Message:   "",
			Role:    "",
		}); err != nil {
			h.sendErrorSse(rw, flusher, "there is something wrong", err)
			return err
		}

		return nil
	}, saveRequestMessage)

	if err != nil {
		h.sendErrorSse(rw, flusher, "there is something wrong", err)
		return err
	}

	return nil
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	id := c.Param("id")

	var chat dto.SseChatMessageRequest

	if err := c.ShouldBindBodyWithJSON(&chat); err != nil {
		h.log.Error("send sse message handler error", zap.Error(err))
		writeError(c, err)
		return
	}

	err := h.sseBroker.SendEvent(id, chat)

	if err != nil {
		h.log.Error("send sse message handler error", zap.Error(err))
		writeError(c, err)
		return
	}
}

func (h *MessageHandler) sendSseEvent(w http.ResponseWriter, flusher http.Flusher, eventType string, data dto.SseEvent) error{
	jsonBytes, err := json.Marshal(data)

	if err != nil {
		return fmt.Errorf("send sse event error: %w", err)
	}

	_, err = w.Write([]byte("event: " + eventType + "\n" + "data: " + string(jsonBytes) + "\n\n"))

	if err != nil {
		return fmt.Errorf("send sse event error: %w", err)
	}

	flusher.Flush()

	return nil
	
}

func (h *MessageHandler) SendAudioMessage(c *gin.Context) {
	id := c.Param("id")


	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		h.log.Error("parse multipart form error", zap.Error(err))
	} else if mf := c.Request.MultipartForm; mf != nil {
		// collect keys
		keys := make([]string, 0, len(mf.Value)+len(mf.File))
		for k := range mf.Value {
			keys = append(keys, k)
		}
		for k := range mf.File {
			keys = append(keys, k)
		}

		h.log.Debug("multipart keys", zap.Any("keys", keys))
		h.log.Debug("multipart values", zap.Any("values", mf.Value))
		h.log.Debug("multipart files", zap.Any("files", mf.File))
	} else {
		h.log.Debug("no multipart form present")
	}

	// From form-data, get value from a field named "file"
	fileHeader, err := c.FormFile("file")

	if err != nil {
		h.log.Error("send sse audio message handler error", zap.Error(err))
		writeError(c, err)
		return
	}

	// Opens the uploaded file and returns a multipart.File. But it might be stored into memory, but you still don't read it
	f, err := fileHeader.Open()

	if err != nil {
		h.log.Error("send sse audio message handler error", zap.Error(err))
		writeError(c, err)
		return
	}

	defer f.Close()

	// Reads everything from f and stores the result in a []byte
	b, err := io.ReadAll(f)

	if err != nil {
		h.log.Error("send sse audio message handler error", zap.Error(err))
		writeError(c, err)
		return
	}


	encoded := base64.StdEncoding.EncodeToString(b)

	chat := dto.SseChatMessageRequest{
		Message: encoded,
		UserId:  c.GetString("userId"),
	}

	if err := h.sseBroker.SendEvent(id, chat); err != nil {
		h.log.Error("send sse audio handler error", zap.Error(err))
		writeError(c, err)
		return
	}

	c.Status(http.StatusOK)

}
