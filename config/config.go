package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/u00io/gazer_node/localstorage"
	"github.com/u00io/gazer_node/utils"
)

type Config struct {
	Units []*ConfigUnit `json:"units"`
}

type ConfigUnit struct {
	Id         string            `json:"id"`
	Type       string            `json:"type"`
	PrivateKey string            `json:"private_key"`
	PublicKey  string            `json:"public_key"`
	Parameters map[string]string `json:"parameters"`
	Translate  bool              `json:"translate"`
}

func NewConfigUnit() *ConfigUnit {
	var c ConfigUnit
	c.Parameters = make(map[string]string)
	return &c
}

var instance Config

func AddUnit(unitConfig *ConfigUnit) string {
	key := utils.NewKey()
	unitConfig.Id = generateId()
	unitConfig.PrivateKey = key.GetPrivateKey()
	unitConfig.PublicKey = key.GetPublicKey()
	instance.Units = append(instance.Units, unitConfig)
	Save()
	return unitConfig.Id
}

func RemoveUnit(unitId string) {
	for i, unit := range instance.Units {
		if unit.Id == unitId {
			instance.Units = append(instance.Units[:i], instance.Units[i+1:]...)
			Save()
			return
		}
	}
}

func generateId() string {
	rndBytes := make([]byte, 16)
	rand.Read(rndBytes)
	return hex.EncodeToString(rndBytes)
}

func configFileDirectory() string {
	return localstorage.Path()
}

func configFilePath() string {
	return configFileDirectory() + "/config.json"
}

func Save() error {
	dir := configFileDirectory()
	os.MkdirAll(dir, 0755)
	bs, _ := json.MarshalIndent(instance, "", "  ")
	return os.WriteFile(configFilePath(), bs, 0644)
}

func Load() error {
	bs, err := os.ReadFile(configFilePath())
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, &instance)
}

func Units() []*ConfigUnit {
	return instance.Units
}

func UnitById(unitId string) *ConfigUnit {
	for _, unit := range instance.Units {
		if unit.Id == unitId {
			return unit
		}
	}
	return nil
}

func (c *ConfigUnit) GetParameterBool(key string, defaultValue bool) bool {
	valueStr := c.GetParameterString(key, fmt.Sprint(defaultValue))
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *ConfigUnit) GetParameterInt64(key string, defaultValue int64) int64 {
	valueStr := c.GetParameterString(key, fmt.Sprint(defaultValue))
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *ConfigUnit) GetParameterFloat64(key string, defaultValue float64) float64 {
	valueStr := c.GetParameterString(key, fmt.Sprint(defaultValue))
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *ConfigUnit) GetParameterString(key string, defaultValue string) string {
	value, exists := c.Parameters[key]
	if !exists {
		return defaultValue
	}
	return value
}

func (c *ConfigUnit) SetParameterBool(key string, value bool) {
	c.SetParameterString(key, fmt.Sprint(value))
}

func (c *ConfigUnit) SetParameterInt64(key string, value int64) {
	c.SetParameterString(key, fmt.Sprint(value))
}

func (c *ConfigUnit) SetParameterFloat64(key string, value float64) {
	c.SetParameterString(key, fmt.Sprint(value))
}

func (c *ConfigUnit) SetParameterString(key string, value string) {
	c.Parameters[key] = value
}
