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

func SerializeObject(i interface{}) []byte {

	serrialized, err := json.Marshal(i)
	if err != nil {
		return nil // TODO: add error handling later
	}
	return serrialized
}

func UnserializeObject(data []byte, i interface{}) {

	err := json.Unmarshal(data, i)
	if err != nil {
		// TODO: add error handling later
	}
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

func (u Utilities) GetMapArrConfigValue(key string) []map[string]interface{} {

	return toArrayMap(u.v.Get(key))
}

func (u Utilities) GetActiveEnv() int {
	return u.GetIntConfigValue("environment.active")
}

func (u Utilities) GetActiveEnvHost() string {
	prefix := u.getActiveEnvPrefix()
	hostKey := fmt.Sprintf("%s.host", prefix)
	return u.GetStringConfigValue(hostKey)
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

func (u Utilities) GetDBUrl() string {

	prefix := u.getActiveEnvPrefix()
	hostKey := fmt.Sprintf("%s.db.host", prefix)
	host := u.GetStringConfigValue(hostKey)
	port := u.GetIntConfigValue("db.port")
	url := fmt.Sprintf("%s:%d", host, port)

	return url
}

func (u Utilities) GetDBDialInfo() *mgo.DialInfo {

	url := u.GetDBUrl()
	timeout := u.GetIntConfigValue("db.timeout")
	dbName := u.GetDBName()
	poolLimit := u.GetIntConfigValue("db.pool_limit")

	return &mgo.DialInfo{
		Addrs:     []string{url},
		Timeout:   time.Duration(timeout) * time.Second,
		Database:  dbName,
		PoolLimit: poolLimit,
	}
}

func (u Utilities) GetDBName() string {

	return u.GetStringConfigValue("db.db_name")
}

func (u Utilities) GetError(
	uuid exception.UUID, errKey string, err error) *exception.ASError {

	errMsg := u.GetErrorMessage(errKey)
	asErr := exception.NewASError(uuid, errMsg, err)
	log.Println(asErr.ErrorMessage())

	return asErr
}

/*
	Private methods
*/

func (u Utilities) getActiveEnvPrefix() string {
	env := u.GetActiveEnv()
	envMap := u.GetMapArrConfigValue("env-map")

	switch env {
	case toIntFromInt64Inteface(envMap[0]["index"]):
		return envMap[0]["type"].(string)
	case toIntFromInt64Inteface(envMap[1]["index"]):
		return envMap[1]["type"].(string)
	default:
		return ""
	}
}

func toArray(i interface{}) []interface{} {
	arr, ok := i.([]interface{})
	if !ok {
		// TODO: add error handling later
	}
	return arr
}

func toMap(i interface{}) map[string]interface{} {
	m, ok := i.(map[string]interface{})
	if !ok {
		// TODO: add error handling later
	}
	return m
}

func toArrayMap(i interface{}) []map[string]interface{} {
	arr := toArray(i)

	mapArr := make([]map[string]interface{}, len(arr))

	for i := 0; i < len(arr); i++ {
		mapArr[i] = toMap(arr[i])
	}

	return mapArr
}

func toIntFromInt64Inteface(i interface{}) int {

	intString := fmt.Sprintf("%d", i)
	integer, _ := strconv.ParseInt(intString, 10, 0)

	return int(integer)
}
