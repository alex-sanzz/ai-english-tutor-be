package handler

import (
	"ai-tutor-backend/internal/apperr"
	"ai-tutor-backend/internal/dto"
	"ai-tutor-backend/internal/log"
	"ai-tutor-backend/internal/models"
	"ai-tutor-backend/internal/usecase"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SessionRoomHandler struct {
	sessionRoomUseCase *usecase.SessionRoomUseCase
	logger log.Logger
}

func NewSessionRoomHandler(sessionRoomUseCase *usecase.SessionRoomUseCase, log log.Logger) *SessionRoomHandler{
	return &SessionRoomHandler{
		sessionRoomUseCase: sessionRoomUseCase,
		logger: log,
	}
}


func (h *SessionRoomHandler) FindAllTopics(c *gin.Context ){
	roomType := c.Query("roomType")

	userId := c.GetString("userId")

	rooms, err := h.sessionRoomUseCase.FindAllRooms(c.Request.Context(), userId, roomType)

	if err != nil {
		h.logger.Error("session room handler find all topics error", zap.Error(err))
		writeError(c, err)
		return
	}

	if rooms == nil {
		rooms = []*models.SessionRoom{}
	}

	c.JSON(200, rooms)



}

func (h *SessionRoomHandler) CreateSessionRoom(c *gin.Context){
	var room dto.CreateSessionRoomRequest

	if err := c.Bind(&room); err != nil {
		h.logger.Error("session room handler create session room error", zap.Error(err))
		writeError(c, apperr.BadRequest("400", "check the fields", err))
		return 
	}

	userId, exist := c.Get("userId")

	if !exist {
		h.logger.Error("session room handler create session room error", zap.Error(fmt.Errorf("user id is not found, after passing auth middleware")))
		writeError(c, apperr.Unauthorized("401", "unauthorized", fmt.Errorf("user id is not found, after passing auth middleware")))
		return 
	}

	createdRoom, err := h.sessionRoomUseCase.CreateSessionRoom(c.Request.Context(),  userId.(string), room.RoomType, room.Icon, room.Topic)

	if err != nil {
		h.logger.Error("session room handler create session room error", zap.Error(err))
		writeError(c, err)
		return 
	}

	c.JSON(201, createdRoom)

}

func (h *SessionRoomHandler) DeleteSessionRoom(c *gin.Context){
	id := c.Param("id")

	err := h.sessionRoomUseCase.DeleteById(c.Request.Context(), id)

	if err != nil {
		h.logger.Error("session room handler delete session room error", zap.Error(err))
		writeError(c, err)
		return 
	}

	c.JSON(204, nil)


}