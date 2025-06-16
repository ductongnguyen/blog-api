package http_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ductong169z/shorten-url/internal/auth"
	authhttp "github.com/ductong169z/shorten-url/internal/auth/delivery/http"

	"github.com/ductong169z/shorten-url/config"
	mock "github.com/ductong169z/shorten-url/internal/auth/mocks"
	"github.com/ductong169z/shorten-url/internal/models"

	"github.com/ductong169z/shorten-url/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandlers_Register(t *testing.T) {
	type mockUseCase struct {
		expCall bool
		input   *models.User
		output  models.User
		err     error
	}
	tcs := map[string]struct {
		givenInput  string
		mockUseCase mockUseCase
		expBody     string
		expErr      error
		expCode     int
	}{
		"success": {
			givenInput: `{
				"username": "test",
				"email": "test@email.com",
				"password": "pass",
				"role": "user"
			}`,
			mockUseCase: mockUseCase{
				expCall: true,
				input: &models.User{
					Username: "test",
					Email:    "test@email.com",
					Password: "pass",
					Role:     models.RoleUser,
				},
				output: models.User{
					ID:       1,
					Username: "test",
					Email:    "test@email.com",
					Password: "pass",
					Role:     models.RoleUser,
				},
				err: nil,
			},
			expBody: `{"message":"Success","result":{"id":1,"username":"test","email":"test@email.com","role":"user","created_at":"0001-01-01 00:00:00","updated_at":"0001-01-01 00:00:00"}}`,

			expErr:  nil,
			expCode: http.StatusOK,
		},
		"invalid input": {
			givenInput: `{"username": "", "email": "invalid", "password": "", "role": "user"}`,
			mockUseCase: mockUseCase{
				expCall: false,
			},
			expBody: `{"message":"Internal server error"}`,
			expErr:  nil,
			expCode: http.StatusInternalServerError,
		},
		"usecase error": {
			givenInput: `{
				"username": "test2",
				"email": "test2@email.com",
				"password": "pass2",
				"role": "user"
			}`,
			mockUseCase: mockUseCase{
				expCall: true,
				input: &models.User{
					Username: "test2",
					Email:    "test2@email.com",
					Password: "pass2",
					Role:     models.RoleUser,
				},
				err: assert.AnError,
			},
			expBody: `{"message":"Internal server error"}`,
			expErr:  nil,
			expCode: http.StatusInternalServerError,
		},
		"duplicate email": {
			givenInput: `{
				"username": "test3",
				"email": "duplicate@email.com",
				"password": "pass3",
				"role": "user"
			}`,
			mockUseCase: mockUseCase{
				expCall: true,
				input: &models.User{
					Username: "test3",
					Email:    "duplicate@email.com",
					Password: "pass3",
					Role:     models.RoleUser,
				},
				err: errors.New("email already exists"),
			},
			expBody: `{"message":"Internal server error"}`,
			expErr:  nil,
			expCode: http.StatusInternalServerError,
		},
		"duplicate username": {
			givenInput: `{
				"username": "duplicate",
				"email": "duplicate@email.com",
				"password": "pass4",
				"role": "user"
			}`,
			mockUseCase: mockUseCase{
				expCall: true,
				input: &models.User{
					Username: "duplicate",
					Email:    "duplicate@email.com",
					Password: "pass4",
					Role:     models.RoleUser,
				},
				err: errors.New("username already exists"),
			},
			expBody: `{"message":"Internal server error"}`,
			expErr:  nil,
			expCode: http.StatusInternalServerError,
		},
		"invalid email": {
			givenInput: `{
				"username": "test4",
				"email": "invalid",
				"password": "pass4",
				"role": "user"
			}`,
			mockUseCase: mockUseCase{
				expCall: false,
			},
			expBody: `{"message":"Internal server error"}`,
			expErr:  nil,
			expCode: http.StatusInternalServerError,
		},
		"empty username": {
			givenInput:  `{"username": "", "email": "emptyuser@email.com", "password": "pass", "role": "user"}`,
			mockUseCase: mockUseCase{expCall: false},
			expBody:     `{"message":"Internal server error"}`,
			expErr:      nil,
			expCode:     http.StatusInternalServerError,
		},
		"empty password": {
			givenInput:  `{"username": "test5", "email": "emptypass@email.com", "password": "", "role": "user"}`,
			mockUseCase: mockUseCase{expCall: false},
			expBody:     `{"message":"Internal server error"}`,
			expErr:      nil,
			expCode:     http.StatusInternalServerError,
		},
		"empty email": {
			givenInput:  `{"username": "test6", "email": "", "password": "pass6", "role": "user"}`,
			mockUseCase: mockUseCase{expCall: false},
			expBody:     `{"message":"Internal server error"}`,
			expErr:      nil,
			expCode:     http.StatusInternalServerError,
		},
		"empty role": {
			givenInput:  `{"username": "test7", "email": "role@email.com", "password": "pass7", "role": ""}`,
			mockUseCase: mockUseCase{expCall: false},
			expBody:     `{"message":"Internal server error"}`,
			expErr:      nil,
			expCode:     http.StatusInternalServerError,
		},
		"invalid role": {
			givenInput:  `{"username": "test8", "email": "invrole@email.com", "password": "pass8", "role": "superuser"}`,
			mockUseCase: mockUseCase{expCall: false},
			expBody:     `{"message":"Internal server error"}`,
			expErr:      nil,
			expCode:     http.StatusInternalServerError,
		},
		"admin role": {
			givenInput: `{"username": "admin1", "email": "admin@email.com", "password": "adminpass", "role": "admin"}`,
			mockUseCase: mockUseCase{
				expCall: true,
				input: &models.User{
					Username: "admin1",
					Email:    "admin@email.com",
					Password: "adminpass",
					Role:     models.RoleAdmin,
				},
				output: models.User{
					ID:       2,
					Username: "admin1",
					Email:    "admin@email.com",
					Password: "adminpass",
					Role:     models.RoleAdmin,
				},
				err: nil,
			},
			expBody: `{"message":"Success","result":{"id":2,"username":"admin1","email":"admin@email.com","role":"admin","created_at":"0001-01-01 00:00:00","updated_at":"0001-01-01 00:00:00"}}`,
			expErr:  nil,
			expCode: http.StatusOK,
		},
		"all fields empty": {
			givenInput:  `{"username": "", "email": "", "password": "", "role": ""}`,
			mockUseCase: mockUseCase{expCall: false},
			expBody:     `{"message":"Internal server error"}`,
			expErr:      nil,
			expCode:     http.StatusInternalServerError,
		},
		"email format edge": {
			givenInput:  `{"username": "test9", "email": "test9@", "password": "pass9", "role": "user"}`,
			mockUseCase: mockUseCase{expCall: false},
			expBody:     `{"message":"Internal server error"}`,
			expErr:      nil,
			expCode:     http.StatusInternalServerError,
		},
		"extra field": {
			givenInput: `{"username": "test10", "email": "test10@email.com", "password": "pass10", "role": "user", "extra": "field"}`,
			mockUseCase: mockUseCase{
				expCall: true,
				input: &models.User{
					Username: "test10",
					Email:    "test10@email.com",
					Password: "pass10",
					Role:     models.RoleUser,
				},
				output: models.User{
					ID:       3,
					Username: "test10",
					Email:    "test10@email.com",
					Password: "pass10",
					Role:     models.RoleUser,
				},
				err: nil,
			},
			expBody: `{"message":"Success","result":{"id":3,"username":"test10","email":"test10@email.com","role":"user","created_at":"0001-01-01 00:00:00","updated_at":"0001-01-01 00:00:00"}}`,
			expErr:  nil,
			expCode: http.StatusOK,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			// Given
			cfg := &config.Config{}
			apiLogger := logger.NewApiLogger(cfg)
			apiLogger.InitLogger()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mock.NewMockUseCase(ctrl)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer([]byte(tc.givenInput)))
			c.Request.Header.Add("Content-Type", "application/json")

			if tc.mockUseCase.expCall {
				mockUseCase.EXPECT().Register(gomock.Any(), gomock.Eq(tc.mockUseCase.input)).Return(&tc.mockUseCase.output, tc.mockUseCase.err)
			}

			// When
			h := authhttp.NewHandlers(cfg, mockUseCase, apiLogger)
			h.Register(c)

			// Then
			if tc.expErr != nil {
				assert.Equal(t, w.Code, tc.expCode)
			} else {
				assert.NoError(t, tc.expErr)
				assert.Equal(t, tc.expCode, w.Code)
				assert.JSONEq(t, tc.expBody, w.Body.String())
			}

		})
	}
}

