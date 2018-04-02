package util

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/spf13/viper"

	"../exception"
)

type (
	Utilities struct {
		v *viper.Viper
	}
)

func NewUtilities(v *viper.Viper) *Utilities {
	// will have another function to set up Redis pool
	return &Utilities{v}
}

/*
	Utility methods
*/

func (u Utilities) GetIntConfigValue(key string) int {
	return u.v.GetInt(key)
}

func (u Utilities) GetStringConfigValue(key string) string {
	return u.v.GetString(key)
}

func (u Utilities) GetBooleanConfigValue(key string) bool {
	return u.v.GetBool(key)
}

func (u Utilities) GetMapArrConfigValue(key string) (
	[]map[string]interface{}, *exception.ASError) {

	return toArrayMap(u.v.Get(key), &u)
}

func (u Utilities) GetActiveEnv() int {
	return u.GetIntConfigValue("environment.active")
}

func (u Utilities) GetActiveEnvHost() (string, *exception.ASError) {

	prefix, asErr := u.getActiveEnvPrefix()
	if asErr != nil {
		return "", nil
	}
	hostKey := fmt.Sprintf("%s.host", prefix)
	return u.GetStringConfigValue(hostKey), nil
}

func (u Utilities) GetMessage(key string) string {
	msgKey := fmt.Sprintf("messages.%s", key)
	return u.GetStringConfigValue(msgKey)
}

func (u Utilities) GetErrorMessage(key string) string {
	msgKey := fmt.Sprintf("messages.error.%s", key)
	return u.GetStringConfigValue(msgKey)
}

func (u Utilities) GetWarningMessage(key string) string {
	msgKey := fmt.Sprintf("messages.warning.%s", key)
	return u.GetStringConfigValue(msgKey)
}

func (u Utilities) GetDBUrl() (string, *exception.ASError) {

	prefix, asErr := u.getActiveEnvPrefix()
	if asErr != nil {
		return "", nil
	}
	hostKey := fmt.Sprintf("%s.db.host", prefix)
	host := u.GetStringConfigValue(hostKey)
	port := u.GetIntConfigValue("db.port")
	url := fmt.Sprintf("%s:%d", host, port)

	return url, nil
}

func (u Utilities) GetDBDialInfo() (*mgo.DialInfo, *exception.ASError) {

	url, asErr := u.GetDBUrl()
	if asErr != nil {
		return nil, asErr
	}
	timeout := u.GetIntConfigValue("db.timeout")
	dbName := u.GetDBName()
	poolLimit := u.GetIntConfigValue("db.pool_limit")

	return &mgo.DialInfo{
		Addrs:     []string{url},
		Timeout:   time.Duration(timeout) * time.Second,
		Database:  dbName,
		PoolLimit: poolLimit,
	}, nil
}

func (u Utilities) GetDBName() string {

	return u.GetStringConfigValue("db.db_name")
}

func (u Utilities) GetDBDocName() string {

	return u.GetStringConfigValue("db.doc_name")
}

func (u Utilities) GetError(
	uuid exception.UUID, errKey string, err error) *exception.ASError {

	var errMsg string
	if errKey == "" {
		errMsg = err.Error()
	} else {
		errMsg = u.GetErrorMessage(errKey)
	}

	asErr := exception.NewASError(uuid, errMsg, err)
	log.Println(asErr.ErrorMessage())

	return asErr
}

func (u Utilities) GetWarning(
	uuid exception.UUID, warnKey string) *exception.ASWarning {

	warnMsg := u.GetWarningMessage(warnKey)
	asWarning := exception.NewASWarning(uuid, warnMsg)
	log.Println(asWarning.WarningMessage())

	return asWarning
}

func (u Utilities) SerializeObject(i interface{}) ([]byte, *exception.ASError) {

	serrialized, err := json.Marshal(i)
	if err != nil {
		asErr := u.GetError(
			exception.AS00001, "serialize_object_error", err)
		return nil, asErr
	}
	return serrialized, nil
}

func (u Utilities) UnserializeObject(
	data []byte, i interface{}) *exception.ASError {

	err := json.Unmarshal(data, i)
	if err != nil {
		asErr := u.GetError(
			exception.AS00002, "unserialize_object_error", err)
		return asErr
	}

	return nil
}

func (u Utilities) GetLogFilePath(typ string) string {

	key := fmt.Sprintf("log_file_path.%s", typ)
	return u.GetStringConfigValue(key)
}

/*
	Private methods
*/

func (u Utilities) getActiveEnvPrefix() (string, *exception.ASError) {
	env := u.GetActiveEnv()
	envMap, asErr := u.GetMapArrConfigValue("env-map")
	if asErr != nil {
		return "", asErr
	}

	switch env {
	case toIntFromInt64Inteface(envMap[0]["index"]):
		return envMap[0]["type"].(string), nil
	case toIntFromInt64Inteface(envMap[1]["index"]):
		return envMap[1]["type"].(string), nil
	default:
		asErr1 := u.GetError(exception.AS00006, "invalid_env_type", nil)
		return "", asErr1
	}
}

func toArray(i interface{}, u *Utilities) (
	[]interface{}, *exception.ASError) {

	arr, ok := i.([]interface{})
	if !ok {
		asErr := u.GetError(exception.AS00004, "to_array_error", nil)
		return nil, asErr
	}
	return arr, nil
}

func toMap(i interface{}, u *Utilities) (
	map[string]interface{}, *exception.ASError) {

	m, ok := i.(map[string]interface{})
	if !ok {
		asErr := u.GetError(exception.AS00005, "to_map_error", nil)
		return nil, asErr
	}
	return m, nil
}

func toArrayMap(i interface{}, u *Utilities) ([]map[string]interface{}, *exception.ASError) {
	arr, asErr := toArray(i, u)
	if asErr != nil {
		return nil, asErr
	}

	mapArr := make([]map[string]interface{}, len(arr))

	for i := 0; i < len(arr); i++ {
		mapArr[i], asErr = toMap(arr[i], u)
		if asErr != nil {
			return nil, asErr
		}
	}

	return mapArr, nil
}

func toIntFromInt64Inteface(i interface{}) int {

	intString := fmt.Sprintf("%d", i)
	integer, _ := strconv.ParseInt(intString, 10, 0)

	return int(integer)
}
