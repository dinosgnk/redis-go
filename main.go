package main

func main() {
	server := NewServer(Config{
		ListenAddr: ":6379",
	})
	server.Start()
}
