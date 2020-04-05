package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gobot.io/x/gobot/drivers/gpio"
)

// Server serves the default routes via HTTP.
type Server struct {
	servo           *gpio.ServoDriver
	currentPosition uint8
	router          chi.Router
}

// New returns a new Server instance with default routes.
func New(servo *gpio.ServoDriver) *Server {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	s := Server{
		servo:           servo,
		router:          r,
		currentPosition: 0,
	}

	r.Get("/{angle}", s.HelloServer)

	return &s
}

// Run ...
func (s *Server) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%v", "9999"), s.router)
}

// HelloServer ...
func (s *Server) HelloServer(w http.ResponseWriter, r *http.Request) {

	sangle := chi.URLParam(r, "angle")
	s.r.
	angle, err := strconv.Atoi(sangle)
	if err != nil {
		fmt.Fprintf(w, "not a valid angle %s", sangle)
	}
	s.servo.Move(uint8(angle))

	fmt.Fprintf(w, "Moving servo to %d", angle)
}