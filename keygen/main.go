package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.TimeOnly})
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	pem, err := x509.MarshalECPrivateKey(pk)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	fileName, ok := os.LookupEnv("FILE_NAME")
	if !ok {
		fileName = "./key.pem"
		logger.Info().Str("file_name", fileName).Msg("No file name specified. Using default.")
	}
	out, err := os.Create(fileName)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	defer out.Close()
	_, err = out.Write(pem)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
}
