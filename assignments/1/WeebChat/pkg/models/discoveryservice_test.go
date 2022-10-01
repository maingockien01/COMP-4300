package models

import (
	"testing"
)

func TestAddChatService(t *testing.T) {
	discoveryService := NewDiscoveryService()
	chatService := NewChatService("test")
	err := discoveryService.AddChatService(chatService)

	chatServices := discoveryService.ChatServices

	if err != nil && len(chatServices) != 1 {
		t.Error("Chat service was not added to discovery service")
	}

	if chatServices[0].LastPing.IsZero() {
		t.Error("Chat service LastPing was not updated")
	}
}

func TestDuplicateChatServiceWithTheSamename(t *testing.T) {
	discoveryService := NewDiscoveryService()
	chatService1 := NewChatService("test")
	chatService2 := NewChatService("test")
	err1 := discoveryService.AddChatService(chatService1)
	err2 := discoveryService.AddChatService(chatService2)

	chatServices := discoveryService.ChatServices

	if err1 != nil && err2 == nil && len(chatServices) != 1 {
		t.Error("Duplicated chat service was added to discovery service")
	}
}

func TestGetChatService(t *testing.T) {
	discoveryService := NewDiscoveryService()
	chatService := NewChatService("test")
	err := discoveryService.AddChatService(chatService)

	if err != nil {
		t.Error("Chat service was not added to discovery service")
	}

	retrievedChatService := discoveryService.GetChatService("test")

	if retrievedChatService == nil {
		t.Error("Chat service was not retrieved from discovery service")
	}

}

func TestPingChatServiceUpdateNewLastActive(t *testing.T) {
	discoveryService := NewDiscoveryService()
	chatService := NewChatService("test")
	discoveryService.AddChatService(chatService)

	time1 := discoveryService.GetChatService("test").LastPing

	discoveryService.PingChatService("test")

	time2 := discoveryService.GetChatService("test").LastPing

	if !time2.After(time1) {
		t.Error("Chat service LastPing was not updated")
	}
}

func TestDeleteChatService(t *testing.T) {
	discoveryService := NewDiscoveryService()
	chatService := NewChatService("test")
	discoveryService.AddChatService(chatService)

	err := discoveryService.DeleteChatService("test")

	if err != nil {
		t.Error("Chat service was not deleted from discovery service")
	}
}
