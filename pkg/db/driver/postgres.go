package driver

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// workspaceContext contains the state that must propagate across process
// boundaries.
type workspaceContext struct {
	id string
}

const contextKey = "::TENANT_ID::"

// fromContext returns the workspaceContext stored in a context, or an error if
// there isn't one.
func fromContext(ctx context.Context) (*workspaceContext, error) {
	workspaceCtx, ok := ctx.Value(contextKey).(*workspaceContext)
	if !ok {
		return nil, errors.New("tenancy: unable to retrieve workspace context")
	}

	return workspaceCtx, nil
}

// WithID returns a new context with the given workspace id attached.
func WithID(parent context.Context, id string) context.Context {
	return context.WithValue(parent, contextKey, &workspaceContext{id: id})
}

// Driver is the Postgres database driver for Multi-Tenancy.
type Driver struct{}

// Open opens a new connection to the database. name is a connection string.
// Most users should only use it through database/sql package from the standard
// library.
func (d *Driver) Open(name string) (driver.Conn, error) {
	return open(name)
}

func init() {
	sql.Register("postgres-tenancy", &Driver{})
}

type conn struct {
	driver.Conn
}

func open(name string) (driver.Conn, error) {
	c, err := pq.Open(name)
	if err != nil {
		return nil, err
	}

	return &conn{
		Conn: c,
	}, nil
}

// Prepare implements driver.Conn.Prepare.
func (c *conn) Prepare(query string) (driver.Stmt, error) {
	return c.Conn.Prepare(query)
}

// Close implements driver.Conn.Close.
func (c *conn) Close() error {
	return c.Conn.Close()
}

// BeginTx implements driver.ConnBeginTx.BeginTx.
func (c *conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return c.Conn.(driver.ConnBeginTx).BeginTx(ctx, opts)
}

// Query implements driver.Queryer.Query.
func (c *conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	return nil, errors.New("driver.Queryer.Query not supported")
}

// QueryContext implements driver.QueryerContext.QueryContext.
func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	workspaceCtx, err := fromContext(ctx)
	if err != nil {
		return nil, err
	}

	useStmt := useStatement(workspaceCtx.id)
	if len(args) > 0 {
		if _, err := c.Conn.(driver.QueryerContext).QueryContext(ctx, useStmt, nil); err != nil {
			return nil, err
		}
	} else {
		query = useStmt + ";" + query
	}
	return c.Conn.(driver.QueryerContext).QueryContext(ctx, query, args)
}

// Exec implements driver.Execer.Exec.
func (c *conn) Exec(query string, args []driver.Value) (driver.Result, error) {
	return nil, errors.New("driver.Execer.Exec not supported")
}

// ExecContext implements driver.ExecerContext.ExecContext.
func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	workspaceCtx, err := fromContext(ctx)
	if err != nil {
		return nil, err
	}

	useStmt := useStatement(workspaceCtx.id)
	if len(args) > 0 {
		if _, err := c.Conn.(driver.ExecerContext).ExecContext(ctx, useStmt, nil); err != nil {
			return nil, err
		}
	} else {
		query = useStmt + ";" + query
	}
	return c.Conn.(driver.ExecerContext).ExecContext(ctx, query, args)
}

func useStatement(workspaceID string) string {
	// escape quotes
	pos := 0
	buf := make([]byte, len(workspaceID)*2)
	for i := 0; i < len(workspaceID); i++ {
		c := workspaceID[i]
		if c == '\'' {
			buf[pos] = '\''
			buf[pos+1] = '\''
			pos += 2
		} else {
			buf[pos] = c
			pos++
		}
	}

	return fmt.Sprintf("SET app.current_tenant='%s'", string(buf[:pos]))
}
