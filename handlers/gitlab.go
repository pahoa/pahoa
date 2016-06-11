package handlers

import (
	"errors"
	"log"
	"regexp"

	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"

	"github.com/pahoa/pahoa/core"
)

var gitClient *gitlab.Client

const GITLAB_CONCURRENT = 10

func client(config *viper.Viper) *gitlab.Client {
	if gitClient != nil {
		return gitClient
	}

	git := gitlab.NewClient(nil, config.GetString("gitlab.token"))

	gitlabBaseURL := config.GetString("gitlab.url")
	if gitlabBaseURL != "" {
		git.SetBaseURL(gitlabBaseURL)
	}

	gitClient = git

	return gitClient
}

type MergeCardBranchesTaskRunner struct {
	git           *gitlab.Client
	regexCardName *regexp.Regexp
}

func (m MergeCardBranchesTaskRunner) Run(item interface{}) error {
	proj := item.(*gitlab.Project)
	log.Printf("Gitlab - Finding branches for repostory %d", proj.ID)

	merges, err := listOpenedMergeRequests(m.git, proj.ID)
	if err != nil {
		log.Printf("Gitlab - Error listing merge requests %#v\n", err)
		return err
	}

	matchBranches := map[string]bool{}
	for _, merge := range merges {
		if merge.TargetBranch == "develop" && m.regexCardName.Match([]byte(merge.SourceBranch)) {
			matchBranches[merge.SourceBranch] = true
		}
	}

	branches, err := listBranches(m.git, proj.ID)
	if err != nil {
		log.Printf("Gitlab - Error listing branches %#v\n", err)
		return err
	}

	log.Printf("Gitlab - Found %d branches for repostory %d", len(branches), *proj.ID)

	for _, branch := range branches {
		if matchBranches[branch.Name] || !m.regexCardName.Match([]byte(branch.Name)) {
			continue
		}

		opts := &gitlab.CreateMergeRequestOptions{
			Title:           "Merge branch '" + branch.Name + "' into 'develop'",
			SourceBranch:    branch.Name,
			TargetBranch:    "develop",
			TargetProjectID: *proj.ID,
		}

		_, _, err := m.git.MergeRequests.CreateMergeRequest(*proj.ID, opts)
		if err != nil {
			log.Printf("Gitlab - Error creating merge request %#v\n", err)
			return err
		}

		log.Printf("Gitlab - Created merge request '" + branch.Name + "' into 'develop'")
	}
	return nil
}

func GitlabCreateMergeRequestToDevelop(config *viper.Viper, card *core.Card) error {
	log.Printf("Gitlab - Create merge request to develop - card: %s", card.ID)
	git := client(config)

	q := core.NewQueue(
		GITLAB_CONCURRENT,
		MergeCardBranchesTaskRunner{
			git:           git,
			regexCardName: regexp.MustCompile(`^` + card.ID + `\/.*$`),
		},
	).Run()

	// seed the queue
	go listAllProjects(git, q)

	if err := q.WaitWorkers(); err != nil {
		log.Printf("Gitlab - Error %#v\n", err)
		return err
	}

	log.Printf("Gitlab - Finished - card: %s", card.ID)
	return nil
}

func listAllProjects(git *gitlab.Client, q *core.Queue) {
	user, _, _ := git.Users.CurrentUser()
	if user == nil {
		q.Errors <- errors.New("Not found user by Gitlab token configuration")
		close(q.Jobs)
		return
	}

	options := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
		},
	}

	nextPage := true

	for i := 1; nextPage; i++ {
		options.Page = i

		projects, _, err := listProjects(git, user.IsAdmin, options)
		if err != nil {
			q.Errors <- err
			close(q.Jobs)
			return
		}

		count := len(projects)

		switch true {
		case count == 0:
			nextPage = false
		case count > 0:
			for _, proj := range projects {
				if proj.DefaultBranch == nil {
					continue
				}
				q.Jobs <- proj
			}
		}
	}
	close(q.Jobs)
}

func listProjects(git *gitlab.Client, isAdmin bool, options *gitlab.ListProjectsOptions) ([]*gitlab.Project, *gitlab.Response, error) {
	if isAdmin {
		return git.Projects.ListAllProjects(options)
	} else {
		return git.Projects.ListProjects(options)
	}
}

func listBranches(git *gitlab.Client, pid *int) ([]*gitlab.Branch, error) {
	result, _, err := git.Branches.ListBranches(*pid)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func listOpenedMergeRequests(git *gitlab.Client, pid *int) ([]*gitlab.MergeRequest, error) {
	opts := &gitlab.ListMergeRequestsOptions{
		State: "opned",
		ListOptions: gitlab.ListOptions{
			PerPage: 999999, // TODO: paging
		},
	}

	result, _, err := git.MergeRequests.ListMergeRequests(*pid, opts)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func init() {
	Register(core.ActionCreateMergeRequestToDevelop, GitlabCreateMergeRequestToDevelop)
}
