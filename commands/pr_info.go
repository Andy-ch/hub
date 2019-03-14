package commands

import (
        "fmt"

        "github.com/github/hub/github"
        "github.com/github/hub/ui"
        "github.com/github/hub/utils"
)

func prInfo(command *Command, args *Args) {
        localRepo, err := github.LocalRepo()
        utils.Check(err)

        project, err := localRepo.MainProject()
        utils.Check(err)

        gh := github.NewClient(project.Host)

        args.NoForward()
        words := args.Words()
        if len(words) == 0 {
                utils.Check(fmt.Errorf("Error: no pull request number given"))
        }

        prNumberString := words[0]
        pr, err := gh.PullRequest(project, prNumberString)
        utils.Check(err)
        flagPullRequestFormat := args.Flag.Value("--format")
        if !args.Flag.HasReceived("--format") {
                flagPullRequestFormat = "%pC%>(8)%i%Creset  %t%  l%nRequested reviewers: %rs%n%n%b%n"
        }
        colorize := colorizeOutput(args.Flag.HasReceived("--color"), args.Flag.Value("--color"))
        ui.Print(formatPullRequest(*pr, flagPullRequestFormat, colorize))
        comments, err := gh.FetchPRComments(project, prNumberString)
        utils.Check(err)
        reviewsToComments := make(map[int]int)
        for i, comment := range comments {
                reviewsToComments[comment.ReviewId] = i
        }
        reviews, err := gh.FetchReviews(project, prNumberString)
        utils.Check(err)
        for _, review := range reviews {
                ui.Printf("%s %s on %s\n", review.User.Login, review.State, review.CreatedAt)
                if review.State != "APPROVED" {
                        comment := comments[reviewsToComments[review.Id]]
                        ui.Printf("%s\n\n%s\n", comment.DiffHunk, comment.Body)
                }
                ui.Printf("\n")
        }
}
