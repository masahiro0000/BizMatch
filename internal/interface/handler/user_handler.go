package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/domain"
	"github.com/masahiro0000/BizMatch/internal/infrastructure/repository"
	"github.com/masahiro0000/BizMatch/internal/interface/handler/util"
	"github.com/masahiro0000/BizMatch/internal/usecase"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler (u *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: u,
	}
}

func (h *UserHandler) ShowSignupForm(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{
	})
}

// Signup processes the signup form submission.
func (h *UserHandler) Signup(c *gin.Context) {
	username := c.PostForm("username")
	displayName := c.PostForm("displayName")
	password := c.PostForm("password")

	err := h.userUsecase.Signup(c, username, displayName, password)
	if err != nil {
		var status int
		// Determine the HTTP status code based on the error type.
		switch {
		case errors.Is(err, domain.ErrInvalidInput) || errors.Is(err, domain.ErrPasswordTooShort):
			status = http.StatusBadRequest
		case errors.Is(err, domain.ErrUserAlreadyExists):
			status = http.StatusConflict
		default:
			status = http.StatusInternalServerError
		}

		// Re-render the signup page with an error message and pre-filled form values.
		c.HTML(status, "signup.html", gin.H{
			"error": err.Error(),
			"username": username,
			"displayName": displayName,
		})
		return
	}

	// On successful signup, redirect the user to the registration info page.
	c.Redirect(http.StatusFound, "/users/register_info")
}

func (h *UserHandler) ShowLoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Attempt to login using username and password.
	user, err := h.userUsecase.Login(c, username, password)
	if err != nil {
		var status int
		// Adjust the HTTP status code according to the type of error.
		switch {
		case errors.Is(err, repository.ErrUserNotFound) || errors.Is(err, domain.ErrIncorrectPassword):
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}

		c.HTML(status, "login.html", gin.H{
			"error": err.Error(),
			"username": username,
		})
	}

	if user != nil {
		// Create a new session and store the authenticated user data.
		session := sessions.Default(c)
		session.Set("user", user)

		// Save the session.
		if err := session.Save(); err != nil {
			c.String(http.StatusInternalServerError, "セッション保存エラーが発生しました")
		}
	}

	// If successful, redirect the user to their mypage.
	c.Redirect(http.StatusFound, "/users/mypage")
}

func (h *UserHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.String(http.StatusInternalServerError, "セッションの削除に失敗しました")
		return
	}

	c.Redirect(http.StatusFound, "users/login")
}

func (h *UserHandler) ShowMypage(c *gin.Context) {
	session := sessions.Default(c)

	// Get the user data stored in the session.
	user := session.Get("user")
	currentUser, ok := user.(*domain.User)

	// If no user is found in the session, redirect to the login page.
	if !ok || currentUser == nil {
		c.Redirect(http.StatusFound, "/users/login")
		return
	}

	c.HTML(http.StatusOK, "mypage.html", gin.H{
		"user": currentUser,
	})
}

func (h *UserHandler) ShowRegisterInfo(c *gin.Context) {
	// Retrieve the current user from the session.
	currentUser, err := util.GetCurrentUser(c)
	if err != nil {
		log.Printf("Fail to get current user: %v", err)
		c.Redirect(http.StatusFound, "/users/login")
		return
	}

	// Update userAge with the currentUser's age.
	var userAge int64 = 0
	if currentUser.Age != nil {
		userAge = *currentUser.Age
	}

	// Update userGender with the currentUser's gender.
	var userGender string
	if currentUser.Gender != nil {
		userGender = *currentUser.Gender
	}

	// Generate a list of ages and gender for display.
	ages := util.GenerateAgeList()
	gender := util.GenerateGenderList()

	c.HTML(http.StatusOK, "register_info.html", gin.H{
		"user": currentUser,
		"userAge": userAge,
		"userGender": userGender,
		"ages": ages,
		"gender": gender,
	})
}

func (h *UserHandler) RegisterInfo(c *gin.Context) {
	// Retrieve the current user from the session.
	currentUser, err := util.GetCurrentUser(c)
	if err != nil {
		c.Redirect(http.StatusFound, "/users/login")
		return
	}

	// Extract form values from the POST request.
	username := c.PostForm("username")
	displayName := c.PostForm("displayName")
	prefectureStr := c.PostForm("prefecture")
	industryStr := c.PostForm("industry")
	jobStr := c.PostForm("job")
	positionStr := c.PostForm("position")
	ageStr := c.PostForm("age")
	gender := c.PostForm("gender")
	file, fileErr := c.FormFile("photo")
	profileDescription := c.PostForm("profileDescription")

	// Parse the string value into nullable integer.
	prefecture, errPrefecture := util.ParseNullableInt(prefectureStr)
	industry, errIndustry := util.ParseNullableInt(industryStr)
	job, errJob := util.ParseNullableInt(jobStr)
	position, errPosition := util.ParseNullableInt(positionStr)
	age, errAge := util.ParseNullableInt(ageStr)

	// Convert the gender and profileDescription string into non-nil pointer.
	gender = *util.StringPtr(gender)
	profileDescription = *util.StringPtr(profileDescription)

	var photoURL string
	// If a file is uploaded, this will be set to the path of the saved image.
	if fileErr == nil {
		fileName := fmt.Sprintf("%d_%s", currentUser.ID, file.Filename)
		savePath := filepath.Join("web", "static", fileName)

		// Attempt to save the uploaded file.
		if saveErr := c.SaveUploadedFile(file, savePath); saveErr != nil {
			log.Printf("Failed to save file: %v", saveErr)
		} else {
			photoURL = "/static/" + fileName
		}
	}

	var photo *string
	if photoURL != "" {
		// If a new photo was uploaded and save, update the current user's photo
		photo = util.StringPtr(photoURL)
		currentUser.Photo = photo
	} else {
		// Otherwise, retain the current user's existing photo.
		photo = currentUser.Photo
	}

	// Check for errors in parsing any of the numeric form variable.
	if errPrefecture != nil || errIndustry != nil || errJob != nil || errPosition != nil || errAge != nil {
		log.Printf("Fail to parse form input. user ID:%d, prefecture:%v, industry:%v, job:%v, position:%v, age:%d",
			currentUser.ID, errPrefecture, errIndustry, errJob, errPosition, errAge)
		c.HTML(http.StatusBadRequest, "register_info.html", gin.H{
			"user": currentUser,
			"ages": util.GenerateAgeList(),
			"gender": util.GenerateGenderList(),
		})
		return
	}

	// Create a new user object with the updated registration information.
	user := &domain.User{
		ID: 		 		currentUser.ID,
		Username: 	 		username,
		DisplayName: 		displayName,
		Prefecture:  		prefecture,
		Industry: 	 		industry,
		Job: 		 		job,
		Position: 	 		position,
		Age:		 		age,
		Gender:		 		&gender,
		Photo:				photo,
		ProfileDescription: &profileDescription,
	}

	// Attempt to update the user's registration info.
	if err := h.userUsecase.RegisterInfo(user); err != nil {
		log.Printf("Fail to update user info in usecase. user ID:%d, error:%s", currentUser.ID, err)
		c.HTML(http.StatusInternalServerError, "register_info.html", gin.H{
			"user": user,
			"ages": util.GenerateAgeList(),
			"gender": util.GenerateGenderList(),
		})
		return
	}

	// Update the session with the new name information.
	util.SetCurrentUser(c, user)
	c.Redirect(http.StatusFound, "/users/mypage")
}