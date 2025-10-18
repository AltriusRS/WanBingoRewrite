package avatar

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// UploadDiscordAvatar downloads the Discord avatar and uploads it to R2, returning the CDN URL
func UploadDiscordAvatar(discordUserID, avatarHash, playerID string) (string, error) {
	if avatarHash == "" {
		return "", nil // No avatar to upload
	}

	// Download avatar from Discord CDN
	avatarURL := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", discordUserID, avatarHash)
	resp, err := http.Get(avatarURL)
	if err != nil {
		return "", fmt.Errorf("failed to download avatar: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download avatar: status %d", resp.StatusCode)
	}

	avatarData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read avatar data: %w", err)
	}

	// Upload to R2
	key := fmt.Sprintf("avatars/%s.png", playerID)

	err = uploadToR2(key, avatarData, "image/png")
	if err != nil {
		return "", fmt.Errorf("failed to upload avatar to R2: %w", err)
	}

	// Return CDN URL
	cdnURL := os.Getenv("CDN_URL")
	if cdnURL == "" {
		return "", fmt.Errorf("CDN_URL environment variable not set")
	}

	// Ensure CDN URL ends with /
	if !strings.HasSuffix(cdnURL, "/") {
		cdnURL += "/"
	}

	return cdnURL + key, nil
}

// uploadToR2 uploads data to R2 bucket
func uploadToR2(key string, data []byte, contentType string) error {
	// Load AWS config for R2
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("auto"), // R2 uses "auto"
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			accessKey := os.Getenv("R2_ACCESS_KEY_ID")
			secretKey := os.Getenv("R2_SECRET_ACCESS_KEY")
			if accessKey == "" || secretKey == "" {
				return aws.Credentials{}, fmt.Errorf("R2_ACCESS_KEY_ID and R2_SECRET_ACCESS_KEY must be set")
			}
			return aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			}, nil
		})),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with R2 endpoint
	accountID := os.Getenv("R2_ACCOUNT_ID")
	if accountID == "" {
		return fmt.Errorf("R2_ACCOUNT_ID environment variable not set")
	}

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true // R2 requires path style
	})

	// Upload object
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String("bingo"),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPublicRead, // Make it publicly accessible
	})

	return err
}
