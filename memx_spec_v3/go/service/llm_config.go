package service

import "memx/db"

// ConfigureLLMsFromEnv は環境変数から OpenAI クライアントを読み込み、要約系 LLM を接続する。
// OPENAI_API_KEY が未設定なら何もしない。
func (s *Service) ConfigureLLMsFromEnv() error {
	client, ok, err := db.NewOpenAIClientFromEnv()
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	s.SetMiniLLM(client)
	s.SetReflectLLM(client)
	return nil
}
