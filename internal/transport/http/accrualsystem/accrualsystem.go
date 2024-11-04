package accrualsystem

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/golovanevvs/gophermart/internal/customerrors"
	"github.com/golovanevvs/gophermart/internal/model"
)

type AccrualSystemInt interface {
	GetAPIOrders(ctx context.Context, orderNumber int) (model.AccrualSystem, error)
}

type accrualSystemStr struct {
	Address string
}

func NewAccrualSystem(address string) *accrualSystemStr {
	return &accrualSystemStr{
		Address: address,
	}
}

func (as *accrualSystemStr) GetAPIOrders(ctx context.Context, orderNumber int) (model.AccrualSystem, error) {
	fmt.Printf("Запуск GetAPIOrders\n")
	// создание экземпляра клиента
	client := &http.Client{}

	// создание запроса, получение ответа
	resp, err := client.Get(as.Address)
	if err != nil {
		return model.AccrualSystem{}, fmt.Errorf("%v: %v", customerrors.ClientError500, err)
	}

	// отложенное закрытие тела ответа
	defer resp.Body.Close()

	// проверка статуса ответа
	switch resp.StatusCode {

	// успешная обработка запроса
	case 200:
		// проверка Content-Type
		if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
			return model.AccrualSystem{}, fmt.Errorf("%v: %v; требуется application/json", customerrors.InvalidContentType400, contentType)
		}

		// десериализация тела ответа в структуру
		var accruelSystem model.AccrualSystem
		if err := json.NewDecoder(resp.Body).Decode(&accruelSystem); err != nil {
			return model.AccrualSystem{}, fmt.Errorf("%v: %v", customerrors.DecodeJSONError500, err)
		}

		return accruelSystem, nil

	// если заказ не зарегистрирован в системе
	case 204:
		return model.AccrualSystem{}, fmt.Errorf("%v", customerrors.ASOrderNotRegistered204)

	// если превышено количество запросов к сервису
	case 429:
		// проверка Content-Type
		if contentType := resp.Header.Get("Content-Type"); contentType != "text/plain" {
			return model.AccrualSystem{}, fmt.Errorf("%v: %v; требуется text/plain", customerrors.InvalidContentType400, contentType)
		}

		// получение Retry-After
		contentType := resp.Header.Get("Retry-After")
		retryAfter, err := strconv.Atoi(string(contentType))
		if err != nil {
			return model.AccrualSystem{}, fmt.Errorf("%v", customerrors.AtoiError500)
		}

		// чтение тела ответа
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return model.AccrualSystem{}, fmt.Errorf("%v", customerrors.ResponseBodyError500)
		}
		message := string(body)

		return model.AccrualSystem{
			RetryAfter: retryAfter,
			Message:    message,
		}, fmt.Errorf("%v", customerrors.ASTooManyRequests429)

	// если возникла внутренняя ошибка сервера
	default:
		return model.AccrualSystem{}, fmt.Errorf("%v", customerrors.InternalServerError500)
	}
}
