package provider

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rfleming71/terraform-provider-mayan-edms/client"
)

func resourceWorkflowTemplateTransitionTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceWorkflowTemplateTransitionTriggerCreate,
		Read:   resourceWorkflowTemplateTransitionTriggerRead,
		Delete: resourceWorkflowTemplateTransitionTriggerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceWorkflowTemplateTransitionTriggerImport,
		},

		Schema: map[string]*schema.Schema{
			"workflow_template": {
				Description: "ID of the workflow template.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"transition": {
				Description: "ID of the workflow template transition (composite ID format: workflow_template_id-transition_id).",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"event_type_id": {
				Description: "Event type that triggers this transition (e.g., workflow_instance.created, tag.attach).",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"event_type_label": {
				Description: "Human-readable label for the event type.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceWorkflowTemplateTransitionTriggerCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)

	workflowTemplateId := d.Get("workflow_template").(int)
	transitionCompositeId := d.Get("transition").(string)
	eventTypeId := d.Get("event_type_id").(string)

	_, transitionId, err := breakCompositeId(transitionCompositeId)
	if err != nil {
		return fmt.Errorf("invalid transition ID format: %s", transitionCompositeId)
	}

	trigger, err := c.CreateWorkflowTransitionTrigger(workflowTemplateId, transitionId, eventTypeId)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v-%v-%v", workflowTemplateId, transitionId, trigger.ID))

	return resourceWorkflowTemplateTransitionTriggerRead(d, m)
}

func resourceWorkflowTemplateTransitionTriggerRead(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)

	workflowTemplateId, transitionId, triggerId, err := getTriggerIdInformation(d)
	if err != nil {
		return err
	}

	trigger, err := c.GetWorkflowTransitionTrigger(workflowTemplateId, transitionId, triggerId)
	if err != nil {
		return err
	}

	return triggerToData(workflowTemplateId, transitionId, trigger, d)
}

func resourceWorkflowTemplateTransitionTriggerDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(client.MayanEdmsClient)

	workflowTemplateId, transitionId, triggerId, err := getTriggerIdInformation(d)
	if err != nil {
		return err
	}

	err = c.DeleteWorkflowTransitionTrigger(workflowTemplateId, transitionId, triggerId)
	if err == nil {
		d.SetId("")
	}

	return err
}

func resourceWorkflowTemplateTransitionTriggerImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(client.MayanEdmsClient)

	workflowTemplateId, transitionId, triggerId, err := getTriggerIdInformation(d)
	rd := []*schema.ResourceData{d}
	if err != nil {
		return rd, err
	}

	trigger, err := c.GetWorkflowTransitionTrigger(workflowTemplateId, transitionId, triggerId)
	if err != nil {
		return rd, err
	}

	err = triggerToData(workflowTemplateId, transitionId, trigger, d)
	return rd, err
}

func triggerToData(workflowTemplateId int, transitionId int, trigger *client.WorkflowTemplateTransitionTrigger, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%v-%v-%v", workflowTemplateId, transitionId, trigger.ID))

	if err := d.Set("workflow_template", workflowTemplateId); err != nil {
		return err
	}

	transitionCompositeId := fmt.Sprintf("%v-%v", workflowTemplateId, transitionId)
	if err := d.Set("transition", transitionCompositeId); err != nil {
		return err
	}

	if err := d.Set("event_type_id", trigger.EventType.ID); err != nil {
		return err
	}

	if err := d.Set("event_type_label", trigger.EventType.Label); err != nil {
		return err
	}

	return nil
}

func getTriggerIdInformation(d *schema.ResourceData) (int, int, int, error) {
	return breakTripleCompositeId(d.Id())
}

func breakTripleCompositeId(id string) (int, int, int, error) {
	ids := strings.Split(id, "-")
	if len(ids) != 3 {
		return 0, 0, 0, fmt.Errorf("expected id format: workflow_template_id-transition_id-trigger_id, got: %s", id)
	}

	part1, err := strconv.Atoi(ids[0])
	if err != nil {
		return 0, 0, 0, err
	}

	part2, err := strconv.Atoi(ids[1])
	if err != nil {
		return 0, 0, 0, err
	}

	part3, err := strconv.Atoi(ids[2])
	if err != nil {
		return 0, 0, 0, err
	}

	return part1, part2, part3, nil
}
