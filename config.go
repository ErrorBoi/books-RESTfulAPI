package main

import "os"

// CreateConfig writes configs in environment variables
func CreateConfig() {
	os.Setenv("PORT", "7878")
}
