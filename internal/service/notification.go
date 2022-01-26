package service

import (
	"HundredToFive/internal/domain"
	"HundredToFive/internal/repository"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
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

func (n *NotificationService) DownloadNotification(file *excelize.File) (*excelize.File, error) {
	page := domain.Pagination{}
	list, err := n.repos.GetAll(page)

	if err != nil {
		return nil, fmt.Errorf("service.DownloadNotification: %w", err)
	}
	id := 2

	getters := map[int]string{1: "все", 2: "выброчно"}
	for _, value := range list.Data {

		file.SetCellValue("Sheet1", "A"+strconv.Itoa(id), value.Id)
		file.SetCellValue("Sheet1", "B"+strconv.Itoa(id), value.Title)
		file.SetCellValue("Sheet1", "C"+strconv.Itoa(id), value.Text)
		file.SetCellValue("Sheet1", "D"+strconv.Itoa(id), value.Status)
		file.SetCellValue("Sheet1", "E"+strconv.Itoa(id), value.Link)
		file.SetCellValue("Sheet1", "F"+strconv.Itoa(id), value.Reference)
		file.SetCellValue("Sheet1", "G"+strconv.Itoa(id), getters[value.Getters])
		id = id + 1
	}
	return file, nil
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

func (n *NotificationService) StoreUsersToken(userId *int, token string) (int, error) {
	return n.repos.StoreUsersToken(userId, token)
}

func (n *NotificationService) GetAllRegistrationTokens() ([]string, error) {
	return n.repos.GetAllRegistrationTokens()
}

func (n *NotificationService) GetNotificationByDate(time int64) ([]domain.Notification, error) {
	return n.repos.GetNotificationByDate(time)
}
