package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"memx/api"
	"memx/db"
	"memx/service"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "api":
		cmdAPI(os.Args[2:])
	case "in":
		cmdIn(os.Args[2:])
	case "out":
		cmdOut(os.Args[2:])
	case "gc":
		cmdGC(os.Args[2:])
	case "summarize":
		cmdSummarize(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `mem - memx CLI (v1.3)

Usage:
  mem api serve   [--addr 127.0.0.1:7766] [--short short.db] ...
  mem in short    --title TITLE [--body BODY | --stdin] [--tag TAG ...]
  mem out search  QUERY
  mem out show    NOTE_ID
  mem gc short    [--dry-run] [--enable-gc]
  mem summarize   NOTE_ID [--json]
  mem summarize   --ids ID1,ID2,... [--json]

Global (for client-mode):
  --api-url http://127.0.0.1:7766   # if set, CLI calls HTTP API
  --json                           # JSON output

DB flags (in-proc / server):
  --short short.db
  --chronicle chronicle.db
  --memopedia memopedia.db
  --archive archive.db

GC flags:
  --dry-run      Show planned operations without executing
  --enable-gc    Enable GC execution (disabled by default)
`)
}

// -------------------- common --------------------

type commonFlags struct {
	apiURL string
	json   bool

	short     string
	chronicle string
	memopedia string
	archive   string
}

func (c *commonFlags) bind(fs *flag.FlagSet) {
	fs.StringVar(&c.apiURL, "api-url", "", "HTTP API base URL (if set, CLI uses HTTP client)")
	fs.BoolVar(&c.json, "json", false, "output JSON")
	fs.StringVar(&c.short, "short", "short.db", "path to short.db")
	fs.StringVar(&c.chronicle, "chronicle", "", "path to chronicle.db")
	fs.StringVar(&c.memopedia, "memopedia", "", "path to memopedia.db")
	fs.StringVar(&c.archive, "archive", "", "path to archive.db")
}

func (c *commonFlags) paths() db.Paths {
	return db.Paths{
		Short:     c.short,
		Chronicle: c.chronicle,
		Memopedia: c.memopedia,
		Archive:   c.archive,
	}
}

func makeClient(ctx context.Context, cf *commonFlags) (api.Client, func(), error) {
	if cf.apiURL != "" {
		return api.NewHTTPClient(cf.apiURL), func() {}, nil
	}
	svc, err := service.New(cf.paths())
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() { _ = svc.Close() }
	return api.NewInProcClient(svc), cleanup, nil
}

func printJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

// -------------------- api --------------------

func cmdAPI(args []string) {
	if len(args) < 1 {
		usage()
		os.Exit(2)
	}
	switch args[0] {
	case "serve":
		fs := flag.NewFlagSet("mem api serve", flag.ExitOnError)
		cf := &commonFlags{}
		cf.bind(fs)
		addr := fs.String("addr", "127.0.0.1:7766", "listen address")
		_ = fs.Parse(args[1:])

		svc, err := service.New(cf.paths())
		if err != nil {
			log.Fatal(err)
		}
		defer func() { _ = svc.Close() }()

		srv := api.NewHTTPServer(svc)
		h := srv.Handler()

		log.Printf("memx API listening on http://%s", *addr)
		log.Fatal(http.ListenAndServe(*addr, h))
	default:
		usage()
		os.Exit(2)
	}
}

// -------------------- in --------------------

func cmdIn(args []string) {
	if len(args) < 1 {
		usage()
		os.Exit(2)
	}
	switch args[0] {
	case "short":
		fs := flag.NewFlagSet("mem in short", flag.ExitOnError)
		cf := &commonFlags{}
		cf.bind(fs)
		title := fs.String("title", "", "note title")
		body := fs.String("body", "", "note body")
		summary := fs.String("summary", "", "note summary")
		stdin := fs.Bool("stdin", false, "read body from stdin")
		tags := multiStringFlag{}
		fs.Var(&tags, "tag", "tag (repeatable)")
		sourceType := fs.String("source-type", "manual", "source type")
		origin := fs.String("origin", "", "origin")
		sourceTrust := fs.String("source-trust", "user_input", "source trust")
		sensitivity := fs.String("sensitivity", "internal", "sensitivity")
		_ = fs.Parse(args[1:])

		b := *body
		if *stdin {
			b = readAllStdin()
		}

		ctx := context.Background()
		client, cleanup, err := makeClient(ctx, cf)
		if err != nil {
			log.Fatal(err)
		}
		defer cleanup()

		resp, apiErr := client.NotesIngest(ctx, api.NotesIngestRequest{
			Title:       *title,
			Body:        b,
			Summary:     *summary,
			SourceType:  *sourceType,
			Origin:      *origin,
			SourceTrust: *sourceTrust,
			Sensitivity: *sensitivity,
			Tags:        tags,
		})
		if apiErr != nil {
			log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
		}

		if cf.json {
			printJSON(resp)
			return
		}
		fmt.Printf("ok id=%s\n", resp.Note.ID)
	default:
		usage()
		os.Exit(2)
	}
}

