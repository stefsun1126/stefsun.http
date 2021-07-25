package configuration

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type config struct {
	viper *viper.Viper
}

// 要找的config名稱
func (c *config) setConfigName(in string) {
	c.viper.SetConfigName(in)
}

// 設定viper找配置文件的路徑
func (c *config) addConfigPath(in string) {
	c.viper.AddConfigPath(in)
}

// 查找並讀取config
func (c *config) readInConfig() error {
	if err := c.viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

// 從viper讀取值並轉為string
func (c *config) GetString(key string) string {
	return c.viper.GetString(key)
}

// 返回config instance
func New(filename string, env ...string) *config {
	config := &config{
		viper: viper.New(),
	}

	if len(env) > 0 && env[0] != "" {
		filename = filename + fmt.Sprintf(".%s", strings.ToLower(env[0]))
	}

	// Viper check ENV variables for all.
	config.viper.AutomaticEnv()
	// 重寫key
	config.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// 找類型是yaml的config
	config.viper.SetConfigType("yaml")
	// 設定viper找配置文件的路徑
	config.addConfigPath("./configuration")
	// 要找的config名稱
	config.setConfigName(filename)
	// // 查找並讀取config
	config.readInConfig()
	return config
}
