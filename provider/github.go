package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
	// . "github.com/go-git/go-git/v5/_examples"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-resty/resty/v2"
	"github.com/maahsome/changelog-pr/common"
)

// Github - Structure to hold stuff
type Github struct {
	Provider string
}

var lastTag *plumbing.Reference
var numRegex = regexp.MustCompile(`#(\d+) from`)

type PRBody struct {
	Body string `json:"body"`
}

func parsePRNumber(msg string) (uint, error) {
	matches := numRegex.FindAllStringSubmatch(msg, 1)
	if len(matches) == 0 || len(matches[0]) < 2 {
		return 0, fmt.Errorf("could not find PR number in commit message")
	}
	u64, err := strconv.ParseUint(matches[0][1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse PR number %q from commit message: %v", matches[0][1], err)
	}
	return uint(u64), nil
}

// GetChangeLogSincePR - Get the changelog details from the PR description
func (p *Github) GetChangeLogFromPR(src string, sincePR string, release string, auth AuthToken, fileName string) (string, error) {
	path := src

	r, err := git.PlainOpen(path)
	if err != nil {
		return "", errors.New("failed generation of changelog")
	}

	tagrefs, err := r.Tags()
	if err != nil {
		return "", errors.New("failed generation of changelog")
	}

	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		if strings.HasSuffix(t.Name().String(), sincePR) {
			common.Logger.Info(t.Hash())
			lastTag = t
		}
		return nil
	})
	if err != nil {
		return "", errors.New("failed generation of changelog")
	}

	// Gets the HEAD history from HEAD, just like this command:
	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	if err != nil {
		return "", errors.New("failed generation of changelog")
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), Order: git.LogOrderCommitterTime})
	if err != nil {
		return "", errors.New("failed generation of changelog")
	}

	findingHash := true
	PRs := []string{}
	err = cIter.ForEach(func(c *object.Commit) error {
		if c.Hash == lastTag.Hash() {
			findingHash = false
			return nil
		}
		if findingHash {
			if strings.HasPrefix(c.Message, "Merge pull request #") {
				pr, err := parsePRNumber(strings.Split(c.Message, "\n")[0])
				if err != nil {
					common.Logger.WithError(err).Error("Bad PR Parse")
				}
				PRs = append(PRs, fmt.Sprintf("%d", pr))
				common.Logger.Info(fmt.Sprintf("%s %s\n", c.ID(), strings.Split(c.Message, "\n")[0]))
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return "", errors.New("failed generation of changelog")
	}

	// curl -sH "Accept: application/vnd.github.v3+json" https://api.github.com/repos/splicemachine/splicectl/pulls/5 | jq -r '.body'
	// Figure out how to AUTHENTICATE for PRIVATE REPOS
	restClient := resty.New()

	changeLog := common.Changelog{}
	changeLog.Version = release

	for _, p := range PRs {
		uri := fmt.Sprintf("https://api.github.com/repos/splicemachine/splicectl/pulls/%s", p)
		resp, resperr := restClient.R().
			SetHeader("Accept", "application/vnd.github.v3+json").
			SetHeader("Authorization", fmt.Sprintf("token %s", auth.GithubToken)).
			Get(uri)

		if resperr != nil {
			common.Logger.WithError(resperr).Error("Error getting PR")
		}

		var body PRBody

		marshErr := json.Unmarshal(resp.Body(), &body)
		if marshErr != nil {
			common.Logger.Error("Could not unmarshall data", marshErr)
		}

		err := common.ParseMarkdown(body.Body, p, &changeLog)
		if err != nil {
			common.Logger.Error("Could not parse the markdown")
		}
	}

	markdown, err := changeLog.Template()
	if err != nil {
		return "", errors.New("failed generation of changelog")
	}

	if len(fileName) > 0 {
		err := changeLog.WriteFile(fileName)
		if err != nil {
			return "", errors.New("failed to write to the output file")
		}
		return "Changelog data has been saved.", nil
	}

	return string(markdown[:]), nil
}
