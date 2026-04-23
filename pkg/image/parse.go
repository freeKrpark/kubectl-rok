package image

import "strings"

func ParseDigest(imageID string) string {
	if idx := strings.Index(imageID, "@"); idx != -1 {
		return imageID[idx+1:]
	}
	return ""
}

func ParseImage(image string) (repository, tag string) {
	searchFrom := 0
	if idx := strings.LastIndex(image, "/"); idx != -1 {
		searchFrom = idx
	}

	if colon := strings.Index(image[searchFrom:], ":"); colon != -1 {
		pos := searchFrom + colon
		return image[:pos], image[pos+1:]
	}

	return image, "latest"
}

func ShortImage(image string) string {
	return ImageBasename(image)
}

func ImageBasename(image string) string {
	if idx := strings.LastIndex(image, "/"); idx != -1 {
		return image[idx+1:]
	}
	return image
}
