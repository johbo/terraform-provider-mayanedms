# Trigger transition when workflow instance is created
resource "mayanedms_workflow_template_transition_trigger" "on_workflow_created" {
  workflow_template = mayanedms_workflow_template.example.id
  transition        = mayanedms_workflow_template_transition.initial.id
  event_type_id     = "workflow_instance.created"
}

# Trigger transition when a tag is attached
resource "mayanedms_workflow_template_transition_trigger" "on_tag_attach" {
  workflow_template = mayanedms_workflow_template.example.id
  transition        = mayanedms_workflow_template_transition.tag_based.id
  event_type_id     = "tag.attach"
}
