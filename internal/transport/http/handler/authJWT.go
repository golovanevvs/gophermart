package handler

import (
	"context"
	"net/http"
	"strings"
)

// авторизация пользователя по токену заголовка Authorization
func (hd *handlerStr) authByJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// получение заголовка Authorization
		header := r.Header.Get("Authorization")

		// обработка ошибок
		if header == "" {
			http.Error(w, "заголовок авторизации отсутствует", http.StatusUnauthorized)
			return
		}

		headerSplit := strings.Split(header, " ")

		if len(headerSplit) != 2 || headerSplit[0] != "Bearer" {
			http.Error(w, "некорректный заголовок авторизации", http.StatusUnauthorized)
			return
		}

		if len(headerSplit[1]) == 0 {
			http.Error(w, "заголовок авторизации не содержит токен", http.StatusUnauthorized)
			return
		}

		// получение userID из JWT
		userID, err := hd.sv.GetUserIDFromJWT(headerSplit[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// оформление передачи userID с помощью контекста
		ctx := context.WithValue(r.Context(), "userIDContextKey", userID)

		// точка входа основного хендлера
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
