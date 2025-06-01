package support

var (
	RelativePath = ""
	RootPath     = ""

	EnvFilePath          = ".env"
	EnvFileEncryptPath   = ".env.encrypted"
	EnvFileEncryptCipher = "AES-256-CBC"

	DontVerifyEnvFileExists    = false
	DontVerifyEnvFileWhiteList = []string{"key:generate", "env:decrypt"}
)
