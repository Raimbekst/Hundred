package domain

import "errors"

var (
	ErrUserAlreadyExist           = errors.New("пользователь с номером или почтой уже существует")
	ErrPasswordNotMatch           = errors.New("пароли не совпадают")
	ErrNotFound                   = errors.New("не найдено")
	ErrCityNotFound               = errors.New("город не найден")
	ErrUserNotExist               = errors.New("пользователя с такими данными не существует")
	ErrUserBlocked                = errors.New("пользователь заблокирован")
	ErrWinnerAlreadyExistInRaffle = errors.New("Розыгрыш имеет победителя ")
	ErrUpdateNotification         = errors.New("изменить отправленный уведомлений не получается")
	ErrCheckBlocked               = errors.New("чек заблокирован")
	ErrTokenAlreadyExist          = errors.New("токен уже существует")
)
