package models

func GetAllComments() ([]string, error) {
	return Client.LRange("comments",0,-1).Result()
}

func AddComment(comment string) error {
	return Client.LPush("comments", comment).Err()
}