// -------------------- out --------------------

func cmdOut(args []string) {
	if len(args) < 1 {
		usage()
		os.Exit(2)
	}
	switch args[0] {
	case "search":
		fs := flag.NewFlagSet("mem out search", flag.ExitOnError)
		cf := &commonFlags{}
		cf.bind(fs)
		topK := fs.Int("k", 20, "top k")
		_ = fs.Parse(args[1:])
		query := strings.Join(fs.Args(), " ")

		ctx := context.Background()
		client, cleanup, err := makeClient(ctx, cf)
		if err != nil {
			log.Fatal(err)
		}
		defer cleanup()

		resp, apiErr := client.NotesSearch(ctx, api.NotesSearchRequest{Query: query, TopK: *topK})
		if apiErr != nil {
			log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
		}
		if cf.json {
			printJSON(resp)
			return
		}
		for _, n := range resp.Notes {
			fmt.Printf("%s\t%s\n", n.ID, n.Title)
		}

	case "show":
		fs := flag.NewFlagSet("mem out show", flag.ExitOnError)
		cf := &commonFlags{}
		cf.bind(fs)
		_ = fs.Parse(args[1:])
		if len(fs.Args()) < 1 {
			log.Fatal("NOTE_ID is required")
		}
		id := fs.Args()[0]

		ctx := context.Background()
		client, cleanup, err := makeClient(ctx, cf)
		if err != nil {
			log.Fatal(err)
		}
		defer cleanup()

		n, apiErr := client.NotesGet(ctx, id)
		if apiErr != nil {
			log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
		}
		if cf.json {
			printJSON(n)
			return
		}
		fmt.Printf("# %s\n\n%s\n", n.Title, n.Body)

	default:
		usage()
		os.Exit(2)
	}
}

// -------------------- summarize --------------------

func cmdSummarize(args []string) {
	fs := flag.NewFlagSet("mem summarize", flag.ExitOnError)
	cf := &commonFlags{}
	cf.bind(fs)
	ids := fs.String("ids", "", "comma-separated note IDs for batch summarization")
	_ = fs.Parse(args)

	ctx := context.Background()
	client, cleanup, err := makeClient(ctx, cf)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	// Batch summarization
	if *ids != "" {
		idList := strings.Split(*ids, ",")
		for i, id := range idList {
			idList[i] = strings.TrimSpace(id)
		}
		resp, apiErr := client.SummarizeBatch(ctx, api.SummarizeBatchRequest{IDs: idList})
		if apiErr != nil {
			log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
		}
		if cf.json {
			printJSON(resp)
			return
		}
		fmt.Printf("Summary (%d notes):\n%s\n", resp.NoteCount, resp.Summary)
		return
	}

	// Single note summarization
	if len(fs.Args()) < 1 {
		log.Fatal("NOTE_ID is required (or use --ids for batch)")
	}
	id := fs.Args()[0]

	resp, apiErr := client.Summarize(ctx, id)
	if apiErr != nil {
		log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
	}
	if cf.json {
		printJSON(resp)
		return
	}
	fmt.Printf("# %s\n\nSummary: %s\n", resp.Note.Title, resp.Note.Summary)
}

// -------------------- gc --------------------

func cmdGC(args []string) {
	if len(args) < 1 {
		usage()
		os.Exit(2)
	}
	switch args[0] {
	case "short":
		fs := flag.NewFlagSet("mem gc short", flag.ExitOnError)
		cf := &commonFlags{}
		cf.bind(fs)
		dryRun := fs.Bool("dry-run", false, "dry-run mode (show planned operations without executing)")
		enableGC := fs.Bool("enable-gc", false, "enable GC execution (disabled by default)")
		_ = fs.Parse(args[1:])

		ctx := context.Background()
		client, cleanup, err := makeClient(ctx, cf)
		if err != nil {
			log.Fatal(err)
		}
		defer cleanup()

		// feature flag チェック
		if !*enableGC && !*dryRun {
			log.Fatal("GC is disabled by default. Use --enable-gc to enable or --dry-run to preview.")
		}

		resp, apiErr := client.GCRun(ctx, api.GCRunRequest{
			Target: "short",
			Options: api.GCOptions{DryRun: *dryRun},
		})
		if apiErr != nil {
			log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
		}

		if cf.json || *dryRun {
			fmt.Println(resp.Status)
			return
		}
		fmt.Printf("GC completed: %s\n", resp.Status)

	default:
		usage()
		os.Exit(2)
	}
}

// -------------------- helpers --------------------

type multiStringFlag []string

func (m *multiStringFlag) String() string { return strings.Join(*m, ",") }
func (m *multiStringFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

func readAllStdin() string {
	b := &strings.Builder{}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		b.WriteString(s.Text())
		b.WriteByte('\n')
	}
	return b.String()
}
