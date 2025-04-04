package main

import (
	"dantaautotool/config"
	"dantaautotool/internal/listener"
	"dantaautotool/internal/service"
	"dantaautotool/pkg/utils/http"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// The ASCII art of "HELLO" is generated by https://patorjk.com/software/taag/
const HELLO = "\n" +
	"ooooo ooooo ooooooooooo ooooo       ooooo         ooooooo   \n" +
	" 888   888   888    88   888         888        o888   888o \n" +
	" 888ooo888   888ooo8     888         888        888     888 \n" +
	" 888   888   888    oo   888      o  888      o 888o   o888 \n" +
	"o888o o888o o888ooo8888 o888ooooo88 o888ooooo88   88ooo88   \n" +
	"                                                            \n"

func main() {

	// Record the start time
	t := time.Now()
	log.Info().Msgf("[main] Start time: %s", t.Format(time.RFC3339))

	// Initialize logger with default log level (Info)
	// If user specifies a log level in the command line arguments, we will override it later
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	fmt.Print(HELLO)

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Err(err).Msgf("[main] Failed to load environment variables: %s", err)
	}

	// Load configuration
	config.LoadConfig()

	// Initialize Lark client
	err = http.InitLarkClient()
	if err != nil {
		log.Fatal().Err(err).Msg("[main] Failed to initialize Lark client")
		return
	}

	// Initialize services
	larkIMService := service.NewLarkIMService()
	larkEmailService := service.NewLarkEmailService()
	larkDocService := service.NewLarkDocService()
	githubService := service.NewGithubService()
	dantaService := service.NewDantaService(larkDocService, larkEmailService, githubService)

	// Initialize listeners
	larkListener := listener.NewLarkListener(larkDocService, larkIMService, dantaService)
	if larkListener == nil {
		log.Fatal().Msg("[main] Failed to create LarkListener")
		return
	}

	if err := larkListener.Start(); err != nil {
		log.Fatal().Err(err).Msg("[main] Failed to start LarkListener")
	}

	// Record the end time
	t = time.Now()
	log.Info().Msgf("[main] End time: %s", t.Format(time.RFC3339))

	// Wait for signal to exit
	<-make(chan struct{})
	os.Exit(0)
}
