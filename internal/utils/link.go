package utils

func GenLink(frontend string, groupID string) string {
	return frontend + "/group/invite/" + groupID
}
