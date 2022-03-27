package types

type Response map[string]interface{}

type ReqestParams map[string]string

type User struct {
	PublicKey string `yaml:"public-key"`
	SecretKey string `yaml:"secret-key"`
}
