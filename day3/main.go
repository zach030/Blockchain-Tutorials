package main

func main() {
	blockchain := NewBlockchain()
	cli := CLI{blockchain}
	cli.Run()
}
