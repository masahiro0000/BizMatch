package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/masahiro0000/BizMatch/internal/domain"
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
		case errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrIncorrectPassword):
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
	// Retrieve the current user from the session.
	currentUser, err := util.GetCurrentUser(c)
	if err != nil {
		log.Printf("Fail to get current user: %v", err)
		c.Redirect(http.StatusFound, "/users/login")
		return
	}

	// Get the candidate users for the current user.
	candidates, err := h.userUsecase.GetRecommendedUsers(currentUser.ID)
	if err != nil {
		log.Printf("Fail to get candidate users. error:%v", err)
		c.HTML(http.StatusInternalServerError, "mypage.html", gin.H{
			"user": currentUser,
		})
	}
	// Score the candidate users based on the current user's preferences.
	var scoredUser = h.userUsecase.ScoreUsers(currentUser, candidates)
	c.HTML(http.StatusOK, "mypage.html", gin.H{
		"user": currentUser,
		"recommendedUsers": scoredUser,
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

func (h *UserHandler) ShowChangePasswordForm(c *gin.Context) {
	// Retrieve the current user from the session.
	currentUser, err := util.GetCurrentUser(c)
	if err != nil {
		c.Redirect(http.StatusFound, "/users/login")
		return
	}

	c.HTML(http.StatusOK, "change_password.html", gin.H{
		"user": currentUser,
	})
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	// Retrieve the current user from the session.
	currentUser, err := util.GetCurrentUser(c)
	if err != nil {
		c.Redirect(http.StatusFound, "/users/login")
		return
	}

	// Extract the old password, new password, and password confirmation from the POST form.
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")
	confirmPassword := c.PostForm("confirm_password")

	// Attempt to update the user's password using the provided information.
	err = h.userUsecase.UpdatePassword(currentUser.ID, oldPassword, newPassword, confirmPassword)
	if err != nil {
		status := http.StatusInternalServerError

		// If the error is related to the new password being too short, new password mismatches confirmation password or
		// the provided old password mismatches the user's current password, use a bad request status.
		if errors.Is(err, domain.ErrPasswordTooShort) ||
			errors.Is(err, domain.ErrNewPasswordMismatch) ||
			errors.Is(err, domain.ErrOldPasswordMismatch){
			status = http.StatusBadRequest
		}

		c.HTML(status, "change_password.html", gin.H{
			"user": currentUser,
			"error": err.Error(),
		})
		return
	}

	// On successful password update, redirect the user to their mypage.
	c.Redirect(http.StatusFound, "/users/mypage")
}

func (h *UserHandler) ShowSearchUsersForm(c *gin.Context) {
	// Retrieve the current user from the session.
	currentUser, err := util.GetCurrentUser(c)
	if err != nil {
		c.Redirect(http.StatusFound, "/users/login")
		return
	}

	// Generate a list of ages and gender for display.
	ageGroups := []int{20, 30, 40, 50}
	genders := util.GenerateGenderList()

	c.HTML(http.StatusOK, "search_users.html", gin.H{
		"user": currentUser,
		"ageGroups": ageGroups,
		"genders": genders,
	})
}

func (h *UserHandler) SearchUsers(c *gin.Context) {
	// Retrieve the current user from the session.
	currentUser, err := util.GetCurrentUser(c)
	if err != nil {
		c.Redirect(http.StatusFound, "/users/login")
		return
	}

	// Retrieve search criteria from the form input as arrays of strings.
	prefectureStrs := c.PostFormArray("prefectures[]")
	industryStrs := c.PostFormArray("industries[]")
	jobStrs := c.PostFormArray("jobs[]")
	positionStrs := c.PostFormArray("positions[]")
	genderStrs := c.PostFormArray("genders[]")
	ageGroupsStr := c.PostFormArray("age_groups[]")

	// Convert age group strings to integers.
	var ageGroups []int
	for _, ag := range ageGroupsStr {
		v, err := strconv.Atoi(ag)
		if err == nil {
			ageGroups = append(ageGroups, v)
		}
	}

	var filter domain.UserSearchFilter

	// Parse the IDs from string array to integer array.
	p, err := util.ParseNullableInts(prefectureStrs)
	if err != nil {
		c.String(http.StatusBadRequest, "都道府県IDに不正な値があります: %v", err)
		return
	}
	filter.Prefectures = p

	i, err := util.ParseNullableInts(industryStrs)
	if err != nil {
		c.String(http.StatusBadRequest, "業種IDに不正な値があります: %v", err)
		return
	}
	filter.Industries = i

	j, err := util.ParseNullableInts(jobStrs)
	if err != nil {
		c.String(http.StatusBadRequest, "職種IDに不正な値があります: %v", err)
		return
	}
	filter.Jobs = j

	pos, err := util.ParseNullableInts(positionStrs)
	if err != nil {
		c.String(http.StatusBadRequest, "役職IDに不正な値があります: %v", err)
		return
	}
	filter.Positions = pos

	// Assign the age groups and gender filter.
	filter.AgeGroups = ageGroups
	filter.Genders = util.ParseStringsAsPtrs(genderStrs)

	// Exclude the current user from the search result.
	filter.ExcludeUserID = &currentUser.ID

	// Perform the user search using the specified filters.
	users, err := h.userUsecase.SearchUsers(&filter)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "search_users.html", gin.H{
			"error": err,
		})
		return
	}

	c.HTML(http.StatusOK, "search_result.html", gin.H{
		"currentUser": currentUser,
		"users": users,
	})
}

func (h *UserHandler) ShowUserDetail(c *gin.Context) {
	// Retrieve the current user from the session.
	currentUser, err := util.GetCurrentUser(c)
	if err != nil {
		c.Redirect(http.StatusFound, "/users/login")
		return
	}

	// Get the "username" parameter from the URL.
	username := c.Param("username")
	if username == "" {
		c.String(http.StatusBadRequest, "ユーザーIDが取得できませんでした")
		return
	}

	// Retrieve the user using provided username.
	user, err := h.userUsecase.GetUserByUsername(username)
	if err != nil {
		c.String(http.StatusBadRequest, "ユーザーが見つかりません")
		return
	}

	c.HTML(http.StatusOK, "user_detail.html", gin.H{
		"currentUser": currentUser,
		"user": user,
	})
}