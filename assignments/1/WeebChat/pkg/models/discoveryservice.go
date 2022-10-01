package models

import (
	"errors"
	"time"
)

type DiscoveryService struct {
	ChatServices []*ChatService
}

func NewDiscoveryService() *DiscoveryService {
	return &DiscoveryService{}
}

func (d *DiscoveryService) AddChatService(chatService *ChatService) error {
	for _, chat := range d.ChatServices {
		if chat.Name == chatService.Name {
			return errors.New("chat service already exists")
		}
	}

	d.ChatServices = append(d.ChatServices, chatService)
	chatService.LastPing = time.Now()
	return nil
}

func (d *DiscoveryService) PingChatService(Name string) {

	for _, chat := range d.ChatServices {
		if chat.Name == Name {
			chat.LastPing = time.Now()
		}
	}
}

func (d *DiscoveryService) GetChatService(Name string) *ChatService {
	for _, chat := range d.ChatServices {
		if chat.Name == Name {
			return chat
		}
	}
	return nil
}

func (d *DiscoveryService) DeleteChatService(Name string) error {
	for i, chat := range d.ChatServices {
		if chat.Name == Name {
			d.ChatServices = append(d.ChatServices[:i], d.ChatServices[i+1:]...)
			return nil
		}
	}
	return errors.New("chat service does not exist")
}

func (d *DiscoveryService) GetChatServices() []*ChatService {
	return d.ChatServices
}
