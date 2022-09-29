/*
package configuration, keep  environment, configuration parameter here in this package.
*/

package configuration

import "github.com/madhurikadam/app-transcation/pkg/database/postgres"

type Config struct {
	postgres.Config

	HTTPPort       int    `envconfig:"HTTP_PORT" default:"8080"`
	AllowedOrigins string `envconfig:"ALLOWED_ORIGINS" default:"*"`
}
