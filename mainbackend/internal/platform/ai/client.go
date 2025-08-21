// mainbackend/internal/platform/ai/client.go

package ai

// IAIClient, AI istemcilerimiz (Ollama, OpenAI vb.) için ortak arayüzü tanımlar.
// Bu arayüz, servis katmanının hangi AI istemcisinin kullanıldığını bilmeden
// işlem yapabilmesini sağlar.
type IAIClient interface {
	// RefineTranscription, ham metni alır ve AI modeli ile düzenleyip döndürür.
	RefineTranscription(rawText string) (string, error)
}
