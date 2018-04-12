package main

import "math/rand"
import "time"

const serviceName = "Promise"

func main() {
	rand.Seed(time.Now().UnixNano())
	handler := CreateRestHandler()
	defer handler.service.Close()
	handler.StartListening(4000)
}
