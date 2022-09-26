package models

type ChatService struct {
	Name string
	Rooms []*Room
}

func (c *ChatService) CreateRoom (name string) {
	//Create a new room
	//Add room to rooms list
}

func (c *ChatService) DeleteRoom (name string) {
	//Delete room from rooms list
}

func (c *ChatService) GetRoom (name string) {
	//Return room from rooms list
}

func (c *ChatService) GetRooms () {
	//Return rooms list
}