package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	git "github.com/libgit2/git2go/v34"
)

func main() {
	pathRe := regexp.MustCompile(`^/(?:@(.*?)/)?(.*)`)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	r := os.Getenv("REPO")
	if len(r) == 0 {
		r = "/src"
	}

	repo, err := git.OpenRepository(r)
	if err != nil {
		fmt.Print("error: ")
		fmt.Println(err)
		return
	}

	fmt.Println("starting server on port " + port)
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		match := pathRe.FindStringSubmatch(req.URL.Path)
		path := match[2]
		if len(path) == 0 || strings.HasSuffix(path, "/") {
			path = path + "index.html"
		}

		blob, err := getFile(repo, match[1], path)
		if err != nil {
			fmt.Printf("error(%v): %v\n", path, err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(blob)
	})
	http.ListenAndServe(":"+port, nil)
}

// ref: empty for head, branch name, or commit hash (slice)
func getFile(repo *git.Repository, ref, path string) ([]byte, error) {
	if len(ref) == 0 {
		head, _ := repo.Head()
		return getFromTree(repo, path, head.Target())
	}

	branch, err := repo.LookupBranch("origin/"+ref, git.BranchRemote)
	if err == nil {
		return getFromTree(repo, path, branch.Target())
	}

	// revparse needed to resolve shorter hashes
	rev, err := repo.RevparseSingle(ref)
	if err != nil {
		return make([]byte, 0), err
	}
	oid := rev.Id()
	return getFromTree(repo, path, oid)
}

func getFromTree(repo *git.Repository, path string, oid *git.Oid) ([]byte, error) {
	commit, _ := repo.LookupCommit(oid)
	tree, _ := commit.Tree()

	entry, err := tree.EntryByPath(path)
	if err != nil {
		return make([]byte, 0), err
	}

	blob, err := repo.LookupBlob(entry.Id)
	return blob.Contents(), nil
}
