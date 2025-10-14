package stream

import (
	"wanshow-bingo/sse"

	"github.com/gofiber/fiber/v2"
)

func Get(c *fiber.Ctx) error {
	client := sse.NewClient()
	client.Hub = sse.GetChatHub()
	return client.Bind(c)
}

//func BodyWriter(w *bufio.Writer) {
//	//defer HandleError()
//
//	if w == nil {
//		log.Println("[SSE ClientChannel] - No writer available.")
//		return
//	}
//
//	client := make(chan string)
//
//	chatHub := sse.GetChatHub()
//
//	if chatHub == nil {
//		log.Println("[SSE ClientChannel] - No chat hub available.")
//		return
//	}
//
//	// Register (and defer the deregistration of) the client
//	chatHub.RegisterClient(client)
//	defer chatHub.UnregisterClient(client)
//
//	// Keep-alive ticker to prevent client timeouts
//	ticker := time.NewTicker(30 * time.Second)
//	defer ticker.Stop()
//
//	connectEvent := sse.BuildEvent("hub.connected", "Connected to chat hub.")
//
//	Send(client, w, connectEvent.String())
//
//	for {
//		select {
//		case msg, ok := <-client:
//			if !ok {
//				log.Println("[SSE ClientChannel] - Context Errored.")
//				return
//			}
//			Send(client, w, msg)
//
//		case <-ticker.C:
//			msg := chatHub.BuildConnectionCount()
//			Send(client, w, msg)
//		}
//	}
//}
//
//func Send(client chan string, hub, w *bufio.Writer, msg string) {
//	utils.Debugln("[SSE ClientChannel] - Sending message:", msg)
//
//	//defer HandleError()
//	_, err := w.WriteString("data: " + msg + "\n\n")
//	if err != nil {
//		log.Println("[SSE ClientChannel] - Error writing to writer:", err)
//	}
//	err = w.Flush()
//
//	if err != nil {
//
//		log.Println("[SSE ClientChannel] - Error flushing writer:", err)
//	}
//
//	utils.Debugln("[SSE ClientChannel] - Sent message:", msg)
//	return
//}
//
//func HandleError() {
//	log.Println("An error occurred in ChatStream.")
//
//	r := recover()
//
//	if r != nil {
//		log.Printf("Recovered from error in ChatStream - %s\n", r)
//	} else {
//		log.Println("Failed to recover from error in ChatStream.")
//	}
//}
