package aws

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/phzietsman/terraform-provider-awsx/internal/provider/internal/keyvaluetags"
)

func resourceControlTowerAccountVending() *schema.Resource {
	return &schema.Resource{
		Create: resourceControlTowerAccountVendingCreate,
		Read:   resourceControlTowerAccountVendingRead,
		Update: resourceControlTowerAccountVendingUpdate,
		Delete: resourceControlTowerAccountVendingDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"record_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"artefact_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parameters": {
				Type:     schema.TypeMap,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
		},
	}
}

func resourceControlTowerAccountVendingCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).scconn
	input := servicecatalog.ProvisionProductInput{

		AcceptLanguage: aws.String("en"),
		ProvisionToken: aws.String(resource.UniqueId()),

		NotificationArns: []*string{},

		ProductId:              aws.String(d.Get("product_id").(string)),
		ProvisioningArtifactId: aws.String(d.Get("artefact_id").(string)),
		ProvisionedProductName: aws.String(d.Get("name").(string)),

		Tags: keyvaluetags.New(d.Get("tags").(map[string]interface{})).IgnoreAws().ServicecatalogTags(),
	}

	value := keyvaluetags.New(d.Get("parameters").(map[string]interface{}))
	provisioningParameter := make([]*servicecatalog.ProvisioningParameter, len(value))
	index := 0
	for k, v := range value {
		log.Printf("[DEBUG] ProvisioningParameter[%s] Key:%s Value:%s", string(index), k, *v)
		strv := *v
		strk := k
		provisioningParameter[index] = &servicecatalog.ProvisioningParameter{
			Key:   &strk,
			Value: &strv,
		}

		index++
	}

	input.ProvisioningParameters = provisioningParameter

	log.Printf("[DEBUG] Creating Service Catalog Portfolio: %#v", input)
	resp, err := conn.ProvisionProduct(&input)
	if err != nil {
		return fmt.Errorf("Creating Service Catalog Portfolio failed: %s", err.Error())
	}

	recordDetail := *resp.RecordDetail

	d.SetId(*recordDetail.ProvisionedProductId)
	log.Printf("[INFO] Provisioned Product Id: %s", d.Id())

	d.Set("record_id", *recordDetail.RecordId)

	// Wait for the Provisioned Product to become available
	log.Printf("[DEBUG] Waiting for Provisioned Product (%s) to become available", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"CREATED", "IN_PROGRESS"},
		Target:  []string{"SUCCEEDED"},
		Refresh: ProvisionedProductStateRefreshFunc(conn, *recordDetail.RecordId),
		Timeout: 60 * time.Minute,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Provisioned Product (%s) to become available: %s", d.Id(), err)
	}

	return resourceControlTowerAccountVendingRead(d, meta)
}

func resourceControlTowerAccountVendingRead(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*AWSClient).scconn

	// Get Provisioned Product
	inputProvisionedProduct := servicecatalog.DescribeProvisionedProductInput{
		AcceptLanguage: aws.String("en"),
		Id:             aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Reading Service Catalog Provisoned Product: %#v", inputProvisionedProduct)

	respProvisionedProduct, err := conn.DescribeProvisionedProduct(&inputProvisionedProduct)
	if err != nil {
		if scErr, ok := err.(awserr.Error); ok && scErr.Code() == "ResourceNotFoundException" {
			log.Printf("[WARN] Service Catalog Provisioned Product %q not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Reading ServiceCatalog Provisioned Product '%s' failed: %s", *inputProvisionedProduct.Id, err.Error())
	}

	provisionedProductDetail := respProvisionedProduct.ProvisionedProductDetail

	if err := d.Set("created_time", provisionedProductDetail.CreatedTime.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Error setting created_time: %s", err)
	}
	d.Set("arn", provisionedProductDetail.Arn)

	// Get Record with Account Number
	inputRecord := servicecatalog.DescribeRecordInput{
		AcceptLanguage: aws.String("en"),
		Id:             aws.String(d.Get("record_id").(string)),
	}

	respRecord, err := conn.DescribeRecord(&inputRecord)
	if err != nil {
		if scErr, ok := err.(awserr.Error); ok && scErr.Code() == "ResourceNotFoundException" {
			log.Printf("[WARN] Service Catalog Resource %q not found, remove account_id", *inputProvisionedProduct.Id)
			d.Set("account_id", "")
			return nil
		}
		return fmt.Errorf("Reading ServiceCatalog Record '%s' failed: %s", *inputProvisionedProduct.Id, err.Error())
	}

	recordOutputs := respRecord.RecordOutputs

	for _, recordOutput := range recordOutputs {
		if *recordOutput.OutputKey == "AccountId" {
			d.Set("account_id", recordOutput.OutputValue)
		}
	}

	return nil
}

func resourceControlTowerAccountVendingUpdate(d *schema.ResourceData, meta interface{}) error {
	// Do not support updating of vended accounts. I dont fully understand what is
	// supported from an updating point of view.
	return nil
}

func resourceControlTowerAccountVendingDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).scconn
	input := servicecatalog.TerminateProvisionedProductInput{
		ProvisionedProductId: aws.String(d.Id()),
		TerminateToken:       aws.String(resource.UniqueId()),
	}

	log.Printf("[DEBUG] Delete Vended Account (Service Catalog Provisioned Product): %#v", input)
	resp, err := conn.TerminateProvisionedProduct(&input)
	if err != nil {
		return fmt.Errorf("Deleting Vended Account (Service Catalog Provisioned Product) '%s' failed: %s", *input.ProvisionedProductId, err.Error())
	}

	// Wait for the Provisioned Product to become available
	log.Printf("[DEBUG] Waiting for Provisioned Product (%s) to become terminated", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"CREATED", "IN_PROGRESS"},
		Target:  []string{"SUCCEEDED"},
		Refresh: ProvisionedProductStateRefreshFunc(conn, *resp.RecordDetail.RecordId),
		Timeout: 60 * time.Minute,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Provisioned Product (%s) to become available: %s", d.Id(), err)
	}
	return nil
}

// ProvisionedProductStateRefreshFunc returns a resource.StateRefreshFunc
// that is used to watch a Provisioned Product.
func ProvisionedProductStateRefreshFunc(conn *servicecatalog.ServiceCatalog, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		opts := &servicecatalog.DescribeRecordInput{
			Id: aws.String(id),
		}
		resp, err := conn.DescribeRecord(opts)
		if err != nil {
			if isAWSErr(err, "ResourceNotFoundException", "") {
				resp = nil
			} else {
				log.Printf("Error on ProvisionedProductStateRefreshFunc: %s", err)
				return nil, "", err
			}
		}

		if resp == nil {
			// Sometimes AWS just has consistency issues and doesn't see
			// the resource yet. Return an empty state.
			return nil, "", nil
		}

		rec := resp.RecordDetail
		return rec, *rec.Status, nil
	}
}
