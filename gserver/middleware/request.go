package middleware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/go-playground/validator.v9"

	"github.com/tomoyane/grant-n-z/gserver/cache"
	"github.com/tomoyane/grant-n-z/gserver/common/ctx"
	"github.com/tomoyane/grant-n-z/gserver/log"
	"github.com/tomoyane/grant-n-z/gserver/model"
	"github.com/tomoyane/grant-n-z/gserver/service"
)

const (
	Authorization = "Authorization"
	Key           = "Api-Key"
	ContentType   = "Content-Type"
)

var iInstance Interceptor

type Interceptor interface {
	// Http request interceptor
	Intercept(next http.HandlerFunc) http.HandlerFunc

	// Interceptor with user authentication
	InterceptAuthenticateUser(next http.HandlerFunc) http.HandlerFunc

	// Interceptor with operator authentication
	InterceptAuthenticateOperator(next http.HandlerFunc) http.HandlerFunc
}

type InterceptorImpl struct {
	tokenService service.TokenService
	userService  service.UserService
	redisClient  cache.RedisClient
}

func GetInterceptorInstance() Interceptor {
	if iInstance == nil {
		iInstance = NewInterceptor()
	}
	return iInstance
}

func NewInterceptor() Interceptor {
	log.Logger.Info("New `Interceptor` instance")
	log.Logger.Info("Inject `TokenService`, `UserService` to `Interceptor`")
	return InterceptorImpl{
		tokenService: service.GetTokenServiceInstance(),
		userService:  service.GetUserServiceInstance(),
		redisClient:  cache.GetRedisClientInstance(),
	}
}

func (i InterceptorImpl) Intercept(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				err := model.InternalServerError("Failed to request body bind")
				model.WriteError(w, err.ToJson(), err.Code)
			}
		}()

		if err := intercept(w, r); err != nil {
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (i InterceptorImpl) InterceptAuthenticateUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				err := model.InternalServerError("Failed to request body bind")
				model.WriteError(w, err.ToJson(), err.Code)
			}
		}()

		if err := intercept(w, r); err != nil {
			return
		}

		token := r.Header.Get(Authorization)
		authUser, err := i.tokenService.VerifyUserToken(token)
		if err != nil {
			model.WriteError(w, err.ToJson(), err.Code)
			return
		}

		ctx.SetUserId(authUser.UserId)
		ctx.SetUserUuid(authUser.UserUuid.String())
		ctx.SetServiceId(authUser.ServiceId)

		next.ServeHTTP(w, r)
	}
}

func (i InterceptorImpl) InterceptAuthenticateOperator(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				err := model.InternalServerError("Failed to request body bind")
				model.WriteError(w, err.ToJson(), err.Code)
			}
		}()

		if err := intercept(w, r); err != nil {
			return
		}

		token := r.Header.Get(Authorization)
		authUser, err := i.tokenService.VerifyOperatorToken(token)
		if err != nil {
			model.WriteError(w, err.ToJson(), err.Code)
			return
		}

		ctx.SetUserId(authUser.UserId)
		ctx.SetUserUuid(authUser.UserUuid.String())
		ctx.SetServiceId(authUser.ServiceId)

		next.ServeHTTP(w, r)
	}
}

// Intercept http request
// Set api key
func intercept(w http.ResponseWriter, r *http.Request) *model.ErrorResBody {
	if err := validateHeader(r); err != nil {
		model.WriteError(w, err.ToJson(), err.Code)
		return err
	}

	ctx.SetApiKey(r.Header.Get(Key))
	w.Header().Set(ContentType, "application/json")
	return nil
}

// Validate http request header
func validateHeader(r *http.Request) *model.ErrorResBody {
	if r.Method != http.MethodGet {
		if r.Header.Get(ContentType) != "application/json" {
			log.Logger.Info("Not allowed content-type")
			return model.BadRequest("Need to content type is only json.")
		}
	}
	return nil
}

// Bind request body what http request converts to interface
func BindBody(w http.ResponseWriter, r *http.Request, i interface{}) *model.ErrorResBody {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Logger.Info("Cannot bind request body.", err.Error())
		err := model.InternalServerError("Failed to request body bind")
		model.WriteError(w, err.ToJson(), err.Code)
		return err
	}

	if len(body) == 0 {
		err := model.BadRequest("Request is empty.")
		model.WriteError(w, err.ToJson(), err.Code)
		return err
	}

	if err := json.Unmarshal(body, i); err != nil {
		err := model.BadRequest("Request is not json.")
		model.WriteError(w, err.ToJson(), err.Code)
		return err
	}

	return nil
}

// Validate request body
func ValidateBody(w http.ResponseWriter, i interface{}) *model.ErrorResBody {
	if err := validator.New().Struct(i); err != nil {
		log.Logger.Info(err.Error())
		err := model.BadRequest("Invalid request.")
		model.WriteError(w, err.ToJson(), err.Code)
		return err
	}
	return nil
}