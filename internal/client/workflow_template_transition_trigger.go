package client

import (
	"fmt"
	"net/http"
)

type WorkflowTemplateTransitionTrigger struct {
	ID        int       `json:"id"`
	EventType EventType `json:"event_type"`
}

type EventType struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Name  string `json:"name"`
}

type workflowTemplateTransitionTriggerRequest struct {
	EventTypeID string `json:"event_type_id"`
}

func (c *Client) GetWorkflowTransitionTrigger(workflowTemplateId int, transitionId int, triggerId int) (*WorkflowTemplateTransitionTrigger, error) {
	var result WorkflowTemplateTransitionTrigger
	err := c.performRequest(
		fmt.Sprintf("workflow_templates/%v/transitions/%v/triggers/%v/", workflowTemplateId, transitionId, triggerId),
		http.MethodGet,
		nil,
		&result,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) CreateWorkflowTransitionTrigger(workflowTemplateId int, transitionId int, eventTypeId string) (*WorkflowTemplateTransitionTrigger, error) {
	var newTrigger WorkflowTemplateTransitionTrigger
	request := workflowTemplateTransitionTriggerRequest{
		EventTypeID: eventTypeId,
	}
	err := c.performRequest(
		fmt.Sprintf("workflow_templates/%v/transitions/%v/triggers/", workflowTemplateId, transitionId),
		http.MethodPost,
		&request,
		&newTrigger,
	)
	if err != nil {
		return &WorkflowTemplateTransitionTrigger{}, err
	}

	return &newTrigger, nil
}

func (c *Client) DeleteWorkflowTransitionTrigger(workflowTemplateId int, transitionId int, triggerId int) error {
	err := c.performRequest(
		fmt.Sprintf("workflow_templates/%v/transitions/%v/triggers/%v/", workflowTemplateId, transitionId, triggerId),
		http.MethodDelete,
		nil,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}
