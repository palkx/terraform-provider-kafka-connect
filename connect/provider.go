package connect

import (
	"context"
	"crypto/tls"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	kc "github.com/palkx/go-kafka-connect/v3/lib/connectors"
)

func Provider() *schema.Provider {
	log.Printf("[INFO] Creating Provider")
	provider := schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_URL", ""),
			},
			"headers": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				// No DefaultFunc here to read from the env on account of this issue:
				// https://github.com/hashicorp/terraform-plugin-sdk/issues/142
			},
			"basic_auth_username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_BASIC_AUTH_USERNAME", ""),
			},
			"basic_auth_password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_BASIC_AUTH_PASSWORD", ""),
			},
			"tls_skip_insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_TLS_IS_INSECURE", false),
			},
			"tls_root_ca_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_SSL_ROOT_CA_FILE_PATH", ""),
			},
			"tls_auth_crt_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_TLS_AUTH_CRT_PATH", ""),
			},
			"tls_auth_key_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_TLS_AUTH_KEY_PATH", ""),
			},
		},
		ConfigureContextFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"kafka-connect_connector": kafkaConnectorResource(),
		},
	}
	log.Printf("[INFO] Created provider: %v", provider)
	return &provider
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log.Printf("[INFO] Initializing KafkaConnect client")
	addr := d.Get("url").(string)
	c := kc.NewClient(addr)
	user := d.Get("basic_auth_username").(string)
	pass := d.Get("basic_auth_password").(string)
	if user != "" && pass != "" {
		c.SetBasicAuth(user, pass)
	}

	ca_path := d.Get("tls_root_ca_path").(string)
	crt_path := d.Get("tls_auth_crt_path").(string)
	key_path := d.Get("tls_auth_key_path").(string)
	is_tls_insecure := d.Get("tls_skip_insecure").(bool)

	if ca_path != "" {
		c.SetRootCertificate(ca_path)
		log.Printf("[INFO] CA : %s\n", ca_path)
	}

	if crt_path != "" && key_path != "" {
		cert, err := tls.LoadX509KeyPair(crt_path, key_path)
		if err != nil {
			log.Fatalf("[ERROR] client: loadkeys: %s", err)
		} else {
			if is_tls_insecure {
				log.Printf("[WARN] SSl connection is insecure : %t", is_tls_insecure)
				c.SetInsecureSSL()
			}
			log.Printf("[INFO] Cert : %s\nKey: %s", crt_path, key_path)
			c.SetClientCertificates(cert)
		}
	}

	headers := d.Get("headers").(map[string]interface{})
	if headers != nil {
		for k, v := range headers {
			c.SetHeader(k, v.(string))
		}
	}

	return c, nil
}
