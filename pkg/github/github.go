package github

import (
	"encoding/json"
	"github_statistics/internal/log"
	"github_statistics/pkg/model"
	"time"
)

type Event interface {
	Dump() []byte
}

const (
	EventIssueComment             string = "issue_comment"
	EventIssues                          = "issues"
	EventPullRequest                     = "pull_request"
	EventPullRequestReview               = "pull_request_review"
	EventPullRequestReviewComment        = "pull_request_review_comment"
	EventStar                            = "star"
)

func HandleEvent(ty string, data []byte) error {
	var (
		event Event
		err   error
		dev   model.Developer
	)

	switch ty {
	case EventIssueComment:
		var info *model.IssueComment
		info, err = bindIssueComment(data)
		event = info
		dev = model.Developer{
			Id: info.Repository.Owner.ID,
			Name: info.Repository.Owner.Login,
			CreatedAt: time.Now(),
			Project: info.Repository.Name,
			EventType: EventIssueComment,
			Action: info.Action,
		}
	case EventIssues:
		var info *model.Issues
		info, err = bindIssues(data)
		event = info
		dev = model.Developer{
			Id: info.Repository.Owner.ID,
			Name: info.Repository.Owner.Login,
			CreatedAt: time.Now(),
			Project: info.Repository.Name,
			EventType: EventIssues,
			Action: info.Action,
		}
	case EventPullRequest:
		var info *model.PullRequest
		info, err = bindPullRequest(data)
		event = info
		dev = model.Developer{
			Id: info.Repository.Owner.ID,
			Name: info.Repository.Owner.Login,
			CreatedAt: time.Now(),
			Project: info.Repository.Name,
			EventType: EventPullRequest,
			Action: info.Action,
		}
	case EventPullRequestReview:
		var info *model.PullRequestReview
		info, err = bindPullRequestReview(data)
		event = info
		dev = model.Developer{
			Id: info.Repository.Owner.ID,
			Name: info.Repository.Owner.Login,
			CreatedAt: time.Now(),
			Project: info.Repository.Name,
			EventType: EventPullRequestReview,
			Action: info.Action,
		}
	case EventPullRequestReviewComment:
		var info *model.PullRequestReviewComment
		info, err = bindPullRequestReviewComment(data)
		event = info
		dev = model.Developer{
			Id: info.Repository.Owner.ID,
			Name: info.Repository.Owner.Login,
			CreatedAt: time.Now(),
			Project: info.Repository.Name,
			EventType: EventPullRequestReviewComment,
			Action: info.Action,
		}
	case EventStar:
		var info *model.Star
		info, err = bindStar(data)
		event = info
		if info.Action == "deleted" {

		} else {
			dev = model.Developer{
				Id: info.Repository.Owner.ID,
				Name: info.Repository.Owner.Login,
				CreatedAt: time.Now(),
				Project: info.Repository.Name,
				EventType: EventStar,
				Action: info.Action,
			}
		}
	default:
		log.Infof("Unhandled event:%v", data)
	}
	if err != nil {
		log.Errorf("Categorize type:%s found error:%v", ty, err)
		return err
	}

	err = dev.Insert()
	if err != nil {
		log.Errorf("insert developer:%v found error:%v", dev, err)
		return err
	}
	log.Info("[Event] will handle:", string(event.Dump()))
	return nil
}

func bindIssueComment(data []byte) (*model.IssueComment, error) {
	var info model.IssueComment
	err := json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func bindIssues(data []byte) (*model.Issues, error) {
	var info model.Issues
	err := json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func bindPullRequest(data []byte) (*model.PullRequest, error) {
	var info model.PullRequest
	err := json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func bindPullRequestReview(data []byte) (*model.PullRequestReview, error) {
	var info model.PullRequestReview
	err := json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func bindPullRequestReviewComment(data []byte) (*model.PullRequestReviewComment, error) {
	var info model.PullRequestReviewComment
	err := json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func bindStar(data []byte) (*model.Star, error) {
	var info model.Star
	err := json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}
