package main

import "math/rand"
import "time"

const serviceName = "Authentication"

func main() {
	rand.Seed(time.Now().UnixNano())
	handler := CreateRestHandler()
	defer handler.service.Close()
	handler.StartListening(8080)
}