func TestHandlers_Login(t *testing.T) {
	type mockUseCase struct {
		expLoginCall              bool
		loginInput                *models.User
		loginOutput               models.User
		loginErr                  error
		expGenerateRefreshCall    bool
		generateRefreshUserID     int
		generateRefreshToken      string
		generateRefreshExpiry     time.Time
		generateRefreshErr        error
	}

	tcs := map[string]struct {
		givenInput  string
		mockUseCase mockUseCase
		expBody     string
		expErr      error
		expCode     int
	}{
		"success": {
			givenInput: `{
				"username": "test",
				"password": "pass"
			}`,
			mockUseCase: mockUseCase{
				expLoginCall: true,
				loginInput: &models.User{
					Username: "test",
					Password: "pass",
				},
				loginOutput: models.User{
					ID:       1,
					Username: "test",
					Email:    "test@example.com",
					Password: "hashed_pass",
					Role:     models.RoleUser,
				},
				loginErr:               nil,
				expGenerateRefreshCall: true,
				generateRefreshUserID:  1,
				generateRefreshToken:   "refresh_token_123",
				generateRefreshExpiry:  time.Now().Add(time.Hour * 24 * 7),
				generateRefreshErr:     nil,
			},
			expCode: http.StatusOK,
			expBody: `{"token":"*","expires_at":"*","refresh_token":"refresh_token_123","refresh_token_expires_at":"*","user":{"id":1,"username":"test","email":"test@example.com","role":"user","created_at":"*","updated_at":"*"}}`,
		},
		"empty_username": {
			givenInput: `{
				"username": "",
				"password": "pass"
			}`,
			mockUseCase: mockUseCase{},
			expCode: http.StatusInternalServerError,
		},
		"empty_password": {
			givenInput: `{
				"username": "test",
				"password": ""
			}`,
			mockUseCase: mockUseCase{},
			expCode: http.StatusInternalServerError,
		},
		"invalid_credentials": {
			givenInput: `{
				"username": "wrong",
				"password": "wrong"
			}`,
			mockUseCase: mockUseCase{
				expLoginCall: true,
				loginInput: &models.User{
					Username: "wrong",
					Password: "wrong",
				},
				loginErr: auth.ErrInvalidCredentials,
			},
			expCode: http.StatusUnauthorized,
		},
		"invalid_json": {
			givenInput:  `{"username": "test", "password":}`,
			mockUseCase: mockUseCase{},
			expCode:     http.StatusBadRequest,
		},
		"refresh_token_error": {
			givenInput: `{
				"username": "test",
				"password": "pass"
			}`,
			mockUseCase: mockUseCase{
				expLoginCall: true,
				loginInput: &models.User{
					Username: "test",
					Password: "pass",
				},
				loginOutput: models.User{
					ID:       1,
					Username: "test",
					Email:    "test@example.com",
					Password: "hashed_pass",
					Role:     models.RoleUser,
				},
				loginErr:               nil,
				expGenerateRefreshCall: true,
				generateRefreshUserID:  1,
				generateRefreshErr:     errors.New("refresh token error"),
			},
			expCode: http.StatusInternalServerError,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			// Given
			cfg := &config.Config{}
			apiLogger := logger.NewApiLogger(cfg)
			apiLogger.InitLogger()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mock.NewMockUseCase(ctrl)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer([]byte(tc.givenInput)))
			c.Request.Header.Add("Content-Type", "application/json")

			if tc.mockUseCase.expLoginCall {
				mockUseCase.EXPECT().Login(gomock.Any(), gomock.Eq(tc.mockUseCase.loginInput)).Return(&tc.mockUseCase.loginOutput, tc.mockUseCase.loginErr)
				
				if tc.mockUseCase.loginErr == nil && tc.mockUseCase.expGenerateRefreshCall {
					mockUseCase.EXPECT().GenerateRefreshToken(gomock.Any(), tc.mockUseCase.generateRefreshUserID).Return(
						tc.mockUseCase.generateRefreshToken, 
						tc.mockUseCase.generateRefreshExpiry, 
						tc.mockUseCase.generateRefreshErr,
					)
				}
			}

			// When
			h := authhttp.NewHandlers(cfg, mockUseCase, apiLogger)
			h.Login(c)

			// Then
			assert.Equal(t, tc.expCode, w.Code)
		})
	}
}

