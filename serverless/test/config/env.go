package testConfig

var DbEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_NAME":     "cc2_test",
	"DB_PASSWORD": "postgres",
	"DB_PORT":     "5432",
	"DB_USER":     "postgres",
}

var AwsEnv = map[string]string{
	"AWS_SES_REGION": "us-east-1",
}

var EmailEnv = map[string]string{
	"SOURCE_EMAIL_ADDRESS": "test@gpalab.digital",
	"EMAIL_REDIRECT_URL":   "https://example.com",
}

var AprimoEnv = map[string]string{
	"APRIMO_TENANT":        "state-sb1",
	"APRIMO_CLIENT_ID":     "4MS42XXC-4MS4",
	"APRIMO_CLIENT_SECRET": "PW9pdFcT7XJ3gfgaeb72EhVcjmaPxckhyvn3fHVg",
}
