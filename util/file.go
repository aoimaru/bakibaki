package util

func GetGitObjectHeader(buffer *[]byte) (string, error) {

	header := make([]byte, 0)
	for _, buf := range *buffer {
		if buf == 0 {
			break
		}
		header = append(header, buf)
	}
	return string(header), nil

}
