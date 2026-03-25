package sse

import (
	"ai-tutor-backend/internal/dto"
	"ai-tutor-backend/internal/log"
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Client struct {
	// ID is session id
	ID      string
	Channel chan dto.SseChatMessageRequest
	LastSeen time.Time
}

type SseBroker struct {
	Clients       map[string]*Client
	NewClients    chan *Client
	ClosedClients chan string
	mu            sync.RWMutex
	ctx           context.Context
	logger        log.Logger
	cancel        context.CancelFunc
	wg 			  sync.WaitGroup
}

func NewSseBroker(logger log.Logger) *SseBroker {
	ctx, cancel := context.WithCancel(context.Background())
	return &SseBroker{
		Clients:       make(map[string]*Client),
		NewClients:    make(chan *Client),
		ClosedClients: make(chan string),
		logger:        logger,
		ctx:           ctx,
		cancel:        cancel,
	}

}

// Notes:
// broker.ctx.Done() will only be called when someone calls ctx.cancel() or context.CancelFunc

func (broker *SseBroker) Listen() {
	broker.wg.Add(1)
	defer broker.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case client := <-broker.NewClients:
			broker.mu.Lock()

			if existing, ok := broker.Clients[client.ID]; ok {
				close(existing.Channel)
			}

			broker.Clients[client.ID] = client
			
			broker.mu.Unlock()
		case sessionId := <-broker.ClosedClients:
			broker.mu.Lock()

			if client, ok := broker.Clients[sessionId]; ok {
				// close channel
				close(client.Channel)

				// remove a value with key sessionId from hashmap
				delete(broker.Clients, sessionId)
			}

			broker.mu.Unlock()
		case <- ticker.C:
			broker.cleanUpStaleConnection()

		case <-broker.ctx.Done():
			broker.mu.Lock()
			for _, client := range broker.Clients {
				close(client.Channel)
			}
			broker.Clients = make(map[string]*Client)
			broker.mu.Unlock()
			return 

		}

		
	}
}

// remove any client that is inactive for more than 5 minutes
func (broker *SseBroker) cleanUpStaleConnection(){
	broker.mu.Lock()
	defer broker.mu.Unlock()

	now := time.Now()

	for sessionId, client := range broker.Clients {
		if now.Sub(client.LastSeen) > 5*time.Minute {
			// close channel
			close(client.Channel)

			// remove a value with key sessionId from hashmap
			delete(broker.Clients, sessionId)
		}
	}  



}

func (broker *SseBroker) SendEvent(sessionId string, message dto.SseChatMessageRequest) error {
	broker.mu.RLock()
    client, ok := broker.Clients[sessionId]
    broker.mu.RUnlock()

    if !ok {
        broker.logger.Info("client not found", zap.String("sessionId", sessionId))
        return fmt.Errorf("sse broker send event error: session id is not found")
    }

	// how select works ?
	// basically it executes all of the cases
	// in this case, if the first one is waiting too long
	// while the second case is already running
	// when the second case is done before the first one is finished
	// it immediately return fmt.Error, telling that
	// message can't be sent, because the channel with session id
	// is stuck or waiting too long
	select {
		case client.Channel <- message:
			broker.mu.Lock()
			client.LastSeen = time.Now()
			broker.mu.Unlock()
			return nil
		case <-time.After(5 * time.Second):
			broker.logger.Warn("sse broker send event warning: timeout when sending message to channel", zap.String("session_id", sessionId))
			return fmt.Errorf("sse broker send event error: timeout when sending message to channel")
		case <-broker.ctx.Done():
			return fmt.Errorf("sse broker send event error: broker is shutting down")
		
	}

}

func (broker *SseBroker) RegisterClient(sessionId string, clientChan chan dto.SseChatMessageRequest) error {
	client := &Client{
		ID:      sessionId,
		Channel: clientChan,
		LastSeen: time.Now(),
	}

	select {
		case broker.NewClients <- client:
		case <-time.After(5 * time.Second):
			broker.logger.Warn("sse broker register client warning: timeout when sending a new client to register channel", zap.String("session_id", sessionId))
			return fmt.Errorf("sse broker register client error: timeout when sending a new client to register channel")

		case <-broker.ctx.Done():
			return fmt.Errorf("sse broker register client error: broker is shutting down")
		
	}

	return nil

}

func (broker *SseBroker) UnregisterClient(sessionId string) error {
	if _, ok := broker.Clients[sessionId]; !ok {
		return nil 
	} 
	
	broker.ClosedClients <- sessionId

	return nil
}

func (broker *SseBroker) Shutdown(){
	broker.logger.Info("shutting down SSE Broker")
	broker.cancel()
	// block here until all waitgroup are finished
	broker.wg.Wait()

	broker.logger.Info("SSE broker shutdown is completed")
}


