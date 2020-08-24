package restapi

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceApiRawObject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApiRawObjectRead,

		Schema: map[string]*schema.Schema{
			"path": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The API path on top of the base URL set in the provider that represents objects of this type on the API server.",
				Required:    true,
			},
			"query_string": &schema.Schema{
				Type:        schema.TypeString,
				Description: "An optional query string to send when performing the search.",
				Optional:    true,
			},
			"debug": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Whether to emit verbose debug output while working with the API object on the server.",
				Optional:    true,
			},
			"api_response": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The raw body of the HTTP response from the last read of the object.",
				Computed:    true,
			},
		}, /* End schema */

	}
}

func dataSourceApiRawObjectRead(d *schema.ResourceData, meta interface{}) error {
	path := d.Get("path").(string)
	query_string := d.Get("query_string").(string)
	debug := d.Get("debug").(bool)
	client := meta.(*api_client)
	if debug {
		log.Printf("datasource_api_raw_object.go: Data routine called.")
		log.Printf("datasource_api_raw_object.go:\npath: %s\nquery_string: %s", path, query_string)
	}

	opts := &apiObjectOpts{
		path:  path,
		debug: debug,
	}

	obj, err := NewAPIObject(client, opts)
	if err != nil {
		return err
	}

	/* Back to terraform-specific stuff. Create an api_object with the ID and refresh it object */
	if debug {
		log.Printf("datasource_api_raw_object.go: Attempting to construct api_object to refresh data")
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	err = obj.read_raw_object()
	if err == nil {
		set_resource_state(obj, d)
	}
	return err
}
