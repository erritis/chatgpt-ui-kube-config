package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2"
	cdk8skit "github.com/erritis/cdk8s-kit"
)

type DbChartProps struct {
	cdk8s.ChartProps
	network string
}

type WsgiServerChartProps struct {
	cdk8s.ChartProps
	network string
}

type WebServerChartProps struct {
	cdk8s.ChartProps
	network string
}

type ClientChartProps struct {
	cdk8s.ChartProps
	network string
}

type NetworkChartProps struct {
	cdk8s.ChartProps
	network string
}

func NewDbChart(scope constructs.Construct, id string, props *DbChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	storage := cdk8skit.NewLocalStorage(chart, "chatgpt-ui-local-storage")

	dbData := cdk8skit.NewLocalVolume(
		chart,
		*storage.Name(),
		"persistent-volume",
		cdk8s.Size_Gibibytes(jsii.Number(0.1)),
		"/mnt/chatgptdb",
		&[]string{"master-node"},
	)

	db := cdk8skit.NewSecretVolume(
		chart, "db-secret",
		jsii.String("chatgpt-db"),
		jsii.String("{{ .Values.db_name }}"),
	)

	dbUser := cdk8skit.NewSecretVolume(
		chart, "user-secret",
		jsii.String("chatgpt-db-user"),
		jsii.String("{{ .Values.db_username }}"),
	)

	dbPasswd := cdk8skit.NewSecretVolume(
		chart, "passwd-secret",
		jsii.String("chatgpt-db-passwd"),
		jsii.String("{{ .Values.db_password }}"),
	)

	cdk8skit.NewStatefulSet(
		chart,
		id,
		jsii.String("postgres:12.9"),
		jsii.Number(5432),
		jsii.Number(5432),
		props.network,
		&map[*string]*string{
			jsii.String("POSTGRES_DB_FILE"):       jsii.String("/run/secrets/chatgpt-db/chatgpt-db"),
			jsii.String("POSTGRES_USER_FILE"):     jsii.String("/run/secrets/chatgpt-db-user/chatgpt-db-user"),
			jsii.String("POSTGRES_PASSWORD_FILE"): jsii.String("/run/secrets/chatgpt-db-passwd/chatgpt-db-passwd"),
		},
		&map[*string]*cdk8splus26.Volume{
			jsii.String("/var/lib/postgresql/data"):       &dbData.Volume,
			jsii.String("/run/secrets/chatgpt-db"):        &db,
			jsii.String("/run/secrets/chatgpt-db-user"):   &dbUser,
			jsii.String("/run/secrets/chatgpt-db-passwd"): &dbPasswd,
		},
	)

	return chart
}

func NewWsgiServerChart(scope constructs.Construct, id string, props *WsgiServerChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	cdk8skit.NewBackend(
		chart,
		id,
		jsii.String("wongsaang/chatgpt-ui-wsgi-server:latest"),
		jsii.Number(80),
		jsii.Number(8000),
		props.network,
		&map[*string]*string{
			jsii.String("APP_DOMAIN"):                 jsii.String("{{ .Values.django_domain }}"),
			jsii.String("DB_URL"):                     jsii.String("{{ .Values.db_url }}"),
			jsii.String("DJANGO_SUPERUSER_USERNAME"):  jsii.String("{{ .Values.django_superuser_username }}"),
			jsii.String("DJANGO_SUPERUSER_PASSWORD"):  jsii.String("{{ .Values.django_superuser_password }}"),
			jsii.String("DJANGO_SUPERUSER_EMAIL"):     jsii.String("{{ .Values.django_superuser_email }}"),
			jsii.String("SERVER_WORKERS"):             jsii.String("3"),
			jsii.String("WORKER_TIMEOUT"):             jsii.String("180"),
			jsii.String("ACCOUNT_EMAIL_VERIFICATION"): jsii.String("{{ .Values.account_email_verification }}"),
		},
	)

	return chart
}

func NewWebServerChart(scope constructs.Construct, id string, props *WebServerChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	cdk8skit.NewFrontend(
		chart,
		id,
		jsii.String("{{ .Values.django_domain }}"),
		jsii.String("letsencrypt-prod"),
		jsii.String("wongsaang/chatgpt-ui-web-server:v2.5.2"),
		jsii.Number(80),
		jsii.Number(80),
		props.network,
		&map[*string]*string{
			jsii.String("BACKEND_URL"): jsii.String("{{ .Values.backend_url }}"),
		},
	)

	return chart
}

func NewClientChart(scope constructs.Construct, id string, props *ClientChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	cdk8skit.NewFrontend(
		chart,
		id,
		jsii.String("{{ .Values.client_domain }}"),
		jsii.String("letsencrypt-prod"),
		jsii.String("wongsaang/chatgpt-ui-client:latest"),
		jsii.Number(80),
		jsii.Number(80),
		props.network,
		&map[*string]*string{
			jsii.String("SERVER_DOMAIN"):          jsii.String("{{ .Values.server_domain }}"),
			jsii.String("NUXT_PUBLIC_APP_NAME"):   jsii.String("{{ .Values.nuxt_public_app_name }}"),
			jsii.String("NUXT_PUBLIC_TYPEWRITER"): jsii.String("false"),
		},
	)

	return chart
}

func NewNetworkChart(scope constructs.Construct, id string, props *NetworkChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	cdk8skit.NewNetworkPolicy(chart, id, props.network)

	return chart
}

func main() {
	app := cdk8s.NewApp(&cdk8s.AppProps{
		Outdir:              jsii.String("../.helm/templates"),
		OutputFileExtension: jsii.String(".yaml"),
		YamlOutputType:      cdk8s.YamlOutputType_FILE_PER_CHART,
	})

	cprops := cdk8s.ChartProps{
		DisableResourceNameHashes: jsii.Bool(true),
		Namespace:                 jsii.String("chatgpt-ui"),
	}

	NewDbChart(app, "chatgpt-ui-db", &DbChartProps{
		ChartProps: cprops,
		network:    "io.network/chatgpt-ui-network",
	})
	NewWsgiServerChart(app, "chatgpt-ui-wsgi-server", &WsgiServerChartProps{
		ChartProps: cprops,
		network:    "io.network/chatgpt-ui-network",
	})
	NewWebServerChart(app, "chatgpt-ui-web-server", &WebServerChartProps{
		ChartProps: cprops,
		network:    "io.network/chatgpt-ui-network",
	})
	NewClientChart(app, "chatgpt-ui-client", &ClientChartProps{
		ChartProps: cprops,
		network:    "io.network/chatgpt-ui-network",
	})

	NewNetworkChart(app, "chatgpt-ui-network", &NetworkChartProps{
		ChartProps: cprops,
		network:    "io.network/chatgpt-ui-network",
	})

	app.Synth()
}
