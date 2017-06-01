package main

var a int

func display(msg int, c chan bool) {
	c <- true
	a = msg
	println("msg: ", msg)
}

func main() {
	c := make(chan bool)

	for i := 0; i < 1000000; i++ {
		go display(i, c)
	}
	<-c

	println("display first message: ", a)

}
