package core
              
// Identity 身份类型
type Identity string

const (
    AsUser Identity = "user"
    AsBot  Identity = "bot"
)

// CliConfig 配置结构
type CliConfig struct {
    AppID      string
    AppSecret  string
    APIBaseURL string
    UserToken  string
    BotToken   string
    CorpID     string
    DeptUserID string
}
                                                                                                                                                                                                                                                 
// APIRequest 通用请求结构
type APIRequest struct {
    Method string
    Path   string
    Query  map[string]interface{}
    Body   interface{}
    As     Identity
}