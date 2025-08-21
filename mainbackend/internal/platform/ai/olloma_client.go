// mainbackend/internal/platform/ai/ollama_client.go
package ai

import (
	"context"
	"fmt"
	"github.com/ollama/ollama/api"
	"strings"
)

// ollamaClient, IAIClient arayüzünü uygular.
type ollamaClient struct {
	client *api.Client
}

// NewOllamaClient, yerel Ollama sunucusu için yeni bir istemci oluşturur.
func NewOllamaClient() IAIClient {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		panic(fmt.Sprintf("Ollama istemcisi oluşturulamadı: %v", err))
	}
	return &ollamaClient{client: client}
}

// RefineTranscription, ham metni alır ve Ollama modeline göndererek düzenlenmiş halini döndürür.
func (c *ollamaClient) RefineTranscription(rawText string) (string, error) {
	if rawText == "" {
		return "", nil // Boş metin geldiyse, boş döndür.
	}

	prompt := fmt.Sprintf(
		"Sen yetenekli bir Türkçe editörsün. Aşağıdaki metin, bir ses kaydından dönüştürülmüş ham bir transkripttir. "+
			"Görevin; sadece yazım ve dil bilgisi hatalarını düzeltmek, uygun noktalama işaretlerini (virgül, nokta, soru işareti) eklemek "+
			"ve metni temiz, okunabilir paragraflar haline biçimlendirmektir. Metnin anlamını veya kelimelerini kesinlikle değiştirme. "+
			"Eğer metin anlamsız veya boş ise, sadece boş bir metin döndür. Düzeltilmiş metni Türkçe olarak sun."+
			"Çıktı olarak sadece düzenlenmiş metini vermeni istiyorum. Yaptığın değişiklikleri ve düzenlemeleri görmek istemiyorum. İşte transkript: \n\n---\n\n%s",
		rawText,
	)

	req := &api.ChatRequest{
		Model: "deepseek-r1:8b", // Kullandığınız modelin adı (örn: llama3)
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
	closingTag := "</think>"
	if separatorIndex := strings.Index(fullResponse, closingTag); separatorIndex != -1 {
		// Eğer etiket varsa, ondan sonraki kısmı al
		cleanedText := fullResponse[separatorIndex+len(closingTag):]
		// Baştaki ve sondaki boşlukları temizle
		return strings.TrimSpace(cleanedText), nil
	}
	return fullResponse, nil
}
