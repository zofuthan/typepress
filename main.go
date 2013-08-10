package main

import (
	"fmt"
	_ "github.com/achun/db/mysql"
	_ "github.com/achun/typepress/controllers"
	. "github.com/achun/typepress/global"
	"github.com/braintree/manners"
	"net"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	if Conf == nil {
		os.Exit(2)
		return
	}
	port := Conf.GetDefault("blog.port", int64(8080)).(int64)
	signal.Notify(manners.ShutdownChannel)
	for {
		oldListener, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
		if err != nil {
			os.Exit(2)
			return
		}
		GracefulListen = manners.NewListener(oldListener)
		err = manners.Serve(GracefulListen, Mux)
		if err == nil {
			// reload
		} else {
			fmt.Println(err)
			break
		}
	}
}
