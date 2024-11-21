package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/NorskHelsenett/ror-ms-switchboard/internal/httpserver"
	"github.com/NorskHelsenett/ror-ms-switchboard/internal/rabbitmq/msswitchboardrabbitmqdefinitions"
	"github.com/NorskHelsenett/ror-ms-switchboard/internal/rabbitmq/msswitchboardrabbitmqhandler"
	"github.com/NorskHelsenett/ror-ms-switchboard/internal/ror"
	"github.com/NorskHelsenett/ror-ms-switchboard/internal/switchboardconfig"
	"github.com/NorskHelsenett/ror-ms-switchboard/internal/switchboardconnections"

	"github.com/NorskHelsenett/ror/pkg/config/configconsts"
	"github.com/NorskHelsenett/ror/pkg/rlog"
	"github.com/NorskHelsenett/ror/pkg/telemetry/trace"

	"github.com/spf13/viper"

	// https://blog.devgenius.io/know-gomaxprocs-before-deploying-your-go-app-to-kubernetes-7a458fb63af1
	"go.uber.org/automaxprocs/maxprocs"
)

func init() {
	_, _ = maxprocs.Set(maxprocs.Logger(rlog.Infof))
}

func main() {
	cancelChan := make(chan os.Signal, 1)
	stop := make(chan struct{})
	// catch SIGETRM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	rlog.Info("Switchboard micro service starting")
	switchboardconfig.Load()

	ror.SetupRORClient()
	switchboardconnections.InitConnections()
	msswitchboardrabbitmqdefinitions.InitOrDie()

	if viper.GetBool(configconsts.ENABLE_TRACING) {
		go func() {
			trace.ConnectTracer(stop, "ror-switchboard", viper.GetString(configconsts.OPENTELEMETRY_COLLECTOR_ENDPOINT))
			sig := <-cancelChan
			_, _ = fmt.Println()
			_, _ = fmt.Println(sig)
			stop <- struct{}{}
		}()
	}

	go func() {
		httpserver.InitHttpServer()
		sig := <-cancelChan
		_, _ = fmt.Println()
		_, _ = fmt.Println(sig)
		stop <- struct{}{}
	}()

	msswitchboardrabbitmqhandler.StartListening()

	sig := <-cancelChan
	rlog.Info("caught signal", rlog.Any("signal", sig))
}
