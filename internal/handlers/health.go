package handlers

import (
	"encoding/json"
	"net/http"
	"sync/atomic" // имей ввиду она не супер быстрая (если я правильно помню), если у тебя атомиков будет по коду то будет не очень все хорошо
	"time"
)

type ResponseInfo struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
}

func (router *RouterStruct) livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // не ну в целом мне нравится что мы живы всегда, даже когда нас дудосят или мы еще не прошли все инициализации
}

// Readiness: сервис готов только после успешного старта
func (router *RouterStruct) readinessHandler(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&router.ready) == 1 { // к тому же, нафига тебе тут инт, есть же атомик булеан? Она легче и быстрее и лучше.
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Service not ready", http.StatusServiceUnavailable)
	}
}

func (router *RouterStruct) infoHandler(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&router.ready) == 1 { // вчитался, по жопе надо за такое бить, у тебя структура одна? RouterStruct, зачем ты тогда её разбил по файлам?
		// Хочется читаемости? Сделай отдельный объект который уже в основном handler своем добавишь, 
		// а то получается я если не посмотрю все файлы, то не узнаю что тут ты залезаешь внутрь router.ready
		// получается как будто ты лезешь в чужой объект из этого файла (хоть я и понимаю что это не так)
		info := ResponseInfo{
			Version:   "1.0.0",
			Timestamp: time.Now().Format(time.RFC3339),
			Status:    "running",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	} else {
		http.Error(w, "Service not ready", http.StatusServiceUnavailable)
	}

}
