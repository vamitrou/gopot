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
	"strings"
	"github.com/BurntSushi/toml"
	"time"
)

type Config struct {
	Ports       []int
	SMTP_Server string
	Mail_From   string
	Mail_To     string
	Timeout     int64
}

var conf Config
var conn_cache = make(map[string]int64)

func send_mail(ip string, port string) bool {
        if conf.SMTP_Server == "" || conf.Mail_From == "" || conf.Mail_To == "" {
            fmt.Printf("SMTP settings not configured.. \n")
            return false
        }
	c, err := smtp.Dial(conf.SMTP_Server)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	// Set the sender and recipient.
	c.Mail(conf.Mail_From)
	c.Rcpt(conf.Mail_To)
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
        return true
}

func report(ip string, port string) {
	fmt.Printf("Connection attempt from %s at port %s\n", ip, port)
	now := time.Now().Unix()

	key := fmt.Sprintf("%s%s", ip, port)
	v, found := conn_cache[key]
	if !found || v+conf.Timeout < now {
            sent := send_mail(ip, port)
		if sent {
                    fmt.Printf("Mail sent. \n")
                }
		conn_cache[key] = now
	}
}

func handleConnection(c net.Conn, port string) {
	ip := strings.Split(c.RemoteAddr().String(), ":")[0]
	report(ip, port)
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

func load_config(path string) {
	b, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Fatal(err)
	}
	tomlData := string(b)
	if _, err := toml.Decode(tomlData, &conf); err != nil {
		log.Fatal(err)
	}
}

func main() {
        load_config("config.toml")
	for i := 0; i < len(conf.Ports); i++ {
		PORT := fmt.Sprintf(":%d", conf.Ports[i])
		go serve(PORT)
	}

	for {
		time.Sleep(10)
	}
}
