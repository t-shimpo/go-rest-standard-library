package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/t-shimpo/go-rest-standard-library/models"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

func respondWithJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJson(w, status, map[string]string{"error": message})
}

// `GET /users`
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // デフォルト値
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0 // デフォルト値
	}

	users, err := models.GetUsers(limit, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "ユーザー取得中にエラーが発生しました")
		return
	}

	respondWithJson(w, http.StatusOK, users)
}

// `POST /users`
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	var req CreateUserRequest
	defer r.Body.Close()

	// json デコード
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "無効なリクエストボディ")
		return
	}

	// バリデーション
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)

	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "名前は必須です")
		return
	}

	if req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "メールは必須です")
		return
	}

	// DB に保存
	user, err := models.CreateUser(req.Name, req.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "ユーザー作成中にエラーが発生しました")
		return
	}

	respondWithJson(w, http.StatusCreated, user)
}

// `GET /users/{id}`
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// URL から ID を取得
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "IDが必要です")
		return
	}

	// ID を整数に変換
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "IDは数値である必要があります")
		return
	}

	// DB からユーザー取得
	user, err := models.GetUserByID(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			respondWithError(w, http.StatusNotFound, "ユーザーが見つかりません")
		} else {
			respondWithError(w, http.StatusInternalServerError, "ユーザー取得中にエラーが発生しました")
		}
		return
	}

	respondWithJson(w, http.StatusOK, user)
}

// `PATCH /users/{id}“
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// URL から ID を取得
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "IDが必要です")
		return
	}

	// ID を整数に変換
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "IDは数値である必要があります")
		return
	}

	var req UpdateUserRequest
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "無効なリクエストボディ")
		return
	}

	if req.Name == nil && req.Email == nil {
		respondWithError(w, http.StatusBadRequest, "更新するフィールドを指定してください")
		return
	}

	user, err := models.UpdateUser(id, req.Name, req.Email)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			respondWithError(w, http.StatusNotFound, "ユーザーが見つかりません")
		} else {
			respondWithError(w, http.StatusInternalServerError, "ユーザー更新中にエラーが発生しました")
		}
		return
	}

	respondWithJson(w, http.StatusOK, user)
}

// `DELETE /users/{id}`
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// URL から ID を取得
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "IDが必要です")
		return
	}

	// ID を整数に変換
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "IDは数値である必要があります")
		return
	}

	err = models.DeleteUser(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			respondWithError(w, http.StatusNotFound, "ユーザーが見つかりません")
		} else {
			respondWithError(w, http.StatusInternalServerError, "ユーザー削除中にエラーが発生しました")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
