package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/achun/db/mysql"
	"github.com/braintree/manners"
	"github.com/gorilla/sessions"

	_ "blog"
	"controllers"
	. "global"
	_ "i18n"
	_ "i18n/mysql"
	_ "i18n/mysql/zh-cn"
	_ "i18n/zh-cn"
)

func main() {
	if Conf == nil {
		os.Exit(2)
		return
	}

	if BlogId != 0 {
		err := OpenDb()
		if err == nil {
			_, err = Db.Collection("users")
		}
		if err != nil {
			panic(err)
		}
		SessionStore = sessions.NewFilesystemStore(Conf.GetDefault("session.path", "").(string), []byte(Conf.GetDefault("session.secret", "").(string)))
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
