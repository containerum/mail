package routes

import (
	"net/http"

	"time"

	"strings"

	"math/rand"
	"strconv"

	"git.containerum.net/ch/grpc-proto-files/auth"
	"git.containerum.net/ch/grpc-proto-files/common"
	"git.containerum.net/ch/mail-templater/upstreams"
	"git.containerum.net/ch/user-manager/models"
	"git.containerum.net/ch/user-manager/utils"
	chutils "git.containerum.net/ch/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserCreateRequest struct {
	UserName  string `json:"username" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	Referral  string `json:"referral" binding:"url"`
	ReCaptcha string `json:"recaptcha" binding:"required"`
}

type UserCreateResponse struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	IsActive bool   `json:"is_active"`
}

type ActivateRequest struct {
	Link string `json:"link" binding:"required"`
}

type ResendLinkRequest struct {
	UserName string `json:"username" binding:"required,email"`
}

type UserInfoByIDGetResponse struct {
	Login string             `json:"login"`
	Data  models.ProfileData `json:"data"`
}

type BlacklistedUserEntry struct {
	Login string `json:"login"`
	ID    string `json:"id"`
}

type BlacklistGetResponse struct {
	BlacklistedUsers []BlacklistedUserEntry `json:"blacklist_users"`
}

type LinksGetResponse struct {
	Links []models.Link `json:"links"`
}

type UserInfoGetResponse struct {
	Login     string             `json:"login"`
	Data      models.ProfileData `json:"data"`
	ID        string             `json:"id"`
	IsActive  bool               `json:"is_active"`
	CreatedAt time.Time          `json:"created_at"`
}

type UserListEntry struct {
	ID            string             `json:"id"`
	Login         string             `json:"login"`
	Referral      string             `json:"referral"`
	Role          models.UserRole    `json:"role"`
	Access        string             `json:"access"`
	CreatedAt     time.Time          `json:"created_at"`
	DeletedAt     time.Time          `json:"deleted_at"`
	BlacklistedAt time.Time          `json:"blacklisted_at"`
	Data          models.ProfileData `json:"data"`
	IsActive      bool               `json:"is_active"`
	IsInBlacklist bool               `json:"is_in_blacklist"`
	IsDeleted     bool               `json:"is_deleted"`
}

type UserListGetResponse struct {
	Users []UserListEntry `json:"users"`
}

const (
	userNotFound            = "user %s was not found"
	userWithIDNotFound      = "user with id %s was not found"
	userAlreadyExists       = "user %s is already registered"
	userNotPartiallyDeleted = "user %s is not partially deleted"
	domainInBlacklist       = "email domain %s is in blacklist"
	linkNotFound            = "link %s was not found or already used or expired"
	profilesNotFound        = "profiles not found"
	waitForResend           = "can`t resend link now, please wait %d seconds"
)

func userCreateHandler(ctx *gin.Context) {
	var request UserCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, chutils.NewError(err.Error()))
		return
	}

	domain := strings.Split(request.UserName, "@")[1]
	blacklisted, err := svc.DB.IsDomainBlacklisted(domain)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if blacklisted {
		ctx.AbortWithStatusJSON(http.StatusForbidden, chutils.NewErrorF(domainInBlacklist, request.UserName))
		return
	}

	user, err := svc.DB.GetUserByLogin(request.UserName)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if user != nil {
		ctx.AbortWithStatusJSON(http.StatusConflict, chutils.NewErrorF(userAlreadyExists, request.UserName))
		return
	}

	salt := utils.GenSalt(request.UserName, request.UserName, request.UserName) // compatibility with old client db
	passwordHash := utils.GetKey(request.Password, salt)
	newUser := &models.User{
		Login:        request.UserName,
		PasswordHash: passwordHash,
		Salt:         salt,
		Role:         models.RoleUser,
		IsActive:     false,
		IsDeleted:    false,
	}

	var link *models.Link

	err = svc.DB.Transactional(func(tx *models.DB) error {
		if err := svc.DB.CreateUser(newUser); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return err
		}

		if err := svc.DB.CreateProfile(&models.Profile{
			User:      newUser,
			Referral:  request.Referral,
			Access:    "rw",
			CreatedAt: time.Now().UTC(),
		}); err != nil {
			return err
		}

		link, err = svc.DB.CreateLink(models.LinkTypeConfirm, 24*time.Hour, newUser)
		return err
	})

	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = svc.MailClient.SendConfirmationMail(&upstreams.Recipient{
		ID:        newUser.ID,
		Name:      request.UserName,
		Email:     request.UserName,
		Variables: map[string]string{"CONFIRM": link.Link},
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = svc.DB.Transactional(func(tx *models.DB) error {
		link.SentAt.Time = time.Now().UTC()
		link.SentAt.Valid = true
		return tx.UpdateLink(link)
	})

	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, UserCreateResponse{
		ID:       newUser.ID,
		Login:    newUser.Login,
		IsActive: newUser.IsActive,
	})
}

