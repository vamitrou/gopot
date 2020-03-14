package main

import (
	//"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/smtp"
	//"strconv"
	//"strings"
	"github.com/BurntSushi/toml"
	"time"
)

var smtp_server = "tmu-econ.mail.allianz:25"
var mail_from = "adp-security@allianz.de"
var mail_to = "vasileios.mitrousis@allianz.com"

type Config struct {
	Ports       []int
	SMTP_Server string
	Mail_From   string
	Mail_To     string
}

var conf Config

func send_mail(ip string, port string) {
	c, err := smtp.Dial(smtp_server)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	// Set the sender and recipient.
	c.Mail(mail_from)
	c.Rcpt(mail_to)
	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	defer wc.Close()
	buf := bytes.NewBufferString(fmt.Sprintf("Connection attempt from %s at port %s\n", ip, port))
	if _, err = buf.WriteTo(wc); err != nil {
		log.Fatal(err)
	}
}

func report(ip string, port string) {
	fmt.Printf("Connection attempt from %s at port %s\n", ip, port)
	//send_mail(ip, port)
}

func handleConnection(c net.Conn, port string) {
	report(c.RemoteAddr().String(), port)
	//for {
	//	netData, err := bufio.NewReader(c).ReadString('\n')
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}

	//	temp := strings.TrimSpace(string(netData))
	//	if temp == "STOP" {
	//		break
	//	}

	//	result := strconv.Itoa(rand.Intn(255)) + "\n"
	//	c.Write([]byte(string(result)))
	//}
	c.Close()
}

func serve(port string) {
	l, err := net.Listen("tcp4", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	rand.Seed(time.Now().Unix())

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		handleConnection(c, port)
	}
}

func main() {
	b, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Fatal(err)
	}
	tomlData := string(b)
	if _, err := toml.Decode(tomlData, &conf); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(conf.Ports); i++ {
		PORT := fmt.Sprintf(":%d", conf.Ports[i])
		go serve(PORT)
	}

	for {
		time.Sleep(10)
	}
}
