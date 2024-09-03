package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirin7/go_final_project/constants"
	"github.com/sirin7/go_final_project/helpers"
)

// Получение следующей даты задачи
func (h *Handler) GetDateTask(w http.ResponseWriter, r *http.Request) {
	nowstr := r.URL.Query().Get("now")
	if nowstr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "now date missing"}`))
		return
	}

	now, err := time.Parse(constants.FormatDate, nowstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid now date format"}`))
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "date missing"}`))
		return
	}

	repeat := r.URL.Query().Get("repeat")
	if repeat == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "repeat missing"}`))
		return
	}

	nextDate, err := helpers.NextDate(now, date, repeat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	// Возвращаем следующую дату
	w.Header().Set("Content-Type", "application/json, charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%s", nextDate)))
}
