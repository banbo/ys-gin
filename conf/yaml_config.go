package conf

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// YamlConfiger 使用 gopkg.in/yaml.v3 解析 YAML 配置
type YamlConfiger struct {
	data map[string]interface{}
	sync.RWMutex
}

// NewYamlConfiger 从文件创建 YAML 配置
func NewYamlConfiger(filename string) (*YamlConfiger, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return NewYamlConfigerData(buf)
}

// NewYamlConfigerData 从字节数据创建 YAML 配置
func NewYamlConfigerData(buf []byte) (*YamlConfiger, error) {
	var data map[string]interface{}
	if err := yaml.Unmarshal(buf, &data); err != nil {
		return nil, err
	}
	return &YamlConfiger{data: data}, nil
}

// getData 按 "." 分隔的 key 路径获取值
func (c *YamlConfiger) getData(key string) (interface{}, error) {
	if key == "" {
		return nil, errors.New("key is empty")
	}
	c.RLock()
	defer c.RUnlock()

	keys := strings.Split(key, ".")
	tmpData := c.data
	for i, k := range keys {
		v, ok := tmpData[k]
		if !ok {
			return nil, fmt.Errorf("not exist key %q", key)
		}
		if i == len(keys)-1 {
			return v, nil
		}
		switch m := v.(type) {
		case map[string]interface{}:
			tmpData = m
		default:
			return nil, fmt.Errorf("not exist key %q", key)
		}
	}
	return nil, fmt.Errorf("not exist key %q", key)
}

func (c *YamlConfiger) Set(key, val string) error {
	c.Lock()
	defer c.Unlock()
	keys := strings.Split(key, ".")
	tmpData := c.data
	for i, k := range keys {
		if i == len(keys)-1 {
			tmpData[k] = val
			return nil
		}
		if v, ok := tmpData[k]; ok {
			if m, ok := v.(map[string]interface{}); ok {
				tmpData = m
			} else {
				return fmt.Errorf("key %q is not a section", key)
			}
		} else {
			newMap := make(map[string]interface{})
			tmpData[k] = newMap
			tmpData = newMap
		}
	}
	return nil
}

func (c *YamlConfiger) String(key string) string {
	v, err := c.getData(key)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func (c *YamlConfiger) Strings(key string) []string {
	v := c.String(key)
	if v == "" {
		return nil
	}
	return strings.Split(v, ";")
}

func (c *YamlConfiger) Int(key string) (int, error) {
	v, err := c.getData(key)
	if err != nil {
		return 0, err
	}
	switch vv := v.(type) {
	case int:
		return vv, nil
	case int64:
		return int(vv), nil
	case float64:
		return int(vv), nil
	case string:
		return strconv.Atoi(vv)
	}
	return 0, fmt.Errorf("not int value for key %q", key)
}

func (c *YamlConfiger) Int64(key string) (int64, error) {
	v, err := c.getData(key)
	if err != nil {
		return 0, err
	}
	switch vv := v.(type) {
	case int:
		return int64(vv), nil
	case int64:
		return vv, nil
	case float64:
		return int64(vv), nil
	case string:
		return strconv.ParseInt(vv, 10, 64)
	}
	return 0, fmt.Errorf("not int64 value for key %q", key)
}

func (c *YamlConfiger) Bool(key string) (bool, error) {
	v, err := c.getData(key)
	if err != nil {
		return false, err
	}
	switch vv := v.(type) {
	case bool:
		return vv, nil
	case string:
		return strconv.ParseBool(vv)
	}
	return false, fmt.Errorf("not bool value for key %q", key)
}

func (c *YamlConfiger) Float(key string) (float64, error) {
	v, err := c.getData(key)
	if err != nil {
		return 0, err
	}
	switch vv := v.(type) {
	case float64:
		return vv, nil
	case int:
		return float64(vv), nil
	case int64:
		return float64(vv), nil
	case string:
		return strconv.ParseFloat(vv, 64)
	}
	return 0, fmt.Errorf("not float64 value for key %q", key)
}

func (c *YamlConfiger) DefaultString(key string, defaultVal string) string {
	v := c.String(key)
	if v == "" {
		return defaultVal
	}
	return v
}

func (c *YamlConfiger) DefaultStrings(key string, defaultVal []string) []string {
	v := c.Strings(key)
	if v == nil {
		return defaultVal
	}
	return v
}

func (c *YamlConfiger) DefaultInt(key string, defaultVal int) int {
	v, err := c.Int(key)
	if err != nil {
		return defaultVal
	}
	return v
}

func (c *YamlConfiger) DefaultInt64(key string, defaultVal int64) int64 {
	v, err := c.Int64(key)
	if err != nil {
		return defaultVal
	}
	return v
}

func (c *YamlConfiger) DefaultBool(key string, defaultVal bool) bool {
	v, err := c.Bool(key)
	if err != nil {
		return defaultVal
	}
	return v
}

func (c *YamlConfiger) DefaultFloat(key string, defaultVal float64) float64 {
	v, err := c.Float(key)
	if err != nil {
		return defaultVal
	}
	return v
}

func (c *YamlConfiger) DIY(key string) (interface{}, error) {
	return c.getData(key)
}

func (c *YamlConfiger) GetSection(section string) (map[string]string, error) {
	c.RLock()
	defer c.RUnlock()

	v, ok := c.data[section]
	if !ok {
		return nil, fmt.Errorf("section %q not found", section)
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("section %q is not a map", section)
	}
	result := make(map[string]string, len(m))
	for k, val := range m {
		result[k] = fmt.Sprintf("%v", val)
	}
	return result, nil
}

func (c *YamlConfiger) SaveConfigFile(filename string) error {
	c.RLock()
	defer c.RUnlock()

	buf, err := yaml.Marshal(c.data)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, buf, 0644)
}
