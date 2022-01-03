package service

import (
	"HundredToFive/internal/domain"
	"HundredToFive/internal/repository"
)

type NotificationService struct {
	repos repository.Notification
}

func NewNotificationService(repos repository.Notification) *NotificationService {
	return &NotificationService{repos: repos}
}

func (n *NotificationService) Create(noty domain.Notification) (int, error) {
	return n.repos.Create(noty)
}
func (n *NotificationService) CreateForUser(noty domain.Notification) ([]string, int, error) {
	return n.repos.CreateForUser(noty)
}

func (n *NotificationService) GetAll(page domain.Pagination) (*domain.GetAllNotificationsResponse, error) {
	return n.repos.GetAll(page)
}

func (n *NotificationService) GetById(id int) (*domain.Notification, error) {
	return n.repos.GetById(id)
}

func (n *NotificationService) Update(id int, inp domain.Notification) error {
	return n.repos.Update(id, inp)
}

func (n *NotificationService) Delete(id int) error {
	return n.repos.Delete(id)
}

func (n *NotificationService) StoreUsersToken(token string) (int, error) {
	return n.repos.StoreUsersToken(token)
}
