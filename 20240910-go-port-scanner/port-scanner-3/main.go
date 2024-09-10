package main

import (
	"context"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"golang.org/x/xerrors"
)

const (
	wellKnownPorts = 1024
	allPorts       = 65535
)

var timeout = 3 * time.Second

// This time our command struct has a few fields, we can use these to store our flag values.
type scanCmd struct {
	host          string
	shouldScanAll bool
}

func main() {
	cli.RunRoot(new(root))
}

// We define our command as a struct.
type root struct{}

// Our command struct needs to implement the Command interface as defined by cdr/cli.
// See: https://pkg.go.dev/go.coder.com/cli#Command for more details.
func (r *root) Run(fl *pflag.FlagSet) {
	fl.Usage()
}

// We'll also need to implement a method for returning a command spec to fully satisfy the
// the aforementioned Command interface.
func (r *root) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "port-scanner",
		Usage: "[subcommand] [flags]",
		Desc:  "A simple port-scanner.",
	}
}

// We can now attach subcommands to our command struct; implementing ParentCommand as a result as defined by cdr/cli.
// See https://pkg.go.dev/go.coder.com/cli#ParentCommand for more details.
func (r *root) Subcommands() []cli.Command {
	return []cli.Command{
		new(scanCmd),
	}
}

// cdr/cli supports subcommand aliases so let's define one in our command spec to provide
// our users with a more succinct input experience.
func (cmd *scanCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "scan",
		Usage:   "[flags]",
		Aliases: []string{"s"},
		Desc:    "Scan a host for open ports.",
	}
}

// When adding flags, hang the following method-signature off your command struct to
// satisfy the FlaggedCommand interface as defined by cdr/cli.
// See https://pkg.go.dev/go.coder.com/cli#FlaggedCommand for more details.
func (cmd *scanCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVar(&cmd.host, "host", "", "host to scan(ip address)")
	fl.BoolVarP(&cmd.shouldScanAll, "all", "a", false, "scan all ports(scans first 1024 if not enabled)")
}

func (cmd *scanCmd) Run(fl *pflag.FlagSet) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if cmd.host == "" {
		fl.Usage()
		log.Fatal("host not provided")
	}

	scanner, err := newScanner(cmd.host, cmd.shouldScanAll)
	if err != nil {
		fl.Usage()
		log.Fatalf("failed to initialize port scanner: %s", err)
	}

	log.Printf("scanning %s...", cmd.host)
	start := time.Now()
	openPorts := scanner.scan(ctx)
	log.Printf("scan completed in %s", time.Since(start))

	if len(openPorts) == 0 {
		log.Printf("%q has no open ports", cmd.host)
		return
	}
	log.Printf("found %d open ports", len(openPorts))
	log.Printf("open-ports: %v", openPorts)
}

// Now let's implement our port scanner.
type scanner struct {
	// We're going to want to scan each port concurrently so let's embed
	// a mutex lock into our scanner to make sure we do this in a thread-safe way.
	sync.Mutex
	host      string
	openPorts []int
	scanAll   bool
}

func newScanner(host string, scanAll bool) (*scanner, error) {
	if net.ParseIP(host) == nil {
		return nil, xerrors.Errorf("%q is an invalid ip address", host)
	}

	return &scanner{
		Mutex:   sync.Mutex{},
		host:    host,
		scanAll: scanAll,
	}, nil
}

func (s *scanner) add(port int) {
	// Since we'll be appending to the same slice from different goroutines,
	// let's make sure we're locking and unlocking between writes.
	s.Lock()
	s.openPorts = append(s.openPorts, port)
	s.Unlock()
}

func (s *scanner) scan(ctx context.Context) []int {
	// Let's use a wait group so we can wait for all of our
	// goroutines to exit before returning our results.
	var wg sync.WaitGroup
	for _, port := range portsToScan(s.scanAll) {
		wg.Add(1)
		// Because the value of the loop-variable 'port' changes on every iteration, we'll want
		// to pass a copy of its value into each new goroutine. We don't need to do this with the
		// 'host' variable because it's not a loop-variable and its value never changes.
		go func(p int) {
			defer wg.Done()
			if isOpen(s.host, p) {
				s.add(p)
			}
		}(port)
	}
	wg.Wait()
	return s.openPorts
}

func portsToScan(shouldScanAll bool) []int {
	max := wellKnownPorts
	if shouldScanAll {
		max = allPorts
	}

	var ports []int
	for port := 1; port < max; port++ {
		ports = append(ports, port)
	}
	return ports
}

func isOpen(host string, port int) bool {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
