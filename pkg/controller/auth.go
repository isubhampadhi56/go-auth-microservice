package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	usermodel "github.com/go-auth-microservice/pkg/model/userModel"
	jwtauth "github.com/go-auth-microservice/pkg/utils/jwtAuth"
	"github.com/go-auth-microservice/pkg/utils/logger"
	"github.com/go-auth-microservice/pkg/utils/validation"
	"github.com/golang-jwt/jwt/v5"
)

type userSignup struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var user userSignup
	log := logger.InitializeAuditLogger()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Error("invalid request body")
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := validation.Validator.Struct(user); err != nil {
		log.Error("data validation failed")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	var userData usermodel.UserSignUp = usermodel.CreateUser(user.Email)
	userData.SetPassword(user.Password)
	if err := userData.Save(); err != nil {
		log.Error("unable to create user")
		http.Error(w, "email already exist", http.StatusConflict)
		return
	}
	json.NewEncoder(w).Encode(userData)
	log.Info("user has been created", userData)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user userSignup
	log := logger.InitializeAuditLogger()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Error("invalid request body")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := validation.Validator.Struct(user); err != nil {
		log.Error("data validation failed")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	var userData usermodel.UserLogin
	userData, err := usermodel.FindUserByEmail(user.Email)
	if err != nil {
		log.Error(`user with email "%v" not found on DB`, user.Email, err)
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}
	err = userData.ValidatePassword(user.Password)
	if err != nil {
		log.Error("invalid login for user %v", user.Email)
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}
	if !userData.GetUserStatus() {
		http.Error(w, "user has been disabled plase contact admin", http.StatusUnauthorized)
		log.Errorf("user %d has been disabled plase contact admin ", userData.GetUserID())
		return
	}
	claims := jwt.MapClaims{}
	claims["userId"] = userData.GetUserID()
	accessToken, err := jwtauth.GetAccessTokenHandler().CreateToken(claims)
	if err != nil {
		log.Error("error creating token ", err)
		return
	}
	refreshToken, err := jwtauth.GetRefreshTokenHandler().CreateToken(claims)
	if err != nil {
		log.Error("error creating token ", err)
		return
	}
	res := map[string]interface{}{}
	res["accesstoken"] = accessToken
	res["refreshtoken"] = refreshToken
	json.NewEncoder(w).Encode(res)
	log.Infof("user with ID %v and Email %v has loggedIn successfully", userData.GetUserID(), user.Email)
}

func RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	log := logger.InitializeAuditLogger()
	refreshToken := r.Header.Get("RefreshToken")
	claim, err := jwtauth.GetRefreshTokenHandler().VerifyToken(refreshToken)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		log.Error("invalid Refresh Token", err)
		return
	}
	tokenCreatedAt, _ := claim["iat"].(float64)
	userId, ok := claim["userId"].(float64)
	if !ok {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		log.Error("invalid or missing userId in Refresh Token")
		return
	}
	var userData usermodel.UserLogin
	userData, err = usermodel.FindUserByID(uint64(userId))
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		log.Error("refresh token denied")
		return
	}
	if !userData.GetUserStatus() {
		http.Error(w, "user has been disabled plase contact admin", http.StatusUnauthorized)
		log.Errorf("user %d has been disabled plase contact admin ", userData.GetUserID())
		return
	}
	if int64(tokenCreatedAt) < userData.GetUserLastUpdated().Unix() {
		http.Error(w, "user has been updated please relogin", http.StatusUnauthorized)
		log.Errorf("user %d has been updated please relogin", userData.GetUserID())
		return
	}
	claims := jwt.MapClaims{}
	claims["userId"] = userData.GetUserID()
	accessToken, err := jwtauth.GetAccessTokenHandler().CreateToken(claims)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		log.Error("failed to generate token", err)
	}
	res := map[string]interface{}{}
	res["accesstoken"] = accessToken
	json.NewEncoder(w).Encode(res)
	log.Info("refresh Token has been generated for user ID %v", userData.GetUserID())
}

func CheckIfSessionValid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, _ := ctx.Value("userId").(uint64)
	w.Write([]byte("user auth is valid for ID " + strconv.Itoa(int(userId))))

}