func TestHandlers_GetUserByID(t *testing.T) {
	type mockUseCase struct {
		expCall bool
		userID  int
		output  models.User
		err     error
	}

	tcs := map[string]struct {
		userIDParam string
		mockUseCase mockUseCase
		expBody     string
		expCode     int
	}{
		"success": {
			userIDParam: "1",
			mockUseCase: mockUseCase{
				expCall: true,
				userID:  1,
				output: models.User{
					ID:       1,
					Username: "test",
					Email:    "test@example.com",
					Role:     models.RoleUser,
				},
				err: nil,
			},
			expBody: `{"message":"Success","result":{"id":1,"username":"test","email":"test@example.com","role":"user","created_at":"0001-01-01 00:00:00","updated_at":"0001-01-01 00:00:00"}}`,
			expCode: http.StatusOK,
		},
		"admin_user": {
			userIDParam: "2",
			mockUseCase: mockUseCase{
				expCall: true,
				userID:  2,
				output: models.User{
					ID:       2,
					Username: "admin",
					Email:    "admin@example.com",
					Role:     models.RoleAdmin,
				},
				err: nil,
			},
			expBody: `{"message":"Success","result":{"id":2,"username":"admin","email":"admin@example.com","role":"admin","created_at":"0001-01-01 00:00:00","updated_at":"0001-01-01 00:00:00"}}`,
			expCode: http.StatusOK,
		},
		"zero_id": {
			userIDParam: "0",
			mockUseCase: mockUseCase{
				expCall: true,
				userID:  0,
				err:     auth.ErrUserNotFound,
			},
			expCode: http.StatusNotFound,
		},
		"negative_id": {
			userIDParam: "-1",
			mockUseCase: mockUseCase{
				expCall: true,
				userID:  -1,
				err:     auth.ErrUserNotFound,
			},
			expCode: http.StatusNotFound,
		},
		"invalid_id": {
			userIDParam: "abc",
			mockUseCase: mockUseCase{
				expCall: false,
			},
			expCode: http.StatusBadRequest,
		},
		"user_not_found": {
			userIDParam: "999",
			mockUseCase: mockUseCase{
				expCall: true,
				userID:  999,
				err:     auth.ErrUserNotFound,
			},
			expCode: http.StatusNotFound,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			// Given
			cfg := &config.Config{}
			apiLogger := logger.NewApiLogger(cfg)
			apiLogger.InitLogger()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mock.NewMockUseCase(ctrl)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest(http.MethodGet, "/auth/user/"+tc.userIDParam, nil)
			c.Params = []gin.Param{{
				Key:   "userId",
				Value: tc.userIDParam,
			}}

			if tc.mockUseCase.expCall {
				mockUseCase.EXPECT().GetUserByID(gomock.Any(), tc.mockUseCase.userID).Return(&tc.mockUseCase.output, tc.mockUseCase.err)
			}

			// When
			h := authhttp.NewHandlers(cfg, mockUseCase, apiLogger)
			h.GetUserByID(c)

			// Then
			assert.Equal(t, tc.expCode, w.Code)
			if tc.expBody != "" {
				assert.JSONEq(t, tc.expBody, w.Body.String())
			}
		})
	}
}

