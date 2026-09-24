package access

import "github.com/gofiber/fiber/v2"

type Server struct {
	service *Service
	app     *fiber.App
}

func NewServer(service *Service) *Server {
	server := &Server{
		service: service,
		app:     fiber.New(),
	}
	server.routes()
	return server
}

func (s *Server) App() *fiber.App {
	return s.app
}

func (s *Server) routes() {
	s.app.Post("/login", s.handleLogin)
	s.app.Post("/signup", s.handleSignup)

	protected := s.app.Group("", s.authenticate)
	protected.Get("/me", s.handleMe)
	protected.Post("/users/registration-link", s.handleRegistrationLink)
}