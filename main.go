package main

func main() {
	done := make(chan struct{})
	defer close(done)
	NewIDGenerator(done)
}
