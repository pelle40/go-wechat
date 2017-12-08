package config
//版本
const VERSION="1.0"
//微信配置
type WechatItem struct{
	Appid string `json:"appid"`
	Appkey string `json:"appkey"`
	InterfaceToken string `json:"token"`
}

var WechatConfig  map[string]*WechatItem = map[string]*WechatItem{
		"top_learn":&WechatItem{"wxb52ac09696b62641","5df4e9762fdbc14e1caa2995122fe5cd","top_learn_token"}}


func init() {
}
