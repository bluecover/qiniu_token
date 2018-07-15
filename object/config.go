package object

type qiniuPersistentOps struct {
	Pfop       string `mapstructure:"pfop"`
	SaveBucket string `mapstructure:"save_bueket"`
	SaveKey    string `mapstructure:"save_key"`
}

type qiniuCategory struct {
	Bucket             string                        `mapstructure:"bucket"`
	SaveKey            string                        `mapstructure:"save_key"`
	Scope              string                        `mapstructure:"scope"`
	IsPrefixalScope    int64                         `mapstructure:"is_prefixal_scope"`
	MimeLimit          string                        `mapstructure:"mime_limit"`
	FsizeLimit         int64                         `mapstructure:"fsize_limit"`
	FsizeMin           int64                         `mapstructure:"fsize_min"`
	InsertOnly         int64                         `mapstructure:"insert_only"`
	PersistentOps      map[string]qiniuPersistentOps `mapstructure:"persistent_ops"`
	PersistentPipeline string                        `mapstructure:"persistent_pipeline"`
	ReturnBody         []string                      `mapstructure:"return_body"`
}

type qiniuConfig struct {
	AccessKey          string                   `mapstructure:"access_key"`
	SecretKey          string                   `mapstructure:"secret_key"`
	TokenDuration      int64                    `mapstructure:"token_duration"`
	PrivateURLDuration int64                    `mapstructure:"private_url_duration"`
	Domain             map[string]string        `mapstructure:"domain"`
	Category           map[string]qiniuCategory `mapstructure:"category"`
}
