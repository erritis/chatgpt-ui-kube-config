package main

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	configs "github.com/erritis/cdk8skit/v3/cdk8skit/configs"
	deployments "github.com/erritis/cdk8skit/v3/cdk8skit/deployments"
	networks "github.com/erritis/cdk8skit/v3/cdk8skit/networks"
	statefulsets "github.com/erritis/cdk8skit/v3/cdk8skit/statefulsets"
	storages "github.com/erritis/cdk8skit/v3/cdk8skit/storages"
	volumes "github.com/erritis/cdk8skit/v3/cdk8skit/volumes"
)

type DbChartProps struct {
	cdk8s.ChartProps
	Environment string
	Network     string
	StorageName string
}

type WsgiServerChartProps struct {
	cdk8s.ChartProps
	Network string
}

type WebServerChartProps struct {
	cdk8s.ChartProps
	Network       string
	ClusterIssuer string
}

type ClientChartProps struct {
	cdk8s.ChartProps
	Network       string
	ClusterIssuer string
}

type NetworkChartProps struct {
	cdk8s.ChartProps
	Network string
}

func NewDbChart(scope constructs.Construct, id string, props *DbChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	if props.Environment == "Production" {
		storages.NewLocalStorage(chart, props.StorageName, &storages.LocalStorageProps{})
		lpvr := volumes.NewLocalVolume(
			chart,
			"persistent-volume",
			jsii.String("/mnt/chatgptdb"),
			&volumes.LocalVolumeProps{
				StorageClassName: jsii.String(props.StorageName),
			},
		)

		statefulsets.NewPostgres(
			chart,
			id,
			&statefulsets.PostgresProps{
				Image:            jsii.String("postgres:12.9"),
				PrefixSecretName: jsii.String("chatgpt-db"),
				DBConfig: &statefulsets.DBConfig{
					Name:     jsii.String("{{ .Values.Db.Name }}"),
					Username: jsii.String("{{ .Values.Db.Username }}"),
					Password: jsii.String("{{ .Values.Db.Password }}"),
				},
				VolumeConfig: &statefulsets.VolumeConfig{
					Volume: &lpvr.Volume,
				},
				Network: jsii.String(props.Network),
			},
		)
	}

	if props.Environment == "Development" {
		statefulsets.NewPostgres(
			chart,
			id,
			&statefulsets.PostgresProps{
				Image:            jsii.String("postgres:12.9"),
				PrefixSecretName: jsii.String("chatgpt-db"),
				DBConfig: &statefulsets.DBConfig{
					Name:     jsii.String("{{ .Values.Db.Name }}"),
					Username: jsii.String("{{ .Values.Db.Username }}"),
					Password: jsii.String("{{ .Values.Db.Password }}"),
				},
				VolumeConfig: &statefulsets.VolumeConfig{
					StorageClassName: jsii.String(props.StorageName),
				},
				Network: jsii.String(props.Network),
			},
		)
	}

	// if props.Environment != "Production" && props.Environment != "Development" {
	// 	statefulsets.NewKubePostgres(
	// 		chart,
	// 		id,
	// 		&statefulsets.KubePostgresProps{
	// 			Image:            jsii.String("postgres:12.9"),
	// 			PrefixSecretName: jsii.String("chatgpt-db"),
	// 			DBConfig: &statefulsets.KubeDBConfig{
	// 				Name:     jsii.String("{{ .Values.Db.Name }}"),
	// 				Username: jsii.String("{{ .Values.Db.Username }}"),
	// 				Password: jsii.String("{{ .Values.Db.Password }}"),
	// 			},
	// 			Network: jsii.String(props.Network),
	// 		},
	// 	)
	// }

	return chart
}

