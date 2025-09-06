package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/traP-jp/1m25_10/backend/internal/domain"
)

type ImageRepository interface {
	GetImages(ctx context.Context, token string) ([]domain.Image, error)
}

type traqMessage struct {
	ID        string         `json:"id"`
	UserID    string         `json:"userId"`
	Content   string         `json:"content"`
	CreatedAt string         `json:"createdAt"`
	UpdatedAt string         `json:"updatedAt"`
	Embedded  []traqEmbedded `json:"embedded"`
}

type traqEmbedded struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Raw  string `json:"raw"`
}

type traqMessagesResponse struct {
	Messages []traqMessage `json:"messages"`
	HasMore  bool          `json:"hasMore"`
}

func (r *sqlRepositoryImpl) GetImages(ctx context.Context, token string) ([]domain.Image, error) {
	// traQ APIでhasImage=trueのメッセージを検索
	messages, err := r.searchTraqMessages(token, true)
	if err != nil {
		return nil, err
	}

	// メッセージから画像IDを抽出
	imageIDs := extractImageIDs(messages)

	// 画像情報を構築（仮の実装）
	var images []domain.Image
	for _, id := range imageIDs {
		// 実際にはメッセージ情報から適切に構築
		image := domain.Image{
			ID:      id,
			Creator: uuid.Nil, // 仮
			Post: domain.Post{
				ID:      uuid.Nil, // 仮
				Content: "",       // 仮
			},
		}
		images = append(images, image)
	}

	return images, nil
}

// searchTraqMessages はtraQ APIでメッセージを検索するヘルパー関数
func (r *sqlRepositoryImpl) searchTraqMessages(token string, hasImage bool) ([]traqMessage, error) {
	baseURL := "https://q.trap.jp/api/v3/messages"
	params := url.Values{}
	if hasImage {
		params.Add("hasImage", "true")
	}

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("traQ API returned status %d", resp.StatusCode)
	}

	var response traqMessagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Messages, nil
}

// extractImageIDs はメッセージから画像のUUIDを抽出するヘルパー関数
func extractImageIDs(messages []traqMessage) []uuid.UUID {
	var imageIDs []uuid.UUID
	for _, msg := range messages {
		for _, embedded := range msg.Embedded {
			if embedded.Type == "file" {
				// embedded.Raw からファイル情報を解析
				// traQの埋め込みファイル形式を仮定
				if id, err := uuid.Parse(embedded.ID); err == nil {
					imageIDs = append(imageIDs, id)
				}
			}
		}
	}
	return imageIDs
}
