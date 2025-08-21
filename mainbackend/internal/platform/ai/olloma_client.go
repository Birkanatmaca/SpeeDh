// mainbackend/internal/platform/ai/ollama_client.go

package ai

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
)

// ollamaClient, IAIClient arayüzünü uygular
type ollamaClient struct {
	client *api.Client
}

// NewOllamaClient, yeni bir istemci oluşturur ve IAIClient arayüzü olarak döndürür.
func NewOllamaClient() IAIClient {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		panic(fmt.Sprintf("Ollama istemcisi oluşturulamadı: %v", err))
	}
	return &ollamaClient{client: client}
}

// RefineTranscription metodu, IAIClient arayüzünün bir parçasıdır.
func (c *ollamaClient) RefineTranscription(rawText string) (string, error) {
	prompt := fmt.Sprintf(
		"Sen yetenekli bir editörsün. Aşağıdaki metin, bir ses kaydından dönüştürülmüş ham bir transkripttir. "+
			"Görevin; yazım ve dil bilgisi hatalarını düzeltmek, uygun noktalama işaretlerini (virgül, nokta, soru işareti) eklemek "+
			"ve metni temiz, okunabilir paragraflar halinde biçimlendirmektir. Metnin anlamını özetleme veya değiştirme. "+
			"Metnin orijinal dili Türkçe'dir. Lütfen düzeltilmiş metni Türkçe olarak sun. İşte transkript:\n\n---\n\n%s",
		rawText,
	)

	req := &api.ChatRequest{
		Model: "deepseek-coder:6.7b",
		Messages: []api.Message{
			{Role: "user", Content: prompt},
		},
	}

	var fullResponse string
	respFn := func(res api.ChatResponse) error {
		fullResponse += res.Message.Content
		return nil
	}

	err := c.client.Chat(context.Background(), req, respFn)
	if err != nil {
		return "", fmt.Errorf("Ollama Chat hatası: %w", err)
	}

	if fullResponse == "" {
		return "", fmt.Errorf("Ollama'dan geçerli bir yanıt alınamadı")
	}

	return fullResponse, nil
}
