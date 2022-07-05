package gconsumer

import (
	"tirelease/commons/git"
	"tirelease/internal/entity"
	"tirelease/internal/repository"

	"github.com/google/go-github/v41/github"
)

type PrReleaseNoteConsumer struct {
}

const NONE_RELEASE_NOTE_LABEL = "release-note-none"

// ref: [release\-note](https://book.prow.tidb.io/#/plugins/release-note)
func (consumer PrReleaseNoteConsumer) Consume(event github.PullRequestEvent) error {
	prEntity := entity.ComposePullRequestFromV3(event.PullRequest)
	releaseNote, err := git.ParseReleaseNote(prEntity.Body)

	if err != nil || !releaseNote.IsReleaseNoteConfirmed {
		labels := prEntity.Labels
		for _, label := range *labels {
			if *label.Name == NONE_RELEASE_NOTE_LABEL {
				releaseNote.IsReleaseNoteConfirmed = true
				releaseNote.ReleaseNote = "None"
			}
		}
	}

	if !releaseNote.IsReleaseNoteConfirmed {
		return nil
	}

	prEntity.IsReleaseNoteConfirmed = releaseNote.IsReleaseNoteConfirmed
	prEntity.ReleaseNote = releaseNote.ReleaseNote
	err = repository.CreateOrUpdatePullRequest(prEntity)

	return err
}

func (consumer PrReleaseNoteConsumer) Validate(event github.PullRequestEvent) bool {
	return true
}
