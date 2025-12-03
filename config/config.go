package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"

	"github.com/u00io/gazer_node/localstorage"
	"github.com/u00io/gazer_node/utils"
)

type Config struct {
	Units []*ConfigUnit
}

type ConfigUnit struct {
	Id         string
	Type       string
	PrivateKey string
	PublicKey  string
	Parameters map[string]string
}

var instance Config

func AddUnit(unitType string, parameters map[string]string) string {
	key := utils.NewKey()

	unit := &ConfigUnit{
		Id:         generateId(),
		Type:       unitType,
		Parameters: parameters,
		PrivateKey: key.GetPrivateKey(),
		PublicKey:  key.GetPublicKey(),
	}
	instance.Units = append(instance.Units, unit)
	Save()
	return unit.Id
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