func NewWsgiServerChart(scope constructs.Construct, id string, props *WsgiServerChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	deployments.NewBackend(
		chart,
		id,
		jsii.String("wongsaang/chatgpt-ui-wsgi-server:latest"),
		&deployments.BackendProps{
			PortConfig: &configs.ServicePortConfig{
				ContainerPort: jsii.Number(8000),
			},
			Variables: &map[*string]*string{
				jsii.String("APP_DOMAIN"):                 jsii.String("{{ .Values.WsgiServer.Domain }}"),
				jsii.String("DB_URL"):                     jsii.String("{{ .Values.WsgiServer.DbUrl }}"),
				jsii.String("DJANGO_SUPERUSER_USERNAME"):  jsii.String("{{ .Values.WsgiServer.Django.Superuser.Username }}"),
				jsii.String("DJANGO_SUPERUSER_PASSWORD"):  jsii.String("{{ .Values.WsgiServer.Django.Superuser.Password }}"),
				jsii.String("DJANGO_SUPERUSER_EMAIL"):     jsii.String("{{ .Values.WsgiServer.Django.Superuser.Email }}"),
				jsii.String("SERVER_WORKERS"):             jsii.String("3"),
				jsii.String("WORKER_TIMEOUT"):             jsii.String("180"),
				jsii.String("ACCOUNT_EMAIL_VERIFICATION"): jsii.String("{{ .Values.WsgiServer.AccountEmailVerification }}"),
			},
			Network: jsii.String(props.Network),
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

	deployments.NewFrontend(
		chart,
		id,
		jsii.String("{{ .Values.WebServer.WsgiDomain }}"),
		jsii.String("wongsaang/chatgpt-ui-web-server:latest"),
		&deployments.FrontendProps{
			PortConfig: &configs.ServicePortConfig{
				ContainerPort: jsii.Number(80),
			},
			Variables: &map[*string]*string{
				jsii.String("BACKEND_URL"): jsii.String("{{ .Values.WebServer.BackendUrl }}"),
			},
			ClusterIssuer: &props.ClusterIssuer,
			Network:       jsii.String(props.Network),
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

	deployments.NewFrontend(
		chart,
		id,
		jsii.String("{{ .Values.Client.Domain }}"),
		jsii.String("wongsaang/chatgpt-ui-client:latest"),
		&deployments.FrontendProps{
			PortConfig: &configs.ServicePortConfig{
				ContainerPort: jsii.Number(80),
			},
			Variables: &map[*string]*string{
				jsii.String("SERVER_DOMAIN"):          jsii.String("{{ .Values.Client.ServerUrl }}"),
				jsii.String("NUXT_PUBLIC_APP_NAME"):   jsii.String("{{ .Values.Client.NuxtPublicAppName }}"),
				jsii.String("NUXT_PUBLIC_TYPEWRITER"): jsii.String("false"),
			},
			ClusterIssuer: &props.ClusterIssuer,
			Network:       jsii.String(props.Network),
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

	networks.NewNetworkPolicy(chart, id, props.Network)

	return chart
}

func main() {

	config, err := LoadConfig()

	if err != nil {
		fmt.Println(err)
	}

	app := cdk8s.NewApp(&cdk8s.AppProps{
		Outdir:              jsii.String(config.Outdir),
		OutputFileExtension: jsii.String(".yaml"),
		YamlOutputType:      cdk8s.YamlOutputType_FILE_PER_CHART,
	})

	cprops := cdk8s.ChartProps{
		DisableResourceNameHashes: jsii.Bool(true),
		Namespace:                 jsii.String("chatgpt-ui"),
	}

	NewDbChart(app, "chatgpt-ui-db", &DbChartProps{
		ChartProps:  cprops,
		Network:     "io.network/chatgpt-ui-network",
		Environment: config.Environment,
		StorageName: config.StorageName,
	})
	NewWsgiServerChart(app, "chatgpt-ui-wsgi-server", &WsgiServerChartProps{
		ChartProps: cprops,
		Network:    "io.network/chatgpt-ui-network",
	})
	NewWebServerChart(app, "chatgpt-ui-web-server", &WebServerChartProps{
		ChartProps:    cprops,
		Network:       "io.network/chatgpt-ui-network",
		ClusterIssuer: config.ClusterIssuer,
	})
	NewClientChart(app, "chatgpt-ui-client", &ClientChartProps{
		ChartProps:    cprops,
		Network:       "io.network/chatgpt-ui-network",
		ClusterIssuer: config.ClusterIssuer,
	})

	NewNetworkChart(app, "chatgpt-ui-network", &NetworkChartProps{
		ChartProps: cprops,
		Network:    "io.network/chatgpt-ui-network",
	})

	app.Synth()
}
