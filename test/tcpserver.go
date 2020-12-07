package test

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:4321")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("started tcp server on localhost:4321")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		bReader := bufio.NewReader(conn)
		bWriter := bufio.NewWriter(conn)

		stdInReader := bufio.NewReader(os.Stdin)

		go func() {
			// also write back from STDIN
			for {
				text, _ := stdInReader.ReadString('\n')
				bWriter.WriteString(text)
				bWriter.Flush()
			}
		}()

		// print out everything on the connection
		go func() {
			for {
				str, err := bReader.ReadString('\n')
				if err == io.EOF {
					log.Println("client disconnected")
					break
				}

				if err != nil {
					log.Println(err)
					break
				}

				log.Printf("%s : %s", conn.RemoteAddr().String(), str)
			}
		}()
	}

}
