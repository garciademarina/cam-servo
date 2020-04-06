package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Server serves the default routes via HTTP.
type Server struct {
	horizonalServoChan     chan uint8
	vercaticalServoChan    chan uint8
	horizontalCurrentAngle int
	verticalCurrentAngle   int
	router                 chi.Router
}

// New returns a new Server instance with default routes.
func New(horizonalServoChan, vercaticalServoChan chan uint8) *Server {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	s := Server{
		horizonalServoChan:     horizonalServoChan,
		vercaticalServoChan:    vercaticalServoChan,
		horizontalCurrentAngle: 60,
		verticalCurrentAngle:   30,
		router:                 r,
	}

	r.Get("/move/horizontal/{angle}", s.MoveHorizontally)
	r.Get("/move/vertical/{angle}", s.MoveVertically)
	r.Get("/up", s.MoveUp)
	r.Get("/down", s.MoveDown)
	r.Get("/left", s.MoveLeft)
	r.Get("/right", s.MoveRight)

	go func() {
		time.Sleep(2 * time.Second)
		// init position horizontal server
		s.horizonalServoChan <- uint8(scaleBetween(s.horizontalCurrentAngle, 11, 52, 0, 180))
		// init position vertical server
		s.vercaticalServoChan <- uint8(scaleBetween(s.verticalCurrentAngle, 10, 40, 0, 180))
	}()

	return &s
}

// Run ...
func (s *Server) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%v", "9999"), s.router)
}

// MoveUp ...
func (s *Server) MoveUp(w http.ResponseWriter, r *http.Request) {
	if s.verticalCurrentAngle != 0 {
		s.verticalCurrentAngle -= 10
		s.vercaticalServoChan <- uint8(scaleBetween(s.verticalCurrentAngle, 10, 40, 0, 180))
	}
}

// MoveDown ...
func (s *Server) MoveDown(w http.ResponseWriter, r *http.Request) {
	if s.verticalCurrentAngle != 180 {
		s.verticalCurrentAngle += 10
		s.vercaticalServoChan <- uint8(scaleBetween(s.verticalCurrentAngle, 10, 40, 0, 180))
	}
}

// MoveLeft ...
func (s *Server) MoveLeft(w http.ResponseWriter, r *http.Request) {
	if s.verticalCurrentAngle != 0 {
		s.horizontalCurrentAngle -= 10
		s.horizonalServoChan <- uint8(scaleBetween(s.horizontalCurrentAngle, 10, 40, 0, 180))
	}
}

// MoveRight ...
func (s *Server) MoveRight(w http.ResponseWriter, r *http.Request) {
	if s.verticalCurrentAngle != 180 {
		s.horizontalCurrentAngle += 10
		s.horizonalServoChan <- uint8(scaleBetween(s.horizontalCurrentAngle, 10, 40, 0, 180))
	}
}

// MoveHorizontally ...
func (s *Server) MoveHorizontally(w http.ResponseWriter, r *http.Request) {
	// movement 11-52
	sangle := chi.URLParam(r, "angle")
	angle, err := strconv.Atoi(sangle)
	if err != nil {
		fmt.Fprintf(w, "not a valid angle %s", sangle)
	}

	s.horizonalServoChan <- uint8(scaleBetween(angle, 11, 52, 0, 180))
	fmt.Fprintf(w, "Moving horizontal servo to %d", angle)
}

// MoveVertically ...
func (s *Server) MoveVertically(w http.ResponseWriter, r *http.Request) {
	// movement 0-45
	sangle := chi.URLParam(r, "angle")
	angle, err := strconv.Atoi(sangle)
	if err != nil {
		fmt.Fprintf(w, "not a valid angle %s", sangle)
	}

	s.vercaticalServoChan <- uint8(scaleBetween(angle, 10, 40, 0, 180))
	fmt.Fprintf(w, "Moving vertical servo to %d", angle)
}

func scaleBetween(unscaledNum, minAllowed, maxAllowed, min, max int) uint8 {
	return uint8((maxAllowed-minAllowed)*(unscaledNum-min)/(max-min) + minAllowed)
}
