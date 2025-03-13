package googlesheet

type GoogleSheetConfig struct {
	SecretPath string
}

type GoogleSheet struct {
	config GoogleSheetConfig
}

func (s *GoogleSheet) Write() error { return nil }
