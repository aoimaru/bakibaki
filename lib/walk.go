package lib

// func WalkGitLog(client Client, hash string) () {
// 	GoBuffer, err := client.GetGitObject(hash)
// 	Header, Content, err := Header3Content(&GoBuffer)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	if  !strings.HasPrefix(string(Header), "commit") {
// 		return
// 	}
// 	commit, err := CreateCommitObject(Header, Content)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	if len(commit.Parents) <= 0 {
// 		return
// 	}

// 	for _, parent := range commit.Parents {
// 		fmt.Println(parent.Hash)
// 		WalkGitLog(client, parent.Hash)

// 	}

// }
