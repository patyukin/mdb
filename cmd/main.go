package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/patyukin/mdb/internal/config"
	"github.com/patyukin/mdb/internal/database"
	"github.com/patyukin/mdb/internal/database/compute"
	"github.com/patyukin/mdb/internal/database/compute/parser"
	"github.com/patyukin/mdb/internal/database/storage"
	"github.com/patyukin/mdb/internal/database/storage/engine"
	"github.com/patyukin/mdb/pkg/logger"
	"go.uber.org/zap"
	"log"
	"os"
	"strings"
)

func main() {
	configPath := flag.String("config_path", "", "Config path")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	l, err := logger.InitLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	engn := engine.New()
	strg := storage.New(engn, l)
	prsr := parser.New()
	cmpt := compute.New(prsr, l)

	dbase := database.New(cmpt, strg, l)

	scanner := bufio.NewScanner(os.Stdin)
	l.Info("Database started. Waiting for commands...")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			l.Error("Error reading from input")
			break
		}

		input := scanner.Text()
		input = strings.TrimSpace(input)
		if input == "" {
			l.Info("Empty command received")
			continue
		}

		var result string
		result, err = dbase.HandleQuery(input)
		if err != nil {
			l.Error("failed c.ProcessRequest", zap.Error(err))
		} else if result != "" {
			l.Info("Request processed successfully", zap.String("result", result))
		}
	}

	if err = scanner.Err(); err != nil {
		l.Error("Error reading from input", zap.Error(err))
	}
}
