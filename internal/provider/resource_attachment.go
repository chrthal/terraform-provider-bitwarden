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
	return &schema.Resource{
		Description: "(EXPERIMENTAL) Manages a Vault item's attachment.",

		CreateContext: attachmentCreate,
		ReadContext:   attachmentRead,
		DeleteContext: attachmentDelete,
		Importer:      importAttachmentResource(),

		Schema: map[string]*schema.Schema{
			attributeID: {
				Description: descriptionIdentifier,
				Type:        schema.TypeString,
				Computed:    true,
			},
			attributeAttachmentFile: {
				Description:      descriptionItemAttachmentFile,
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: fileHashComputable,
				StateFunc:        fileHash,
			},
			attributeAttachmentItemID: {
				Description: descriptionItemIdentifier,
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			attributeAttachmentFileName: {
				Description: descriptionItemAttachmentFileName,
				Type:        schema.TypeString,
				Computed:    true,
			},
			attributeAttachmentSize: {
				Description: descriptionItemAttachmentSize,
				Type:        schema.TypeString,
				Computed:    true,
			},
			attributeAttachmentSizeName: {
				Description: descriptionItemAttachmentSizeName,
				Type:        schema.TypeString,
				Computed:    true,
			},
			attributeAttachmentURL: {
				Description: descriptionItemAttachmentURL,
				Type:        schema.TypeString,
				Computed:    true,
			},
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
