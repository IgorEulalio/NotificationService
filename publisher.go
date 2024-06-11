package main

//func main() {
//	// Connect to NATS server
//	nc, err := nats.Connect(nats.DefaultURL)
//	if err != nil {
//		log.Fatalf("Error connecting to NATS: %v", err)
//	}
//	defer nc.Close()
//
//	// Publish a message to the "news" topic
//	err = nc.Publish("news", []byte("Breaking news!"))
//	if err != nil {
//		log.Fatalf("Error publishing message: %v", err)
//	}
//	log.Println("Message published to topic 'news'")
//}
