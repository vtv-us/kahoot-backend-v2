package utils

func GenLink(frontend string, groupID string, inviterID string) string {
	return frontend + "/group/invite/" + groupID + "/" + inviterID
}
