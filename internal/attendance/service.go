package attendance

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"portal/pkg/constants"
	"portal/pkg/consul"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
)

type AttendanceService interface {
	GetAttendanceInfor(ctx context.Context, userID string) ([]*AttendanceUserInfo, error)
	// GetStudentInfor(ctx context.Context, studentID string) (*AttendanceInfor, error)
	// GetTeacherInfor(ctx context.Context, studentID string) (*AttendanceInfor, error)
	// GetStaffInfor(ctx context.Context, studentID string) (*AttendanceInfor, error)
}

type attendanceService struct {
	client *callAPI
}

type callAPI struct {
	client       consul.ServiceDiscovery
	clientServer *api.CatalogService
}

var (
	mainService = "holiday-service"
)

// ===== Helpers =====

func logErr(prefix string, err error) {
	if err != nil {
		log.Printf("[attendanceService] %s: %v", prefix, err)
	}
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func castToFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
		return 0.0
	default:
		return 0.0
	}
}

// ===== Constructors =====

func NewAttendanceService(client *api.Client) AttendanceService {
	mainServiceAPI := NewServiceAPI(client, mainService)
	return &attendanceService{
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
		log.Printf("[attendanceService] LOCAL_TEST= true — override address => localhost")
		service.ServiceAddress = "localhost"
	}

	return &callAPI{
		client:       sd,
		clientServer: service,
	}
}

// ===== Public methods (fail-safe: log và trả về nil, nil) =====

func (u *attendanceService) GetAttendanceInfor(ctx context.Context, userID string) ([]*AttendanceUserInfo, error) {
	if u.client == nil || u.client.clientServer == nil || u.client.client == nil {
		log.Printf("[attendanceService] client not ready (service discovery/server nil)")
		return nil, nil
	}

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok || token == "" {
		log.Printf("[attendanceService] token not found in context")
		return nil, nil
	}

	data, err := u.client.getAttendanceInfor(userID, token)
	if err != nil {
		logErr("getAttendanceInfor call error", err)
		return nil, nil
	}
	if data == nil {
		log.Printf("[attendanceService] getAttendanceInfor: empty data for userID=%s", userID)
		return nil, nil
	}

	// Handle array of attendance records
	var attendanceRecords []interface{}
	if records, ok := data["data"].([]interface{}); ok {
		attendanceRecords = records
	} else {
		log.Printf("[attendanceService] invalid response: 'data' field is not an array")
		return nil, nil
	}

	if len(attendanceRecords) == 0 {
		log.Printf("[attendanceService] no attendance records found for userID=%s", userID)
		return nil, nil
	}

	var attendanceInfos []*AttendanceUserInfo
	for _, record := range attendanceRecords {
		recordMap, ok := record.(map[string]interface{})
		if !ok {
			log.Printf("[attendanceService] invalid record format, skipping")
			continue
		}

		attendanceInfo := &AttendanceUserInfo{
			AttendanceID: getString(recordMap, "id"),
			StudentID:    getString(recordMap, "user_id"),
			Date:         getString(recordMap, "date"),
			Temperature:  castToFloat64(recordMap["temperature"]),
		}

		// Skip invalid records (those without ID)
		if attendanceInfo.AttendanceID == "" {
			log.Printf("[attendanceService] skipping record with missing id")
			continue
		}

		attendanceInfos = append(attendanceInfos, attendanceInfo)
	}

	if len(attendanceInfos) == 0 {
		log.Printf("[attendanceService] no valid attendance records found for userID=%s", userID)
		return nil, nil
	}

	return attendanceInfos, nil
}

func (c *callAPI) getAttendanceInfor(userID string, token string) (map[string]interface{}, error) {
	return c.getJSON(fmt.Sprintf("/api/v1/gateway/student-temperature?student-id=%s", url.QueryEscape(userID)), token)
}

func (c *callAPI) getJSON(endpoint, token string) (map[string]interface{}, error) {

	if c == nil || c.client == nil || c.clientServer == nil {
		log.Printf("[attendanceService] service discovery/client not ready for endpoint %s", endpoint)
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
		log.Printf("[attendanceService] empty response from %s", endpoint)
		return nil, nil
	}

	var myMap map[string]interface{}
	if err := json.Unmarshal([]byte(res), &myMap); err != nil {
		logErr("Error unmarshalling response "+endpoint, err)
		return nil, nil
	}
	return myMap, nil
}
