package logger

import (
	"bytes"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// структура для хранения сведений об ответе
type responseData struct {
	status      int
	size        int
	contentType string
	body        *bytes.Buffer
}

// структура с http.ResponseWriter и responseData
type loggingResponseWriter struct {
	// встраиваем оригинальный http.ResponseWriter
	http.ResponseWriter
	responseData *responseData
}

// переопределяем методы Write и WriteHeader интерфейса http.ResponseWriter
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter, получаем размер
	size, err := r.ResponseWriter.Write(b)
	// захватваем размер
	r.responseData.size += size
	// захватываем тело ответа
	r.responseData.body.Write(b)

	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	// захватываем код статуса
	r.responseData.status = statusCode
	// захватываем тип контента
	r.responseData.contentType = r.Header().Get("Content-Type")
}

// WithLogging - middleware, Функция-оболочка, которая оборачивает http.Handler
// добавляет дополнительный код и возвращает новый http.Handler
func WithLogging(lg *logrus.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// для определения времени обработки запроса
			start := time.Now()

			// создаём экземпляр структуры responseData
			responseData := &responseData{
				status:      0,
				size:        0,
				contentType: "",
				body:        bytes.NewBuffer(nil),
			}

			// создаём экземпляр структуры loggingResponseWriter
			lw := loggingResponseWriter{
				// встраиваем оригинальный http.responseWriter
				ResponseWriter: w,
				responseData:   responseData,
			}

			// записываем данные запроса
			// эндпоинт
			reqURI := r.RequestURI
			// метод запроса
			reqMethod := r.Method
			// тип контента
			reqContentType := r.Header.Get("Content-Type")

			// обслуживание оригинального запроса c внедрённой реализацией http.ResponseWriter
			h.ServeHTTP(&lw, r)

			// записываем данные ответа
			// статус
			resStatus := responseData.status
			// тип контента
			resContentType := responseData.contentType
			// размер
			resSize := responseData.size
			// тело
			resBody := responseData.body.String()

			// время обработки запроса
			duration := time.Since(start)

			// отправляем сведения в логгер
			lg.Debugf("---------------------------------------------------------------")
			lg.Debugf("Request method: %v", reqMethod)
			lg.Debugf("Request URI: %v", reqURI)
			lg.Debugf("Request Content-Type: %v", reqContentType)
			lg.Debugf("Response status: %v", resStatus)
			lg.Debugf("Response Content-Type: %v", resContentType)
			lg.Debugf("Response size: %v", resSize)
			lg.Debugf("Response body: %v", resBody)
			lg.Debugf("Duration: %v", duration)
		})
	}
}
