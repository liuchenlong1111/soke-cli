package core
              
import (
    "encoding/json"
    "os"
    "path/filepath"
)

// GetConfigPath 返回配置文件路径
func GetConfigPath() string {
    home, _ := os.UserHomeDir()
    return filepath.Join(home, ".soke-cli", "config.json")
}

// LoadConfig 加载配置
func LoadConfig() (*CliConfig, error) {
    path := GetConfigPath()                                                                                                                                                                                                                    
    data, err := os.ReadFile(path)
    if err != nil {                                                                                                                                                                                                                            
        return nil, err
    }

    var cfg CliConfig
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}

// SaveConfig 保存配置
func SaveConfig(cfg *CliConfig) error {
    path := GetConfigPath()
    dir := filepath.Dir(path)

    if err := os.MkdirAll(dir, 0700); err != nil {                                                                                                                                                                                             
        return err
    }                                                                                                                                                                                                                                          
                
    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(path, data, 0600)
}