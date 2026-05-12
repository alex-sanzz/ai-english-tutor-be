package handler

import (
	"ai-tutor-backend/internal/apperr"
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/usecase"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ConversationQuestionHandler struct {
	usecase usecase.ConversationQuestionUseCase
	log log.Logger
}

func NewConversationQuestionHandler(usecase usecase.ConversationQuestionUseCase, log log.Logger) *ConversationQuestionHandler{
	return &ConversationQuestionHandler{
		usecase: usecase,
		log: log,
	}
}

func (h *ConversationQuestionHandler) FindById(c *gin.Context){
	id := c.Param("id")

	if id == "" {
		writeError(c, apperr.BadRequest("400", "id is required", fmt.Errorf("id is required")))
	}

	conversation, err := h.usecase.FindById(c.Request.Context(), id)

	if err != nil {
		h.log.Error("conversation question handler error", zap.Error(err))
		writeError(c, err)
	}

	c.JSON(200, conversation)

}

func (h *ConversationQuestionHandler) FindAll(c *gin.Context){
	sessionRoomId := c.Query("sessionRoomId")

	if sessionRoomId == "" {
		writeError(c, apperr.BadRequest("400", "session room id is required", fmt.Errorf("session room id is required")))
	}

	conversations, err := h.usecase.FindAll(c.Request.Context(), sessionRoomId)

	if err != nil {
		h.log.Error("conversation question handler error", zap.Error(err))
		writeError(c, err)
	}

	c.JSON(200, conversations)

}

func (h *ConversationQuestionHandler) FindAllAnsweredQuestion(c *gin.Context){
	sessionRoomId := c.Query("sessionRoomId")

	if sessionRoomId == "" {
		writeError(c, apperr.BadRequest("400", "session room id is required", fmt.Errorf("session room id is required")))
	}

	conversations, err := h.usecase.FindAllAnsweredQuestion(c.Request.Context(), sessionRoomId)

	if err != nil {
		h.log.Error("conversation question handler error", zap.Error(err))
		writeError(c, err)
	}

	c.JSON(200, conversations)

}

func (h *ConversationQuestionHandler) GenerateQuestion(c *gin.Context){
	sessionRoomId := c.Query("sessionRoomId")

	if sessionRoomId == "" {
		writeError(c, apperr.BadRequest("400", "session room id is required", fmt.Errorf("session room id is required")))
	}

	err := h.usecase.GenerateQuestion(c.Request.Context(), sessionRoomId)

	if err != nil {
		writeError(c, err)
	}

	c.JSON(204, nil)

}

func (h *ConversationQuestionHandler) AnswerQuestion(c *gin.Context){
	id := c.Query("id")

	if id == "" {
		writeError(c, apperr.BadRequest("400", "ID is required", fmt.Errorf("id is required")))
	}

	alternateVersion := c.DefaultQuery("alternative-version", "0")
	culturalContext := c.DefaultQuery("cultural-context", "")
	paraphraseVersion := c.DefaultQuery("paraphrase-version", "false")

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

	err = h.usecase.AnswerQuestion(c.Request.Context(), id, alternateVersion, culturalContext, paraphraseVersion, b)

	if err != nil {
		h.log.Error("send sse audio message handler error", zap.Error(err))
		writeError(c, err)
	}



	c.JSON(204, nil)

}





