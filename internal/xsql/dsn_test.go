package xsql

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ydb-platform/ydb-go-sdk/v3/config"
	"github.com/ydb-platform/ydb-go-sdk/v3/internal/xsql/bind"
)

func TestParse(t *testing.T) {
	newConnector := func(opts ...ConnectorOption) *Connector {
		c := &Connector{}
		for _, opt := range opts {
			if err := opt(c); err != nil {
				t.Error(err)
			}
		}
		return c
	}
	compareConfigs := func(t *testing.T, lhs, rhs config.Config) {
		require.Equal(t, lhs.Secure(), rhs.Secure())
		require.Equal(t, lhs.Endpoint(), rhs.Endpoint())
		require.Equal(t, lhs.Database(), rhs.Database())
	}
	for _, tt := range []struct {
		dsn              string
		expOpts          []config.Option
		expConnectorOpts []ConnectorOption
		expErr           error
	}{
		{
			dsn: "grpc://localhost:2135/local",
			expOpts: []config.Option{
				config.WithSecure(false),
				config.WithEndpoint("localhost:2135"),
				config.WithDatabase("/local"),
			},
			expConnectorOpts: nil,
			expErr:           nil,
		},
		{
			dsn: "grpcs://localhost:2135/local/db",
			expOpts: []config.Option{
				config.WithSecure(true),
				config.WithEndpoint("localhost:2135"),
				config.WithDatabase("/local/db"),
			},
			expConnectorOpts: nil,
			expErr:           nil,
		},
		{
			dsn: "grpc://localhost:2135/local?query_mode=scripting",
			expOpts: []config.Option{
				config.WithSecure(false),
				config.WithEndpoint("localhost:2135"),
				config.WithDatabase("/local"),
			},
			expConnectorOpts: []ConnectorOption{
				WithDefaultQueryMode(ScriptingQueryMode),
			},
			expErr: nil,
		},
		{
			dsn: "grpc://localhost:2135/local?query_mode=scripting&go_auto_bind.table_path_prefix=path/to/tables",
			expOpts: []config.Option{
				config.WithSecure(false),
				config.WithEndpoint("localhost:2135"),
				config.WithDatabase("/local"),
			},
			expConnectorOpts: []ConnectorOption{
				WithDefaultQueryMode(ScriptingQueryMode),
				WithBind(bind.TablePathPrefix("path/to/tables")),
			},
			expErr: nil,
		},
		{
			dsn: "grpc://localhost:2135/local?query_mode=scripting&go_auto_bind=numeric&go_auto_bind.table_path_prefix=path/to/tables", //nolint:lll
			expOpts: []config.Option{
				config.WithSecure(false),
				config.WithEndpoint("localhost:2135"),
				config.WithDatabase("/local"),
			},
			expConnectorOpts: []ConnectorOption{
				WithDefaultQueryMode(ScriptingQueryMode),
				WithBind(bind.Numeric().WithTablePathPrefix("path/to/tables")),
			},
			expErr: nil,
		},
		{
			dsn: "grpc://localhost:2135/local?query_mode=scripting&go_auto_bind=positional&go_auto_bind.table_path_prefix=path/to/tables", //nolint:lll
			expOpts: []config.Option{
				config.WithSecure(false),
				config.WithEndpoint("localhost:2135"),
				config.WithDatabase("/local"),
			},
			expConnectorOpts: []ConnectorOption{
				WithDefaultQueryMode(ScriptingQueryMode),
				WithBind(bind.Positional().WithTablePathPrefix("path/to/tables")),
			},
			expErr: nil,
		},
		{
			dsn: "grpc://localhost:2135/local?query_mode=scripting&go_auto_bind=declare&go_auto_bind.table_path_prefix=path/to/tables", //nolint:lll
			expOpts: []config.Option{
				config.WithSecure(false),
				config.WithEndpoint("localhost:2135"),
				config.WithDatabase("/local"),
			},
			expConnectorOpts: []ConnectorOption{
				WithDefaultQueryMode(ScriptingQueryMode),
				WithBind(bind.Declare().WithTablePathPrefix("path/to/tables")),
			},
			expErr: nil,
		},
		{
			dsn: "grpc://localhost:2135/local?query_mode=scripting&go_auto_bind.table_path_prefix=path/to/tables",
			expOpts: []config.Option{
				config.WithSecure(false),
				config.WithEndpoint("localhost:2135"),
				config.WithDatabase("/local"),
			},
			expConnectorOpts: []ConnectorOption{
				WithDefaultQueryMode(ScriptingQueryMode),
				WithBind(bind.TablePathPrefix("path/to/tables")),
			},
			expErr: nil,
		},
	} {
		t.Run("", func(t *testing.T) {
			opts, connectorOpts, err := Parse(tt.dsn)
			if tt.expErr != nil {
				require.ErrorIs(t, err, tt.expErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, newConnector(tt.expConnectorOpts...), newConnector(connectorOpts...))
				compareConfigs(t, config.New(tt.expOpts...), config.New(opts...))
			}
		})
	}
}
