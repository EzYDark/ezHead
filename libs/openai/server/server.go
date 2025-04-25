package openai

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ezydark/ezHead/libs/openai/server/handles"
	"github.com/rs/zerolog/log"
)

type OpenAIServer struct {
	address string
	port    int
}

// Default OpenAI-compatible server Configuration
var (
	oai_server  *OpenAIServer
	def_address = "localhost"
	def_port    = 8080
)

// Initialize OpenAI Server with default configuration
func InitServer() (*OpenAIServer, error) {
	if oai_server != nil {
		return nil, errors.New("OpenAIServer already initialized")
	}

	oai_server = &OpenAIServer{
		address: def_address,
		port:    def_port,
	}

	return oai_server, nil
}

func GetServer() (*OpenAIServer, error) {
	if oai_server == nil {
		return nil, errors.New("OpenAIServer not initialized")
	}
	return oai_server, nil
}

func (oai_server *OpenAIServer) Start() {
	log.Debug().Msgf("Starting OpenAI-compatible server at %s:%d", oai_server.address, oai_server.port)

	http.HandleFunc("/v1/chat/completions", handles.AuthMiddleware(handles.HandleChatCompletions))
	http.HandleFunc("/health", handles.HandleHealthCheck)

	// Channel for server errors
	serverErrors := make(chan error, 1)
	// Start the server in a separate goroutine
	go func() {
		serverErrors <- http.ListenAndServe(fmt.Sprintf("%s:%d", oai_server.address, oai_server.port), nil)
	}()
	// Wait for server errors and log them
	go func() {
		err := <-serverErrors
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("OpenAI-compatible server error while starting:\n%v", err)
		} else if err == http.ErrServerClosed {
			log.Info().Msg("OpenAI-compatible server closed!")
		}
	}()
}

// SetAddress sets the address of the OpenAI Server
func (oai_server *OpenAIServer) SetAddress(address string) (*OpenAIServer, error) {
	if address == "" {
		return nil, errors.New("invalid address")
	}
	oai_server.address = address

	return oai_server, nil
}

// SetPort sets the port of the OpenAI Server
func (oai_server *OpenAIServer) SetPort(port int) (*OpenAIServer, error) {
	if port < 0 || port > 65535 {
		return nil, errors.New("invalid port number (only 0-65535)")
	}
	oai_server.port = port

	return oai_server, nil
}
