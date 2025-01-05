package http

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

type Server struct {
	Port       string
	server     *http.Server
	router     *mux.Router
	queue      *amqp.Queue
	channel    *amqp.Channel
	connection *amqp.Connection
}

// handleIndex handles requests to the root ("/") and prints a personalized message.
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {

	messageBody := "Hello, RabbitMQ!"
	err := s.channel.Publish(
		"",           // Exchange (empty string for default)
		s.queue.Name, // Routing key (queue name)
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(messageBody),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}
	defer s.connection.Close()
	defer s.channel.Close()
}

// NewServer creates and returns a new Server instance.
func NewServer(port string) *Server {
	// Read variables from the environment
	rabbitmqHost := os.Getenv("RABBITMQ_HOST")
	rabbitmqPort := os.Getenv("RABBITMQ_PORT")
	rabbitmqLogin := os.Getenv("RABBITMQ_LOGIN")
	rabbitmqPassword := os.Getenv("RABBITMQ_PASSWORD")
	rabbitmqUserQueue := os.Getenv("RABBITMQ_USER_QUEUE")

	rabbitmqURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitmqLogin, rabbitmqPassword, rabbitmqHost, rabbitmqPort)
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	queue, err := ch.QueueDeclare(
		rabbitmqUserQueue, // Queue name
		true,              // Durable
		false,             // Auto-delete
		false,             // Exclusive
		false,             // No-wait
		nil,               // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	if port == "" {
		port = "8080"
	}

	s := &Server{
		Port:       port,
		server:     &http.Server{},
		router:     mux.NewRouter(),
		queue:      &queue,
		channel:    ch,
		connection: conn,
	}

	// Define routes
	router := s.router.PathPrefix("/user").Subrouter()
	router.HandleFunc("/{name}", s.handleIndex).Methods("GET")

	return s
}

// Open starts the server on the specified port.
func (s *Server) Open() error {
	s.server.Addr = ":" + s.Port
	s.server.Handler = s.router
	return s.server.ListenAndServe()
}