func TestHandlers_RefreshToken(t *testing.T) {
	type mockUseCase struct {
		expValidateCall       bool
		validateInput         string
		validateOutput        models.RefreshToken
		validateErr           error
		expGetUserCall        bool
		getUserID             int
		getUserOutput         models.User
		getUserErr            error
		expGenerateRefreshCall bool
		generateRefreshUserID  int
		generateRefreshToken   string
		generateRefreshExpiry  time.Time
		generateRefreshErr     error
	}

	tcs := map[string]struct {
		givenInput  string
		mockUseCase mockUseCase
		expBody     string
		expCode     int
	}{
		"success": {
			givenInput: `{
				"refresh_token": "valid_refresh_token"
			}`,
			mockUseCase: mockUseCase{
				expValidateCall: true,
				validateInput:   "valid_refresh_token",
				validateOutput: models.RefreshToken{
					ID:        1,
					UserID:    1,
					Token:     "valid_refresh_token",
					ExpiresAt: time.Now().Add(time.Hour * 24),
					Revoked:   false,
				},
				validateErr:    nil,
				expGetUserCall: true,
				getUserID:      1,
				getUserOutput: models.User{
					ID:       1,
					Username: "test",
					Email:    "test@example.com",
					Role:     models.RoleUser,
				},
				getUserErr:             nil,
				expGenerateRefreshCall: true,
				generateRefreshUserID:  1,
				generateRefreshToken:   "new_refresh_token",
				generateRefreshExpiry:  time.Now().Add(time.Hour * 24 * 7),
				generateRefreshErr:     nil,
			},
			expCode: http.StatusOK,
			expBody: `{"token":"*","expires_at":"*","refresh_token":"new_refresh_token","refresh_expires_at":"*"}`,
		},
		"admin_user": {
			givenInput: `{
				"refresh_token": "admin_refresh_token"
			}`,
			mockUseCase: mockUseCase{
				expValidateCall: true,
				validateInput:   "admin_refresh_token",
				validateOutput: models.RefreshToken{
					ID:        5,
					UserID:    2,
					Token:     "admin_refresh_token",
					ExpiresAt: time.Now().Add(time.Hour * 24),
					Revoked:   false,
				},
				validateErr:    nil,
				expGetUserCall: true,
				getUserID:      2,
				getUserOutput: models.User{
					ID:       2,
					Username: "admin",
					Email:    "admin@example.com",
					Role:     models.RoleAdmin,
				},
				getUserErr:             nil,
				expGenerateRefreshCall: true,
				generateRefreshUserID:  2,
				generateRefreshToken:   "new_admin_refresh_token",
				generateRefreshExpiry:  time.Now().Add(time.Hour * 24 * 7),
				generateRefreshErr:     nil,
			},
			expCode: http.StatusOK,
			expBody: `{"token":"*","expires_at":"*","refresh_token":"new_admin_refresh_token","refresh_expires_at":"*"}`,
		},
		"empty_token": {
			givenInput: `{
				"refresh_token": ""
			}`,
			mockUseCase: mockUseCase{},
			expCode: http.StatusInternalServerError,
		},
		"invalid_token": {
			givenInput: `{
				"refresh_token": "invalid_token"
			}`,
			mockUseCase: mockUseCase{
				expValidateCall: true,
				validateInput:   "invalid_token",
				validateErr:     auth.ErrInvalidToken,
			},
			expCode: http.StatusUnauthorized,
		},
		"revoked_token": {
			givenInput: `{
				"refresh_token": "revoked_token"
			}`,
			mockUseCase: mockUseCase{
				expValidateCall: true,
				validateInput:   "revoked_token",
				validateOutput: models.RefreshToken{
					ID:        2,
					UserID:    1,
					Token:     "revoked_token",
					ExpiresAt: time.Now().Add(time.Hour * 24),
					Revoked:   true,
				},
				validateErr: nil,
			},
			expCode: http.StatusUnauthorized,
		},
		"user_not_found": {
			givenInput: `{
				"refresh_token": "valid_token_user_not_found"
			}`,
			mockUseCase: mockUseCase{
				expValidateCall: true,
				validateInput:   "valid_token_user_not_found",
				validateOutput: models.RefreshToken{
					ID:        3,
					UserID:    999,
					Token:     "valid_token_user_not_found",
					ExpiresAt: time.Now().Add(time.Hour * 24),
					Revoked:   false,
				},
				validateErr:    nil,
				expGetUserCall: true,
				getUserID:      999,
				getUserErr:     auth.ErrUserNotFound,
			},
			expCode: http.StatusNotFound,
		},
		"refresh_token_generation_error": {
			givenInput: `{
				"refresh_token": "valid_token_refresh_error"
			}`,
			mockUseCase: mockUseCase{
				expValidateCall: true,
				validateInput:   "valid_token_refresh_error",
				validateOutput: models.RefreshToken{
					ID:        4,
					UserID:    2,
					Token:     "valid_token_refresh_error",
					ExpiresAt: time.Now().Add(time.Hour * 24),
					Revoked:   false,
				},
				validateErr:    nil,
				expGetUserCall: true,
				getUserID:      2,
				getUserOutput: models.User{
					ID:       2,
					Username: "test2",
					Email:    "test2@example.com",
					Role:     models.RoleUser,
				},
				getUserErr:             nil,
				expGenerateRefreshCall: true,
				generateRefreshUserID:  2,
				generateRefreshErr:     errors.New("refresh token generation error"),
			},
			expCode: http.StatusInternalServerError,
		},
		"invalid_json": {
			givenInput:  `{"refresh_token":}`,
			mockUseCase: mockUseCase{},
			expCode:     http.StatusBadRequest,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			// Given
			cfg := &config.Config{}
			apiLogger := logger.NewApiLogger(cfg)
			apiLogger.InitLogger()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mock.NewMockUseCase(ctrl)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer([]byte(tc.givenInput)))
			c.Request.Header.Add("Content-Type", "application/json")

			if tc.mockUseCase.expValidateCall {
				mockUseCase.EXPECT().ValidateRefreshToken(gomock.Any(), tc.mockUseCase.validateInput).Return(&tc.mockUseCase.validateOutput, tc.mockUseCase.validateErr)

				if tc.mockUseCase.validateErr == nil && !tc.mockUseCase.validateOutput.Revoked && tc.mockUseCase.expGetUserCall {
					mockUseCase.EXPECT().GetUserByID(gomock.Any(), tc.mockUseCase.getUserID).Return(&tc.mockUseCase.getUserOutput, tc.mockUseCase.getUserErr)

					if tc.mockUseCase.getUserErr == nil && tc.mockUseCase.expGenerateRefreshCall {
						mockUseCase.EXPECT().GenerateRefreshToken(gomock.Any(), tc.mockUseCase.generateRefreshUserID).Return(
							tc.mockUseCase.generateRefreshToken,
							tc.mockUseCase.generateRefreshExpiry,
							tc.mockUseCase.generateRefreshErr,
						)
					}
				}
			}

			// When
			h := authhttp.NewHandlers(cfg, mockUseCase, apiLogger)
			h.RefreshToken(c)

			// Then
			assert.Equal(t, tc.expCode, w.Code)
		})
	}
}
