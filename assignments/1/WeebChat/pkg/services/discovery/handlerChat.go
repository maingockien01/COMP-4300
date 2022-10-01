package discovery

import (
	"net/http"
)

func (d *DiscoveryServiceServer) HandlerServicesChat(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		chatServices := d.DiscoveryService.GetChatServices()

		ReturnRestResponse(w, chatServices, http.StatusOK)
	}
}