func linkResendHandler(ctx *gin.Context) {
	var request ResendLinkRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, chutils.NewError(err.Error()))
		return
	}

	user, err := svc.DB.GetUserByLogin(request.UserName)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if user == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, chutils.NewErrorF(userNotFound, user.Login))
		return
	}

	link, err := svc.DB.GetLinkForUser(models.LinkTypeConfirm, user)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if link == nil {
		link, err = svc.DB.CreateLink(models.LinkTypeConfirm, 24*time.Hour, user)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	if tdiff := time.Now().UTC().Sub(link.SentAt.Time); link.SentAt.Valid && tdiff < 5*time.Minute {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, chutils.NewErrorF(waitForResend, int(tdiff.Seconds())))
		return
	}

	err = svc.MailClient.SendConfirmationMail(&upstreams.Recipient{
		ID:        user.ID,
		Name:      request.UserName,
		Email:     request.UserName,
		Variables: map[string]string{"CONFIRM": link.Link},
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = svc.DB.Transactional(func(tx *models.DB) error {
		link.SentAt.Time = time.Now().UTC()
		link.SentAt.Valid = true
		return tx.UpdateLink(link)
	})
}

func activateHandler(ctx *gin.Context) {
	var request ActivateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, chutils.NewError(err.Error()))
		return
	}

	link, err := svc.DB.GetLinkFromString(request.Link)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if link == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, chutils.NewErrorF(linkNotFound, request.Link))
		return
	}

	// TODO: send request to billing manager

	err = svc.MailClient.SendActivationMail(&upstreams.Recipient{
		ID:        link.User.ID,
		Name:      link.User.Login,
		Email:     link.User.Login,
		Variables: map[string]string{},
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	tokens, err := svc.AuthClient.CreateToken(ctx, &auth.CreateTokenRequest{
		UserAgent:   ctx.Request.UserAgent(),
		UserId:      &common.UUID{Value: link.User.ID},
		UserIp:      ctx.ClientIP(),
		UserRole:    auth.Role_USER,
		RwAccess:    true,
		Access:      &auth.ResourcesAccess{},
		PartTokenId: nil,
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

func userToBlacklistHandler(ctx *gin.Context) {
	userID := ctx.GetHeader(UserIDHeader)
	user, err := svc.DB.GetUserByID(userID)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if user == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, chutils.NewErrorF(userWithIDNotFound, userID))
		return
	}

	profile, err := svc.DB.GetProfileByUser(user)
	if err != nil || profile == nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// TODO: send request to resource manager

	err = svc.MailClient.SendBlockedMail(&upstreams.Recipient{
		ID:    user.ID,
		Name:  user.Login,
		Email: user.Login,
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = svc.DB.BlacklistUser(user)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusAccepted)
}

func blacklistGetHandler(ctx *gin.Context) {
	blacklisted, err := svc.DB.GetBlacklistedUsers()
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var resp BlacklistGetResponse
	for _, v := range blacklisted {
		resp.BlacklistedUsers = append(resp.BlacklistedUsers, BlacklistedUserEntry{
			Login: v.Login,
			ID:    v.ID,
		})
	}
	ctx.JSON(http.StatusAccepted, resp)
}

func linksGetHandler(ctx *gin.Context) {
	userID := ctx.GetHeader(UserIDHeader)
	user, err := svc.DB.GetUserByID(userID)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if user == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, chutils.NewErrorF(userWithIDNotFound, userID))
	}

	links, err := svc.DB.GetUserLinks(user)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, LinksGetResponse{Links: links})
}

