package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/achun/db/mysql"
	"github.com/braintree/manners"

	"controllers"
	. "global"
)

func main() {
	if Conf == nil {
		os.Exit(2)
		return
	}

	signal.Notify(manners.ShutdownChannel, syscall.SIGINT)
	for {
		oldListener, err := net.Listen("tcp", ":"+Port)
		if err != nil {
			os.Exit(2)
			return
		}
		GracefulListen = manners.NewListener(oldListener)
		err = manners.Serve(GracefulListen, controllers.HandlerMux(Mux.ServeHTTP))
		if err == nil {
			break
		} else {
			fmt.Println(err)
			break
		}
	}
	FireShutDown()
}
