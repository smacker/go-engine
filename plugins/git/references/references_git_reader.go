package references

import (
	"strings"

	"github.com/chrislusf/gleam/util"
	git "gopkg.in/src-d/go-git.v4"
	storer "gopkg.in/src-d/go-git.v4/plumbing/storer"
)

type ReferencesGitReader struct {
	repositoryID string
	refs         storer.ReferenceIter
}

func New(r *git.Repository) *ReferencesGitReader {

	refs, _ := r.References()
	remotes, _ := r.Remotes()

	urls := remotes[0].Config().URLs
	repositoryID := strings.TrimPrefix(urls[0], "https://")

	return &ReferencesGitReader{
		repositoryID: repositoryID,
		refs:         refs,
	}
}

func (r *ReferencesGitReader) ReadHeader() (fieldNames []string, err error) {
	fieldNames = []string{
		"repositoryID",
		"refHash",
		"refName",
	}
	return fieldNames, nil
}

/*
root
 |-- repository_id: string (nullable = false)
 |-- name: string (nullable = false)
 |-- hash: string (nullable = false)
 |-- is_remote: boolean (nullable = false)
*/

func (r *ReferencesGitReader) Read() (row *util.Row, err error) {

	ref, err := r.refs.Next()
	if err != nil {
		return nil, err
	}

	return util.NewRow(util.Now(), r.repositoryID, ref.Hash().String(), ref.Name().String()), nil
}