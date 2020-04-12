package middleware

import (
	"bytes"
	"github.com/tomoyane/grant-n-z/gnz/entity"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"go.etcd.io/etcd/clientv3"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/tomoyane/grant-n-z/gnz/cache"
	"github.com/tomoyane/grant-n-z/gnz/common"
	"github.com/tomoyane/grant-n-z/gnz/ctx"
	"github.com/tomoyane/grant-n-z/gnz/log"
	"github.com/tomoyane/grant-n-z/gnzserver/service"
)

var interceptor Interceptor

func init() {
	log.InitLogger("info")
	ctx.InitContext()
	ctx.SetUserId(1)
	ctx.SetServiceId(1)
	ctx.SetUserUuid(uuid.New())
	ctx.SetApiKey("test")

	stubConnection, _ := gorm.Open("sqlite3", "/tmp/test_grant_nz.db")
	stubEtcdConnection, _ := clientv3.New(clientv3.Config{
		Endpoints:            []string{},
		DialTimeout:          5 * time.Millisecond,
		DialKeepAliveTimeout: 5 * time.Millisecond,
	})

	userService := service.UserServiceImpl{
		UserRepository: StubUserRepositoryImpl{Connection: stubConnection},
		EtcdClient: cache.EtcdClientImpl{
			Connection: stubEtcdConnection,
			Ctx:        ctx.GetCtx(),
		},
	}

	operatorPolicyService := service.OperatorPolicyServiceImpl{
		OperatorPolicyRepository: StubOperatorPolicyRepositoryImpl{Connection: stubConnection},
		UserRepository:           StubUserRepositoryImpl{Connection: stubConnection},
		RoleRepository:           StubRoleRepositoryImpl{Connection: stubConnection},
	}

	ser := service.ServiceImpl{
		EtcdClient: cache.EtcdClientImpl{
			Connection: stubEtcdConnection,
			Ctx:        ctx.GetCtx(),
		},
		ServiceRepository:    StubServiceRepositoryImpl{Connection: stubConnection},
		RoleRepository:       StubRoleRepositoryImpl{Connection: stubConnection},
		PermissionRepository: StubPermissionRepositoryImpl{Connection: stubConnection},
	}

	policyService := service.PolicyServiceImpl{
		EtcdClient: cache.EtcdClientImpl{
			Connection: stubEtcdConnection,
			Ctx:        ctx.GetCtx(),
		},
		PolicyRepository:     StubPolicyRepositoryImpl{Connection: stubConnection},
		PermissionRepository: StubPermissionRepositoryImpl{Connection: stubConnection},
		RoleRepository:       StubRoleRepositoryImpl{Connection: stubConnection},
		ServiceRepository:    StubServiceRepositoryImpl{Connection: stubConnection},
		GroupRepository:      StubGroupRepositoryImpl{Connection: stubConnection},
	}

	roleService := service.RoleServiceImpl{
		EtcdClient: cache.EtcdClientImpl{
			Connection: stubEtcdConnection,
			Ctx:        ctx.GetCtx(),
		},
		RoleRepository: StubRoleRepositoryImpl{Connection: stubConnection},
	}

	permissionService := service.PermissionServiceImpl{
		EtcdClient: cache.EtcdClientImpl{
			Connection: stubEtcdConnection,
			Ctx:        ctx.GetCtx(),
		},
		PermissionRepository: StubPermissionRepositoryImpl{Connection: stubConnection},
	}

	serviceConfig := common.ServerConfig{
		SignedInPrivateKeyBase64: "test key",
	}

	tokenProcessor := TokenProcessorImpl{
		UserService:           userService,
		OperatorPolicyService: operatorPolicyService,
		Service:               ser,
		PolicyService:         policyService,
		RoleService:           roleService,
		PermissionService:     permissionService,
		ServerConfig:          serviceConfig,
	}

	interceptor = InterceptorImpl{tokenProcessor: tokenProcessor}
}

// Test constructor
func TestGetInterceptorInstance(t *testing.T) {
	GetInterceptorInstance()
}

// Test intercept header
func TestInterceptHeader_Error(t *testing.T) {
	writer := StubResponseWriter{}
	request := http.Request{Header: http.Header{}}
	request.Header.Set("Content-Type", "html/text")
	err := interceptHeader(writer, &request)
	if err == nil {
		t.Errorf("Incorrect TestInterceptHeader_Error test.")
		t.FailNow()
	}
}

