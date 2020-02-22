package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/ethanent/protocore-go"
)

var sch protocore.Schema = protocore.Schema{}

func main() {
	sch.Components = append(sch.Components, protocore.Component{
		Name: "tvarint",
		Kind: protocore.Varint,
	})

	sch.Components = append(sch.Components, protocore.Component{
		Name: "tbuf",
		Kind: protocore.Buffer,
	})

	sch.Components = append(sch.Components, protocore.Component{
		Name: "tstr",
		Kind: protocore.String,
	})

	sch.Components = append(sch.Components, protocore.Component{
		Name: "tuint",
		Kind: protocore.UInt,
		Size: 16,
	})

	hostTestServer()
}

func hostTestServer() {
	srv, err := net.Listen("tcp", ":8080")

	if err != nil {
		panic(err)
	}

	for {
		conn, err := srv.Accept()

		fmt.Println("Got conn")

		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	abs := protocore.NewStreamingAbstractor()

	abs.Register("tester", sch)

	go func() {
		for {
			ch := make(chan map[string]interface{}, 1)

			abs.Handle("tester", ch)

			fmt.Println("got msg!", <-ch)
		}
	}()

	go io.Copy(abs, c)
	go io.Copy(c, abs)

	go func() {
		for {
			time.Sleep(500 * time.Millisecond)

			err := abs.Send("tester", map[string]interface{}{
				"tvarint": -58,
				"tbuf":    []byte{56, 64, 69, 62, 42, 255},
				"tstr":    "HeY ThERe! 5236",
				"tuint":   uint(220),
			})

			if err != nil {
				panic(err)
			}
		}
	}()
}

func tests() {
	d, _, err := sch.Parse([]byte{128, 56, 2, 44, 88, 7, 0, 6, 56,
		69, 69, 69, 42, 0, 0, 15, 72, 101,
		89, 32, 84, 72, 69, 114, 101, 33, 32,
		51, 53, 52, 54, 21, 2})

	fmt.Println(d)
	fmt.Println(err)

	b, err := sch.Build(map[string]interface{}{
		"tvarint": -56,
		"tbuf":    []byte{56, 69, 69, 69, 42, 0},
		"tstr":    "HeY THEre! 3546",
		"tuint":   uint(533),
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(b)
}
