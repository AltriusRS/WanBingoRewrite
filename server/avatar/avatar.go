package avatar

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
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
	log.Printf("AVATAR: Starting upload for Discord user %s, avatar hash: %s, player ID: %s", discordUserID, avatarHash, playerID)
	if avatarHash == "" {
		log.Printf("AVATAR: No avatar hash provided, skipping upload")
		return "", nil // No avatar to upload
	}

	// Download avatar from Discord CDN
	// Discord avatars can be PNG, GIF, or WebP. If hash starts with 'a_', it's animated (GIF)
	extension := "png"
	if strings.HasPrefix(avatarHash, "a_") {
		extension = "gif"
	}
	avatarURL := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.%s", discordUserID, avatarHash, extension)
	log.Printf("AVATAR: Downloading from URL: %s", avatarURL)
	resp, err := http.Get(avatarURL)
	if err != nil {
		log.Printf("AVATAR: Failed to download avatar: %v", err)
		return "", fmt.Errorf("failed to download avatar: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("AVATAR: Discord response status: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		log.Printf("AVATAR: Failed to download avatar, status: %d", resp.StatusCode)
		return "", fmt.Errorf("failed to download avatar: status %d", resp.StatusCode)
	}

	avatarData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("AVATAR: Failed to read avatar data: %v", err)
		return "", fmt.Errorf("failed to read avatar data: %w", err)
	}
	log.Printf("AVATAR: Downloaded %d bytes of avatar data", len(avatarData))

	// Upload to R2
	contentType := fmt.Sprintf("image/%s", extension)
	key := fmt.Sprintf("avatars/%s.%s", playerID, extension)

	err = uploadToR2(key, avatarData, contentType)
	if err != nil {
		log.Printf("AVATAR: Failed to upload to R2: %v", err)
		return "", fmt.Errorf("failed to upload avatar to R2: %w", err)
	}
	log.Printf("AVATAR: Successfully uploaded to R2")

	// Return just the key (not full URL) for database flexibility
	log.Printf("AVATAR: Successfully uploaded, returning key: %s", key)
	return key, nil
}

// uploadToR2 uploads data to R2 bucket
func uploadToR2(key string, data []byte, contentType string) error {
	log.Printf("R2: Starting upload for key: %s, size: %d bytes", key, len(data))

	// Load AWS config for R2
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("auto"), // R2 uses "auto"
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			// Use access key + secret (standard R2 authentication)
			accessKey := os.Getenv("R2_ACCESS_KEY_ID")
			secretKey := os.Getenv("R2_SECRET_ACCESS_KEY")
			log.Printf("R2: Access key present: %v (%d chars), Secret key present: %v (%d chars)",
				accessKey != "", len(accessKey), secretKey != "", len(secretKey))

			if accessKey == "" || secretKey == "" {
				return aws.Credentials{}, fmt.Errorf("R2_ACCESS_KEY_ID and R2_SECRET_ACCESS_KEY must be set")
			}

			// Check if access key is the right length (should be 32 chars)
			if len(accessKey) != 32 {
				log.Printf("R2: WARNING - Access key length is %d, expected 32. This may cause authentication issues.", len(accessKey))
			}

			return aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			}, nil
		})),
	)
	if err != nil {
		log.Printf("R2: Failed to load AWS config: %v", err)
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with R2 endpoint
	accountID := os.Getenv("R2_ACCOUNT_ID")
	log.Printf("R2: Account ID present: %v", accountID != "")
	if accountID == "" {
		return fmt.Errorf("R2_ACCOUNT_ID environment variable not set")
	}

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)
	log.Printf("R2: Using endpoint: %s", endpoint)
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true // R2 requires path style
	})

	// Upload object
	bucket := os.Getenv("R2_BUCKET")
	log.Printf("R2: Bucket present: %v, bucket name: %s", bucket != "", bucket)
	if bucket == "" {
		return fmt.Errorf("R2_BUCKET environment variable not set")
	}

	log.Printf("R2: Attempting to upload to bucket: %s, key: %s", bucket, key)
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPublicRead, // Make it publicly accessible
	})

	if err != nil {
		log.Printf("R2: PutObject failed: %v", err)
		return err
	}

	log.Printf("R2: Successfully uploaded to R2")
	return nil
}

// GetAvatarURL constructs the full CDN URL from an avatar key
func GetAvatarURL(avatarKey string) string {
	if avatarKey == "" {
		return ""
	}

	cdnURL := os.Getenv("CDN_URL")
	if cdnURL == "" {
		return "" // Return empty if CDN not configured
	}

	// Ensure CDN URL ends with /
	if !strings.HasSuffix(cdnURL, "/") {
		cdnURL += "/"
	}

	return cdnURL + avatarKey
}

// ExtractAvatarKey extracts just the key from a full CDN URL (for migration)
func ExtractAvatarKey(fullURL string) string {
	if fullURL == "" {
		return ""
	}

	cdnURL := os.Getenv("CDN_URL")
	if cdnURL == "" {
		return fullURL // If no CDN configured, return as-is
	}

	// Ensure CDN URL ends with /
	if !strings.HasSuffix(cdnURL, "/") {
		cdnURL += "/"
	}

	// If the URL starts with our CDN URL, extract just the key
	if strings.HasPrefix(fullURL, cdnURL) {
		return strings.TrimPrefix(fullURL, cdnURL)
	}

	// Otherwise return as-is (might be a legacy URL or different format)
	return fullURL
}