func userInfoGetHandler(ctx *gin.Context) {
	userID := ctx.GetHeader(UserIDHeader)
	user, err := svc.DB.GetUserByID(userID)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if user == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, chutils.NewErrorF(userWithIDNotFound, userID))
		return
	}

	profile, err := svc.DB.GetProfileByUser(user)
	if err != nil || profile == nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, &UserInfoGetResponse{
		Login:     user.Login,
		Data:      profile.Data,
		ID:        user.ID,
		IsActive:  user.IsActive,
		CreatedAt: profile.CreatedAt,
	})
}

func userInfoUpdateHandler(ctx *gin.Context) {
	userID := ctx.GetHeader(UserIDHeader)
	user, err := svc.DB.GetUserByID(userID)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if user == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, chutils.NewErrorF(userWithIDNotFound, userID))
		return
	}

	profile, err := svc.DB.GetProfileByUser(user)
	if err != nil || profile == nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := ctx.ShouldBindWith(&profile.Data, binding.JSON); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, chutils.NewError(err.Error()))
		return
	}

	err = svc.DB.Transactional(func(tx *models.DB) error {
		return tx.UpdateProfile(profile)
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, &UserInfoGetResponse{
		Login:     user.Login,
		Data:      profile.Data,
		ID:        user.ID,
		IsActive:  user.IsActive,
		CreatedAt: profile.CreatedAt,
	})
}

func userListGetHandler(ctx *gin.Context) {
	profiles, err := svc.DB.GetAllProfiles()
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if profiles == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, chutils.NewError(profilesNotFound))
		return
	}

	filters := strings.Split(ctx.Query("filters"), ",")
	var filterFuncs []func(p models.Profile) bool
	for _, filter := range filters {
		switch filter {
		case "active":
			filterFuncs = append(filterFuncs, func(p models.Profile) bool {
				return p.User.IsActive
			})
		case "inactive":
			filterFuncs = append(filterFuncs, func(p models.Profile) bool {
				return !p.User.IsActive
			})
		case "in_blacklist":
			filterFuncs = append(filterFuncs, func(p models.Profile) bool {
				return p.User.IsInBlacklist
			})
		case "deleted":
			filterFuncs = append(filterFuncs, func(p models.Profile) bool {
				return p.User.IsDeleted
			})
		}
	}

	satisfiesFilter := func(p models.Profile) bool {
		ret := true
		for _, v := range filterFuncs {
			ret = ret && v(p)
		}
		return ret
	}

	var resp UserListGetResponse
	for _, v := range profiles {
		if !satisfiesFilter(v) {
			continue
		}
		resp.Users = append(resp.Users, UserListEntry{
			ID:            v.User.ID,
			Login:         v.User.Login,
			Referral:      v.Referral,
			Role:          v.User.Role,
			Access:        v.Access,
			CreatedAt:     v.CreatedAt,
			DeletedAt:     v.DeletedAt.Time,
			BlacklistedAt: v.BlacklistAt.Time,
			Data:          v.Data,
			IsActive:      v.User.IsActive,
			IsInBlacklist: v.User.IsInBlacklist,
			IsDeleted:     v.User.IsDeleted,
		})
	}
	ctx.JSON(http.StatusOK, resp)
}

func partialDeleteHandler(ctx *gin.Context) {
	userID := ctx.GetHeader(UserIDHeader)
	user, err := svc.DB.GetUserByID(userID)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if user == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, chutils.NewErrorF(userWithIDNotFound, userID))
		return
	}

	// TODO: send request to user manager

	// TODO: send request to billing manager

	if _, err := svc.AuthClient.DeleteUserTokens(ctx, &auth.DeleteUserTokensRequest{
		UserId: &common.UUID{Value: user.ID},
	}); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = svc.DB.Transactional(func(tx *models.DB) error {
		user.IsDeleted = true
		return tx.UpdateUser(user)
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusAccepted)
}

func completeDeleteHandler(ctx *gin.Context) {
	userID := ctx.GetHeader(UserIDHeader)
	user, err := svc.DB.GetUserByID(userID)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if user == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, chutils.NewErrorF(userWithIDNotFound, userID))
		return
	}

	if !user.IsDeleted {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, chutils.NewErrorF(userNotPartiallyDeleted, user.Login))
		return
	}

	// TODO: send request to billing manager

	err = svc.DB.Transactional(func(tx *models.DB) error {
		user.Login = user.Login + strconv.Itoa(rand.Int())
		return svc.DB.UpdateUser(user)
	})
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusAccepted)
}
