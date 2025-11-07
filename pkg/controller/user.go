package controller

import (
	"encoding/json"
	"net/http"
	"time"

	tokencache "github.com/go-auth-microservice/pkg/model/tokenCache"
	usermodel "github.com/go-auth-microservice/pkg/model/userModel"
	"github.com/go-auth-microservice/pkg/utils/logger"
)

func GetUserData(w http.ResponseWriter, r *http.Request) {
	log := logger.InitializeAuditLogger()
	ctx := r.Context()
	userId, _ := ctx.Value("userId").(uint64)
	userData, err := usermodel.FindUserByID(userId)
	if err != nil {
		log.Errorf("unable to find user with ID %v ", userId, err)
		return
	}
	json.NewEncoder(w).Encode(userData)
}

func DeActivateUser(w http.ResponseWriter, r *http.Request) {
	log := logger.InitializeAuditLogger()
	ctx := r.Context()
	userId, _ := ctx.Value("userId").(uint64)
	var userData usermodel.UserStatus
	userData, err := usermodel.FindUserByID(userId)
	if err != nil {
		log.Errorf("unable to find user with ID %v ", userId, err)
		return
	}
	err = userData.Disable()
	if err != nil {
		http.Error(w, "user has already been disabled", http.StatusBadRequest)
		log.Error(err)
		return
	}
	err = userData.Save()
	if err != nil {
		http.Error(w, "unable to update user status", http.StatusBadRequest)
		log.Error("unable to update user status ", err)
		return
	}
	var blackListedToken tokencache.BlackListedToken = tokencache.GetBlacklistTokenCache()
	token := r.Header.Get("Authorization")
	blackListedToken.Set(token, time.Now().Add(time.Minute*time.Duration(5)).Unix())
	w.Write([]byte("user has been disabled"))
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	log := logger.InitializeAuditLogger()
	ctx := r.Context()
	userId, _ := ctx.Value("userId").(uint64)
	var data map[string]interface{}
	json.NewDecoder(r.Body).Decode(&data)
	password, ok := data["password"].(string)
	if !ok {
		log.Error("invalid request to change password for user ID ", userId)
		http.Error(w, "new password should be 8 character long", http.StatusBadRequest)
		return
	}
	if password == "" || len(password) < 8 {
		log.Error("invalid request to change password for user ID ", userId)
		http.Error(w, "new password should be 8 character long", http.StatusBadRequest)
		return
	}
	var userData usermodel.UserSignUp
	userData, err := usermodel.FindUserByID(userId)
	if err != nil {
		log.Errorf("unable to find user with ID %v ", userId, err)
		return
	}
	err = userData.SetPassword(password)
	if err != nil {
		http.Error(w, "unable to change password", http.StatusInternalServerError)
		log.Error("unable to change password for ID ", userId, " ", err)
		return
	}
	err = userData.Save()
	if err != nil {
		http.Error(w, "unable to update user password", http.StatusInternalServerError)
		log.Error("unable to update user password for ID ", userId, " ", err)
		return
	}
	var blackListedToken tokencache.BlackListedToken = tokencache.GetBlacklistTokenCache()
	token := r.Header.Get("Authorization")
	blackListedToken.Set(token, time.Now().Add(time.Minute*time.Duration(5)).Unix())
	w.Write([]byte("user password has been changed."))
}
