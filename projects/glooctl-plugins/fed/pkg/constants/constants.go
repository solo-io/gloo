package constants

import (
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

var (
	SubcommandError = eris.New("please select a subcommand")
)

var (
	VIRTUAL_SERVICE_COMMAND = cobra.Command{
		Use:     "virtualservice",
		Aliases: []string{"vs", "virtualservices"},
	}

	ROUTE_TABLE_COMMAND = cobra.Command{
		Use:     "routetable",
		Aliases: []string{"rt", "routetables"},
	}

	UPSTREAM_COMMAND = cobra.Command{
		Use:     "upstream",
		Aliases: []string{"u", "us", "upstreams"},
	}

	UPSTREAM_GROUP_COMMAND = cobra.Command{
		Use:     "upstreamgroup",
		Aliases: []string{"ug", "ugs", "upstreamgroups"},
	}

	PROXY_COMMAND = cobra.Command{
		Use:     "proxy",
		Aliases: []string{"p", "proxies"},
	}

	SETTINGS_COMMAND = cobra.Command{
		Use:     "setting",
		Aliases: []string{"s", "settings"},
	}

	GATEWAY_COMMAND = cobra.Command{
		Use:     "gateway",
		Aliases: []string{"g", "gateways"},
	}

	MATCHABLE_HTTP_GATEWAY_COMMAND = cobra.Command{
		Use:     "matchablehttpgateway",
		Aliases: []string{"hgw", "matchablehttpgateways"},
	}

	AUTH_CONFIG_COMMAND = cobra.Command{
		Use:     "authconfig",
		Aliases: []string{"ac", "authconfig"},
	}

	RATE_LIMIT_CONFIG_COMMAND = cobra.Command{
		Use:     "ratelimitconfig",
		Aliases: []string{"rlc", "ratelimitconfigs"},
	}

	VERSION_COMMAND = cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Print current version",
		Long:    "Get the version of glooctl-fed",
	}
)
