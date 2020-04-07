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
	"os"
	//"strconv"
	"github.com/BurntSushi/toml"
	"gopot/globals"
	"strings"
	"time"
)

func potlog(msg string) {
        fmt.Printf("[+] %s\n", msg)
}

func send_mail(ip string, port string) bool {
	if globals.Conf.SMTP_Server == "" || globals.Conf.Mail_From == "" || len(globals.Conf.Mail_To) == 0 {
		fmt.Printf("SMTP settings not configured.. \n")
		return false
	}
	// Connect to the remote SMTP server.
	c, err := smtp.Dial(globals.Conf.SMTP_Server)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer c.Close()

	// Authentication
	if globals.Conf.SMTP_User != "" && globals.Conf.SMTP_Passwd != "" {
		smtp_hostname := strings.Split(globals.Conf.SMTP_Server, ":")[0]
		auth := smtp.PlainAuth("", globals.Conf.SMTP_User, globals.Conf.SMTP_Passwd, smtp_hostname)
		err := smtp.Auth(auth)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}

	// Set the sender and recipient.
	c.Mail(globals.Conf.Mail_From)
	for _, addr := range globals.Conf.Mail_To {
		c.Rcpt(addr)
	}
	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer wc.Close()
	buf := bytes.NewBufferString(fmt.Sprintf("To: %s\r\n"+
		"Subject: Network intrusion alert in %s\r\n\r\n"+
		"Connection attempt from %s at port %s\n", strings.Join(globals.Conf.Mail_To, ";"), globals.Hostname, ip, port))

	if _, err = buf.WriteTo(wc); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func report(ip string, port string) {
	fmt.Printf("Connection attempt from %s at port %s\n", ip, port)

	if len(globals.Conf.SMTP_Server) > 0 {
		now := time.Now().Unix()
		key := fmt.Sprintf("%s%s", ip, port)
		v, found := globals.Conn_cache[key]
		if !found || v+globals.Conf.Timeout < now {
			fmt.Println("Sending email.")
			sent := send_mail(ip, port)
			if sent {
				fmt.Printf("Mail sent. \n")
				globals.Conn_cache[key] = now
			}
		}
	}
}

func handleConnection(c net.Conn, port string) {
	fmt.Println("DEBUG: handleConnection")
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
	if _, err := toml.Decode(tomlData, &globals.Conf); err != nil {
		log.Fatal(err)
	}
}

func main() {
	potlog("GoPot started")
	load_config("config.toml")
	potlog("config.toml was loaded")

	globals.Hostname, _ = os.Hostname()

	for i := 0; i < len(globals.Conf.Ports); i++ {
		PORT := fmt.Sprintf(":%d", globals.Conf.Ports[i])
		go serve(PORT)
                potlog(fmt.Sprintf("serving on port: %d", globals.Conf.Ports[i]))
	}

	for {
		time.Sleep(10)
	}
}
