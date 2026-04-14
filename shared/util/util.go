package util

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetRandomAvatar returns a random avatar URL from the randomuser.me API
func GetRandomAvatar(index int) string {
	return fmt.Sprintf("https://randomuser.me/api/portraits/lego/%d.jpg", index)
}

func GenRandomID() string {
	return primitive.NewObjectID().Hex()
}
