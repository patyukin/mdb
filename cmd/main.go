package main

import (
	"github.com/patyukin/mdb/internal/compute"
	"github.com/patyukin/mdb/internal/compute/parser"
	"github.com/patyukin/mdb/internal/config"
	"github.com/patyukin/mdb/internal/storage"
	"github.com/patyukin/mdb/internal/storage/engine"
	"github.com/patyukin/mdb/pkg/logger"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
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
	cmpt := compute.New(prsr, strg, l)

	cmpt.Start()
}
