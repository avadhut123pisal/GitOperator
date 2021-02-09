package githubService

import (
	"GitOperator/logger"
	"context"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// HARD-CODED PARAMETERS FOR DEMO

var (
	ctx                = context.Background()
	owner              = "avadhut123pisal"
	repo               = "GolangAdvancedConcepts"
	baseBranch         = "master"
	branch             = "newBranch"
	filesToBeModified  = "README.md"
	authorName         = "Avadhut Pisal"
	authorEmail        = "avadhutpisal47@gmail.com"
	commitMessage      = "First commit"
	pullReqTitle       = "Pull Request By Git Operator"
	pullReqDescription = "Git operator app has modified some files content"
)

// HandleGithubOperations handles the github repo operations
func HandleGithubOperations(accessToken string) error {
	// CREATE CLIENT
	client := createGithubClient(accessToken)
	// CREATE A BRNACH IF NOT EXISTS
	branchRef, err := getBranchReference(client)
	if err != nil {
		logger.ErrorLogger.Printf("ERROR OCCURED IN CREATING A NEW BRANCH: %v", err)
		return err
	}
	// MODIFY FIILES FROM NEW BRANCH
	tree, err := modifyFilesAndGetTree(branchRef, client)
	if err != nil {
		logger.ErrorLogger.Printf("ERROR OCCURED IN CREATING A TREE: %v", err)
		return err
	}
	// PUSH COMMIT / UPDATE REMOTE REPOSITORY CODE
	err = pushCommit(branchRef, tree, client)
	if err != nil {
		logger.ErrorLogger.Printf("ERROR IN UPDATING / COMMITING FILES ON REMOTE: %v", err)
		return err
	}
	err = createPullRequest(client)
	if err != nil {
		logger.ErrorLogger.Printf("ERROR IN CREATING PULL REQUEST: %v", err)
		return err
	}
	return nil
}

func createGithubClient(accessToken string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accessToken,
	})
	httpClient := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(httpClient)
	return client
}

func getBranchReference(client *github.Client) (ref *github.Reference, err error) {
	// TRY TO GET REFERENCE TO THE NEW BRANCH
	if ref, _, err = client.Git.GetRef(ctx, owner, repo, "refs/heads/"+branch); err == nil {
		return ref, nil
	}
	// IF ERR IS THERE, MEANS NO BRANCH IS PRESENT YET WITH THIS NAME
	// CREATE A NEW ONE FROM BASE BRANCH
	var baseRef *github.Reference
	if baseRef, _, err = client.Git.GetRef(ctx, owner, repo, "refs/heads/"+baseBranch); err != nil {
		return nil, err
	}
	newRef := &github.Reference{Ref: github.String("refs/heads/" + branch), Object: &github.GitObject{SHA: baseRef.Object.SHA}}
	ref, _, err = client.Git.CreateRef(ctx, owner, repo, newRef)
	if err != nil {
		logger.ErrorLogger.Printf("ERROR OCCURED IN CREATING A NEW BRANCH: %v", err)
		return nil, err
	}
	return ref, err
}

func modifyFilesAndGetTree(ref *github.Reference, client *github.Client) (tree *github.Tree, err error) {
	// CREATE A TREE FOR THE FILES TO BE COMMITED
	entries := []*github.TreeEntry{}

	for _, fileArg := range strings.Split(filesToBeModified, ",") {
		fileContent := []byte("THIS IS MODIFIED CONTENT BY THE CODE...")
		entries = append(entries, &github.TreeEntry{Path: github.String(fileArg), Type: github.String("blob"), Content: github.String(string(fileContent)), Mode: github.String("100644")})
	}
	tree, _, err = client.Git.CreateTree(ctx, owner, repo, *ref.Object.SHA, entries)
	if err != nil {
		logger.ErrorLogger.Printf("ERROR OCCURED IN CREATING A TREE: %v", err)
		return nil, err
	}
	return tree, err
}

// pushCommit pushes the changes in the repo to remote branch
func pushCommit(ref *github.Reference, tree *github.Tree, client *github.Client) (err error) {
	// Get the parent commit to attach the commit to.
	parent, _, err := client.Repositories.GetCommit(ctx, owner, repo, *ref.Object.SHA)
	if err != nil {
		return err
	}

	parent.Commit.SHA = parent.SHA

	// Create the commit using the tree.
	date := time.Now()
	author := &github.CommitAuthor{Date: &date, Name: github.String(authorName), Email: github.String(authorEmail)}
	commit := &github.Commit{Author: author, Message: github.String(commitMessage), Tree: tree, Parents: []*github.Commit{parent.Commit}}
	newCommit, _, err := client.Git.CreateCommit(ctx, owner, repo, commit)
	if err != nil {
		logger.ErrorLogger.Printf("ERROR IN CREATING COMMIT: %v", err)
		return err
	}
	// Attach the commit to the master branch.
	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(ctx, owner, repo, ref, false)
	if err != nil {
		logger.ErrorLogger.Printf("ERROR IN UPDATING FILES ON REMOTE: %v", err)
		return err
	}
	return err
}

// createPullRequest creates the pull request
func createPullRequest(client *github.Client) error {
	pullReq := &github.NewPullRequest{
		Title:               github.String(pullReqTitle),
		Head:                github.String(branch),
		Base:                github.String(baseBranch),
		Body:                github.String(pullReqDescription),
		MaintainerCanModify: github.Bool(true),
	}
	pr, _, err := client.PullRequests.Create(ctx, owner, repo, pullReq)
	if err != nil {
		logger.ErrorLogger.Printf("ERROR IN CREATING PULL REQUEST: %v", err)
		return err
	}
	logger.DebugLogger.Printf("PULL REQUEST CREATED: %s\n", *pr.URL)
	return nil
}
