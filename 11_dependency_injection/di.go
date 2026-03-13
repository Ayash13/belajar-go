package dependency_injection

import "fmt"

// Define a contract for logging
type Logger interface {
	Log(message string)
}

// Implementation 1: Console Logger
type ConsoleLogger struct{}

func (c ConsoleLogger) Log(message string) {
	fmt.Println("CONSOLE LOG:", message)
}

// Implementation 2: File Logger (simulated)
type FileLogger struct{}

func (f FileLogger) Log(message string) {
	fmt.Println("FILE LOG (writing to file):", message)
}

// Service that needs a logger
type UserService struct {
	logger Logger // Dependency
}

// Inject dependency via constructor
func NewUserService(l Logger) *UserService {
	return &UserService{
		logger: l,
	}
}

func (s *UserService) CreateUser(name string) {
	// ... database logic
	s.logger.Log(fmt.Sprintf("User %s created successfully", name))
}

func Run() {
	consoleLogger := ConsoleLogger{}
	fileLogger := FileLogger{}

	// Injecting ConsoleLogger
	svc1 := NewUserService(consoleLogger)
	svc1.CreateUser("Ayash")

	// Injecting FileLogger
	svc2 := NewUserService(fileLogger)
	svc2.CreateUser("Budi")
}
