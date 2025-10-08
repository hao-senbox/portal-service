package user

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"portal/pkg/constants"
	"portal/pkg/consul"
	"time"

	"github.com/hashicorp/consul/api"
)

type UserService interface {
	GetUserInfor(ctx context.Context, userID string) (*UserInfor, error)
	GetStudentInfor(ctx context.Context, studentID string) (*UserInfor, error)
	GetTeacherInfor(ctx context.Context, studentID string) (*UserInfor, error)
	GetStaffInfor(ctx context.Context, studentID string) (*UserInfor, error)
}

type userService struct {
	client *callAPI
}

type callAPI struct {
	client       consul.ServiceDiscovery
	clientServer *api.CatalogService
}

var (
	mainService = "go-main-service"
)

// ===== Helpers =====

func logErr(prefix string, err error) {
	if err != nil {
		log.Printf("[userService] %s: %v", prefix, err)
	}
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func castToInt64(v interface{}) int64 {
	switch val := v.(type) {
	case float64:
		return int64(val)
	case int:
		return int64(val)
	default:
		return 0
	}
}

func castToBool(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val == "true" || val == "1"
	default:
		return false
	}
}

// ===== Constructors =====

func NewUserService(client *api.Client) UserService {
	mainServiceAPI := NewServiceAPI(client, mainService)
	return &userService{
		client: mainServiceAPI,
	}
}

func NewServiceAPI(client *api.Client, serviceName string) *callAPI {
	sd, err := consul.NewServiceDiscovery(client, serviceName)
	if err != nil {
		logErr("Error creating service discovery", err)
		return &callAPI{}
	}

	var service *api.CatalogService

	for i := 0; i < 10; i++ {
		service, err = sd.DiscoverService()
		if err == nil && service != nil {
			break
		}
		fmt.Printf("Waiting for service %s... retry %d/10\n", serviceName, i+1)
		time.Sleep(3 * time.Second)
	}

	if service == nil {
		fmt.Printf("Service %s not found after retries, continuing anyway...\n", serviceName)
	}

	if os.Getenv("LOCAL_TEST") == "true" {
		log.Printf("[userService] LOCAL_TEST= true — override address => localhost")
		service.ServiceAddress = "localhost"
	}

	return &callAPI{
		client:       sd,
		clientServer: service,
	}
}

// ===== Public methods (fail-safe: log và trả về nil, nil) =====

func (u *userService) GetUserInfor(ctx context.Context, userID string) (*UserInfor, error) {
	if u.client == nil || u.client.clientServer == nil || u.client.client == nil {
		log.Printf("[userService] client not ready (service discovery/server nil)")
		return nil, nil
	}

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok || token == "" {
		log.Printf("[userService] token not found in context")
		return nil, nil
	}

	data, err := u.client.getUserInfor(userID, token)
	if err != nil {
		logErr("getUserInfor call error", err)
		return nil, nil
	}
	if data == nil {
		log.Printf("[userService] getUserInfor: empty data for userID=%s", userID)
		return nil, nil
	}

	innerData, ok := data["data"].(map[string]interface{})
	if !ok {
		log.Printf("[userService] invalid response: missing 'data' field")
		return nil, nil
	}

	var avatar Avatar
	if raw, exists := innerData["avatar"].(map[string]interface{}); exists && raw != nil {
		avatar = Avatar{
			ImageID:  uint64(castToInt64(raw["image_id"])),
			ImageKey: getString(raw, "image_key"),
			ImageUrl: getString(raw, "image_url"),
			Index:    int(castToInt64(raw["index"])),
			IsMain:   castToBool(raw["is_main"]),
		}
	}

	user := &UserInfor{
		UserID:   getString(innerData, "id"),
		UserName: getString(innerData, "name"),
		Avartar:  avatar,
	}

	// Nếu không có id/name, có thể coi là không hợp lệ -> trả nil
	if user.UserID == "" && user.UserName == "" {
		log.Printf("[userService] user data missing id/name")
		return nil, nil
	}

	return user, nil
}

func (u *userService) GetStudentInfor(ctx context.Context, studentID string) (*UserInfor, error) {
	if u.client == nil || u.client.clientServer == nil || u.client.client == nil {
		log.Printf("[userService] client not ready (service discovery/server nil)")
		return nil, nil
	}

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok || token == "" {
		log.Printf("[userService] token not found in context")
		return nil, nil
	}

	data, err := u.client.getStudentInfor(studentID, token)
	if err != nil {
		logErr("getStudentInfor call error", err)
		return nil, nil
	}
	if data == nil {
		log.Printf("[userService] getStudentInfor: empty data for studentID=%s", studentID)
		return nil, nil
	}

	innerData, ok := data["data"].(map[string]interface{})
	if !ok {
		log.Printf("[userService] invalid response: missing 'data' field")
		return nil, nil
	}

	var avatar Avatar
	if raw, exists := innerData["avatar"].(map[string]interface{}); exists && raw != nil {
		avatar = Avatar{
			ImageID:  uint64(castToInt64(raw["image_id"])),
			ImageKey: getString(raw, "image_key"),
			ImageUrl: getString(raw, "image_url"),
			Index:    int(castToInt64(raw["index"])),
			IsMain:   castToBool(raw["is_main"]),
		}
	}

	user := &UserInfor{
		UserID:   getString(innerData, "id"),
		UserName: getString(innerData, "name"),
		Avartar:  avatar,
	}
	if user.UserID == "" && user.UserName == "" {
		log.Printf("[userService] student data missing id/name")
		return nil, nil
	}

	return user, nil
}