// Test intercept header
func TestInterceptHeader_Success(t *testing.T) {
	writer := StubResponseWriter{}
	request := http.Request{Header: http.Header{}}
	request.Header.Set("Content-Type", "application/json")
	err := interceptHeader(writer, &request)
	if err != nil {
		t.Errorf("Incorrect TestInterceptHeader_Success test.")
		t.FailNow()
	}
}

// Test intercept api Key in header
func TestInterceptApiKey_Error(t *testing.T) {
	writer := StubResponseWriter{}
	request := http.Request{Header: http.Header{}}
	request.Header.Set("Api-Key", "")
	err := interceptApiKey(writer, &request)
	if err == nil {
		t.Errorf("Incorrect TestInterceptApiKey_Error test.")
		t.FailNow()
	}
}

// Test intercept api Key in header
func TestInterceptApiKey_Success(t *testing.T) {
	writer := StubResponseWriter{}
	request := http.Request{Header: http.Header{}}
	request.Header.Set("Api-Key", "test_key")
	err := interceptApiKey(writer, &request)
	if err != nil {
		t.Errorf("Incorrect TestInterceptApiKey_Success test.")
		t.FailNow()
	}
}

// Test validate header
func TestValidateHeader_Error(t *testing.T) {
	request := http.Request{Header: http.Header{}}
	request.Header.Set("Content-Type", "html/text")
	err := validateHeader(&request)
	if err == nil {
		t.Errorf("Incorrect TestValidateHeader test.")
		t.FailNow()
	}
}

// Test validate header
func TestValidateHeader_Success(t *testing.T) {
	request := http.Request{Header: http.Header{}}
	request.Header.Set("Content-Type", "application/json")
	err := validateHeader(&request)
	if err != nil {
		t.Errorf("Incorrect TestValidateHeader test.")
		t.FailNow()
	}
}

// Test bind request body
func TestBindBody_Error(t *testing.T) {
	writer := StubResponseWriter{}
	body := ioutil.NopCloser(bytes.NewReader([]byte("")))
	request := http.Request{Header: http.Header{}, Body: body}
	err := BindBody(writer, &request, nil)
	if err == nil {
		t.Errorf("Incorrect TestBindBody_Error test.")
		t.FailNow()
	}

	body = ioutil.NopCloser(bytes.NewReader([]byte("{\"username\":\"test\", \"password\":\"testtest\"}")))
	request = http.Request{Header: http.Header{}, Body: body}
	err = BindBody(writer, &request, nil)
	if err == nil {
		t.Errorf("Incorrect TestBindBody_Error test.")
		t.FailNow()
	}
}

// Test bind request body
func TestBindBody_Success(t *testing.T) {
	writer := StubResponseWriter{}
	body := ioutil.NopCloser(bytes.NewReader([]byte("{\"username\":\"test\", \"email\":\"test@gmail.com\", \"password\":\"testtest\"}")))
	request := http.Request{Header: http.Header{}, Body: body}

	var userEntity *entity.User
	err := BindBody(writer, &request, &userEntity)
	if err != nil {
		t.Errorf("Incorrect TestBindBody_Success test.")
		t.FailNow()
	}
}

// Test bind request body
func TestValidateBody_Error(t *testing.T) {
	writer := StubResponseWriter{}
	user := entity.User{
		Username: "test",
		Email:    "",
		Password: "testtest",
	}
	err := ValidateBody(writer, user)
	if err == nil {
		t.Errorf("Incorrect TestValidateBody_Error test.")
		t.FailNow()
	}
}

// Test bind request body
func TestValidateBody_Success(t *testing.T) {
	writer := StubResponseWriter{}
	user := entity.User{
		Username: "test",
		Email:    "test@gmail.com",
		Password: "testtest",
	}
	err := ValidateBody(writer, user)
	if err != nil {
		t.Errorf("Incorrect TestValidateBody_Success test.")
		t.FailNow()
	}
}

// Test param group id
func TestParamGroupId_Error(t *testing.T) {
	request := http.Request{Header: http.Header{}, URL: &url.URL{}}
	request.URL.Host = "localhost:8080"
	request.URL.Path = "/api/v1/groups/1/user"
	_, err := ParamGroupId(&request)
	if err == nil {
		t.Errorf("Incorrect TestParamGroupId_Error test.")
		t.FailNow()
	}
}

type StubResponseWriter struct {
}

func (w StubResponseWriter) Header() http.Header {
	return http.Header{}
}

func (w StubResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (w StubResponseWriter) WriteHeader(statusCode int) {
}