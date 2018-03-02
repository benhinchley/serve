package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/benhinchley/cmd"
)

func main() {
	p, err := cmd.NewProgram("serve", "simple local webserver", &serveCommand{}, nil)
	if err != nil {
		fatalPrint(err)
	}
	if err := p.ParseArgs(os.Args); err != nil {
		fatalPrint(err)
	}
	if err := p.Run(func(env *cmd.Environment, c cmd.Command, args []string) error {
		if err := c.Run(env.GetDefaultContext(), args); err != nil {
			return fmt.Errorf("%s: %v", c.Name(), err)
		}
		return nil
	}); err != nil {
		fatalPrint(err)
	}
}

func fatalPrint(m error) {
	fmt.Fprintf(os.Stderr, "%v", m)
	os.Exit(1)
}

type serveCommand struct {
	addr string
}

func (c *serveCommand) Name() string { return "serve" }
func (c *serveCommand) Args() string { return "[flags] <path>" }
func (c *serveCommand) Desc() string { return "" }
func (c *serveCommand) Help() string { return "" }

func (c *serveCommand) Register(fs *flag.FlagSet) {
	fs.StringVar(&c.addr, "addr", "localhost:4567", "bind address for serve")
}

func (c *serveCommand) Run(ctx cmd.Context, args []string) error {
	var p string
	if len(args) < 1 {
		p = "."
	} else {
		p = args[0]
	}

	ctx.Stdout().Printf("server running at http://%s", c.addr)
	return http.ListenAndServe(c.addr, http.FileServer(http.Dir(path.Clean(p))))
}
