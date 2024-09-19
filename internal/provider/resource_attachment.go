package provider

import (
    "context"
    "fmt"
    "strings"
    "github.com/hashicorp/go-cty/cty"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAttachment() *schema.Resource {
    resourceAttachmentSchema := attachmentSchema()
    resourceAttachmentSchema[attributeAttachmentFile] = &schema.Schema{
        Description:      descriptionItemAttachmentFile,
        Type:             schema.TypeString,
        Optional:         true,
        ForceNew:         true,
        ValidateDiagFunc: fileHashComputable,
        StateFunc:        fileHash,
    }
    resourceAttachmentSchema[attributeAttachmentContent] = &schema.Schema{
        Description:      descriptionItemAttachmentContent,
        Type:             schema.TypeString,
        Optional:         true,
        ForceNew:         true,
        ValidateDiagFunc: fileHashComputable,
        StateFunc:        fileHash,
    }
    resourceAttachmentSchema[attributeAttachmentItemID] = &schema.Schema{
        Description: descriptionItemIdentifier,
        Type:        schema.TypeString,
        Required:    true,
        ForceNew:    true,
    }

    resourceAttachmentSchema[attributeAttachmentFile].ConflictsWith = []string{attributeAttachmentContent}
    resourceAttachmentSchema[attributeAttachmentContent].ConflictsWith = []string{attributeAttachmentFile}

    return &schema.Resource{
        Description: "Manages an item attachment.",

        CreateContext: attachmentCreate,
        ReadContext:   attachmentRead,
        DeleteContext: attachmentDelete,
        Importer:      importAttachmentResource(),

        Schema: resourceAttachmentSchema,
        CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
            if _, ok := d.GetOk(attributeAttachmentFile); !ok {
                if _, ok := d.GetOk(attributeAttachmentContent); !ok {
                    return fmt.Errorf("either %s or %s must be specified", attributeAttachmentFile, attributeAttachmentContent)
                }
            }
            return nil
        },
    }
}

func importAttachmentResource() *schema.ResourceImporter {
    return &schema.ResourceImporter{
        StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
            split := strings.Split(d.Id(), "/")
            if len(split) != 2 {
                return nil, fmt.Errorf("invalid ID specified, should be in the format <item_id>/<attachment_id>: '%s'", d.Id())
            }
            d.SetId(split[0])
            d.Set(attributeAttachmentItemID, split[1])
            return []*schema.ResourceData{d}, nil
        },
    }
}

func validateOnlyOneOf(attrs ...string) schema.SchemaValidateDiagFunc {
    return func(i interface{}, path cty.Path) diag.Diagnostics {
        var setCount int
        for _, attr := range attrs {
            if v, ok := i.(map[string]interface{})[attr]; ok && v != "" {
                setCount++
            }
        }
        if setCount != 1 {
            return diag.Diagnostics{
                {
                    Severity: diag.Error,
                    Summary:  "Validation Error",
                    Detail:   fmt.Sprintf("Only one of %s must be specified", strings.Join(attrs, ", ")),
                },
            }
        }
        return nil
    }
}

func fileHashComputable(val interface{}, _ cty.Path) diag.Diagnostics {
    _, err := fileSha1Sum(val.(string))
    if err != nil {
        return diag.FromErr(fmt.Errorf("unable to compute hash of file: %w", err))
    }
    return diag.Diagnostics{}
}

func fileHash(val interface{}) string {
    hash, _ := fileSha1Sum(val.(string))
    return hash
}