func (u *userService) GetTeacherInfor(ctx context.Context, teacherID string) (*UserInfor, error) {
	if u.client == nil || u.client.clientServer == nil || u.client.client == nil {
		log.Printf("[userService] client not ready (service discovery/server nil)")
		return nil, nil
	}

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok || token == "" {
		log.Printf("[userService] token not found in context")
		return nil, nil
	}

	data, err := u.client.getTeacherInfor(teacherID, token)
	if err != nil {
		logErr("getTeacherInfor call error", err)
		return nil, nil
	}
	if data == nil {
		log.Printf("[userService] getTeacherInfor: empty data for teacherID=%s", teacherID)
		return nil, nil
	}

	innerData, ok := data["data"].(map[string]interface{})
	if !ok {
		log.Printf("[userService] invalid response: missing 'data' field")
		return nil, nil
	}

	var avatar Avatar
	if raw, exists := innerData["avatar"].(map[string]interface{}); exists && raw != nil {
		avatar = Avatar{
			ImageID:  uint64(castToInt64(raw["image_id"])),
			ImageKey: getString(raw, "image_key"),
			ImageUrl: getString(raw, "image_url"),
			Index:    int(castToInt64(raw["index"])),
			IsMain:   castToBool(raw["is_main"]),
		}
	}

	user := &UserInfor{
		UserID:   getString(innerData, "id"),
		UserName: getString(innerData, "name"),
		Avartar:  avatar,
	}
	if user.UserID == "" && user.UserName == "" {
		log.Printf("[userService] teacher data missing id/name")
		return nil, nil
	}

	return user, nil
}

func (u *userService) GetStaffInfor(ctx context.Context, staffID string) (*UserInfor, error) {
	if u.client == nil || u.client.clientServer == nil || u.client.client == nil {
		log.Printf("[userService] client not ready (service discovery/server nil)")
		return nil, nil
	}

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok || token == "" {
		log.Printf("[userService] token not found in context")
		return nil, nil
	}

	data, err := u.client.getStaffInfor(staffID, token)
	if err != nil {
		logErr("getStaffInfor call error", err)
		return nil, nil
	}
	if data == nil {
		log.Printf("[userService] getStaffInfor: empty data for staffID=%s", staffID)
		return nil, nil
	}

	innerData, ok := data["data"].(map[string]interface{})
	if !ok {
		log.Printf("[userService] invalid response: missing 'data' field")
		return nil, nil
	}

	var avatar Avatar
	if raw, exists := innerData["avatar"].(map[string]interface{}); exists && raw != nil {
		avatar = Avatar{
			ImageID:  uint64(castToInt64(raw["image_id"])),
			ImageKey: getString(raw, "image_key"),
			ImageUrl: getString(raw, "image_url"),
			Index:    int(castToInt64(raw["index"])),
			IsMain:   castToBool(raw["is_main"]),
		}
	}

	user := &UserInfor{
		UserID:   getString(innerData, "id"),
		UserName: getString(innerData, "name"),
		Avartar:  avatar,
	}
	if user.UserID == "" && user.UserName == "" {
		log.Printf("[userService] staff data missing id/name")
		return nil, nil
	}

	return user, nil
}

func (c *callAPI) getUserInfor(userID string, token string) (map[string]interface{}, error) {
	return c.getJSON(fmt.Sprintf("/v1/gateway/users/%s", userID), token)
}

func (c *callAPI) getStudentInfor(studentID string, token string) (map[string]interface{}, error) {
	return c.getJSON(fmt.Sprintf("/v1/gateway/students/%s", studentID), token)
}

func (c *callAPI) getTeacherInfor(teacherID string, token string) (map[string]interface{}, error) {
	return c.getJSON(fmt.Sprintf("/v1/gateway/teachers/%s", teacherID), token)
}

func (c *callAPI) getStaffInfor(staffID string, token string) (map[string]interface{}, error) {
	return c.getJSON(fmt.Sprintf("/v1/gateway/staffs/%s", staffID), token)
}

func (c *callAPI) getJSON(endpoint, token string) (map[string]interface{}, error) {

	if c == nil || c.client == nil || c.clientServer == nil {
		log.Printf("[userService] service discovery/client not ready for endpoint %s", endpoint)
		return nil, nil
	}

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		logErr("Error calling API "+endpoint, err)
		return nil, nil
	}
	if res == "" {
		log.Printf("[userService] empty response from %s", endpoint)
		return nil, nil
	}

	var myMap map[string]interface{}
	if err := json.Unmarshal([]byte(res), &myMap); err != nil {
		logErr("Error unmarshalling response "+endpoint, err)
		return nil, nil
	}
	return myMap, nil
}
