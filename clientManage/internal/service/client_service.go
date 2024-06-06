package service

import (
	"clientManage/internal/domain/model"
	"clientManage/internal/repository"
	"clientManage/internal/transport/messaging"
	"errors"
	"log"
)

type ClientService struct {
	repo      repository.ClientRepository
	messaging messaging.ClientMessaging
}

func NewClientService(repo repository.ClientRepository, messaging messaging.ClientMessaging) *ClientService {
	return &ClientService{repo: repo, messaging: messaging}
}

func (s *ClientService) CreateClient(client *model.Client) error {
	existingUser, err := s.repo.GetClientByEmail(client.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("user with this email already exists")
	}

	err = s.repo.CreateClient(client)
	if err != nil {
		log.Printf("Error creating client: %v", err)
		return err
	}

	err = s.messaging.PublishClientCreated(client)
	if err != nil {
		log.Printf("Error publishing client created message: %v", err)
		return err
	}

	return nil
}

func (s *ClientService) GetClientByID(id int64) (*model.Client, error) {
	client, err := s.repo.GetClientByID(id)
	if err != nil {
		log.Printf("Error getting client by ID: %v", err)
		return nil, err
	}
	return client, nil
}

func (s *ClientService) GetClientByEmail(email string) (*model.Client, error) {
	client, err := s.repo.GetClientByEmail(email)
	if err != nil {
		log.Printf("Error getting client by Email: %v", err)
		return nil, err
	}
	return client, nil
}

func (s *ClientService) UpdateClient(client *model.Client) error {
	err := s.repo.UpdateClient(client)
	if err != nil {
		log.Printf("Error updating client: %v", err)
		return err
	}
	return nil
}

func (s *ClientService) DeleteClient(id int64) error {
	err := s.repo.DeleteClient(id)
	if err != nil {
		log.Printf("Error deleting client: %v", err)
		return err
	}
	return nil
}

func (s *ClientService) ListClients(offset, limit int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.Client, error) {
	clients, err := s.repo.ListClients(offset, limit, filters, sortBy, sortOrder)
	if err != nil {
		log.Printf("Error listing clients: %v", err)
		return nil, err
	}
	return clients, nil
}
