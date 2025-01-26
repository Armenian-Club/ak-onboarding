package drive

import (
	"context"
	"fmt"

	"github.com/Armenian-Club/ak-onboarding/internal/config"
	"google.golang.org/api/drive/v3"
)

func (c *client) AddUser(ctx context.Context, gmail string) error {
	folderID := config.FolderID
	permission := &drive.Permission{
		Type:         "user",
		Role:         "writer", 
		EmailAddress: gmail, 
	}
	_, err := c.srv.Permissions.Create(folderID, permission).
		Fields("id").
		Context(ctx).
		Do()
	if err != nil {
		return fmt.Errorf("failed to add drive permission: %w", err)
	}

	return nil
}
