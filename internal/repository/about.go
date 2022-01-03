package repository

import (
	"HundredToFive/internal/domain"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
)

type AboutRepos struct {
	db *sqlx.DB
}

func NewAboutRepos(db *sqlx.DB) *AboutRepos {
	return &AboutRepos{db: db}
}

func (a *AboutRepos) Create(about domain.AboutUs) error {
	query := fmt.Sprintf(
		`INSERT INTO 
						%s
					(facebook_link,youtube_link,instagram_link,tiktok_link,whatsapp_link,telegram_link,phone_number,phone_number_2) 
				VALUES 
					($1,$2,$3,$4,$5,$6,$7,$8)`, websiteLinks)

	_, err := a.db.Exec(query, about.FacebookLink, about.YoutubeLink, about.InstagramLink, about.TiktokLink, about.WhatsappLink, about.TelegramLink, about.PhoneNumber, about.PhoneNumber2)

	if err != nil {
		return fmt.Errorf("repository.Create: %w", err)
	}
	return nil
}

func (a *AboutRepos) GetAll() ([]*domain.AboutUs, error) {
	var inp []*domain.AboutUs
	query := fmt.Sprintf("SELECT * FROM %s", websiteLinks)

	err := a.db.Select(&inp, query)

	if err != nil {
		return nil, fmt.Errorf("repsoitory.GetAll: %w", err)
	}
	return inp, nil
}

func (a *AboutRepos) Update(about domain.AboutUs) error {
	setValues := make([]string, 0, reflect.TypeOf(domain.AboutUs{}).NumField())

	if about.FacebookLink != "" {
		setValues = append(setValues, fmt.Sprintf("facebook_link=:facebook_link"))
	}

	if about.YoutubeLink != "" {
		setValues = append(setValues, fmt.Sprintf("youtube_link=:youtube_link"))
	}

	if about.InstagramLink != "" {
		setValues = append(setValues, fmt.Sprintf("instagram_link=:instagram_link"))
	}

	if about.TiktokLink != "" {
		setValues = append(setValues, fmt.Sprintf("tiktok_link=:tiktok_link"))
	}
	if about.WhatsappLink != "" {
		setValues = append(setValues, fmt.Sprintf("whatsapp_link=:whatsapp_link"))
	}
	if about.TelegramLink != "" {
		setValues = append(setValues, fmt.Sprintf("telegram_link=:telegram_link"))
	}
	if about.PhoneNumber != "" {
		setValues = append(setValues, fmt.Sprintf("phone_number=:phone_number"))
	}
	if about.PhoneNumber2 != "" {
		setValues = append(setValues, fmt.Sprintf("phone_number_2=:phone_number_2"))
	}

	setQuery := strings.Join(setValues, ", ")

	if setQuery == "" {
		return errors.New("empty body")
	}

	query := fmt.Sprintf("UPDATE %s SET %s", websiteLinks, setQuery)

	_, err := a.db.NamedExec(query, about)

	if err != nil {
		return fmt.Errorf("repository.Update: %w", err)
	}
	return nil

}

func (a *AboutRepos) Delete() error {
	query := fmt.Sprintf("DELETE FROM %s ", websiteLinks)
	_, err := a.db.Exec(query)
	if err != nil {
		return fmt.Errorf("repository.Delete: %w", domain.ErrNotFound)
	}
	return nil
}
