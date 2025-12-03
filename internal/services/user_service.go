package services

import (
	"errors"
	"sync"
	"time"

	"go-microservice/internal/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid input")
)

type UserService struct {
	users  map[int]models.User
	mu     sync.RWMutex
	nextID int
}

func NewUserService() *UserService {
	return &UserService{
		users:  make(map[int]models.User),
		nextID: 1,
	}
}

func (s *UserService) GetAll() ([]models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]models.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users, nil
}

func (s *UserService) GetByID(id int) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (s *UserService) Create(req models.CreateUserRequest) (*models.User, error) {
	if req.Name == "" || req.Email == "" {
		return nil, ErrInvalidInput
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверка уникальности email
	for _, user := range s.users {
		if user.Email == req.Email {
			return nil, errors.New("email already exists")
		}
	}

	now := time.Now()
	user := models.User{
		ID:        s.nextID,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.users[user.ID] = user
	s.nextID++

	return &user, nil
}

func (s *UserService) Update(id int, req models.UpdateUserRequest) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Проверка уникальности email
		for _, u := range s.users {
			if u.ID != id && u.Email == req.Email {
				return nil, errors.New("email already exists")
			}
		}
		user.Email = req.Email
	}

	user.UpdatedAt = time.Now()
	s.users[id] = user

	return &user, nil
}

func (s *UserService) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return ErrUserNotFound
	}

	delete(s.users, id)
	return nil
}
