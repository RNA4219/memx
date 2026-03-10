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
	"path/filepath"
	"strings"

	"memx/api"
	"memx/db"
	"memx/service"
)

func main() {
	log.SetFlags(0)
	if err := loadDotEnvFromHierarchy(); err != nil {
		log.Printf("warning: failed to load .env: %v", err)
	}
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
	case "docs":
		cmdDocs(os.Args[2:])
	case "gc":
		cmdGC(os.Args[2:])
	case "summarize":
		cmdSummarize(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
}

func loadDotEnvFromHierarchy() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	path, ok := findDotEnvPath(cwd)
	if !ok {
		return nil
	}
	return loadDotEnvFile(path)
}

func findDotEnvPath(start string) (string, bool) {
	dir := start
	for {
		candidate := filepath.Join(dir, ".env")
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}

func loadDotEnvFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := parseDotEnvValue(parts[1])
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func parseDotEnvValue(raw string) string {
	value := strings.TrimSpace(raw)
	if len(value) >= 2 {
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}
	}
	return value
}
func usage() {
	fmt.Fprint(os.Stderr, `mem - memx CLI (v1.4)

Usage:
  mem api serve   [--addr 127.0.0.1:7766] [--short short.db] ...
  mem in short    --title TITLE [--body BODY | --stdin] [--tag TAG ...] [--no-llm]
  mem in journal --title TITLE --body BODY --scope SCOPE [--tag TAG ...] [--no-llm]
  mem in knowledge --title TITLE --body BODY --scope SCOPE [--tag TAG ...] [--pinned] [--no-llm]
  mem out search  QUERY
  mem out show    NOTE_ID
  mem out recall  QUERY [--top-k N] [--range N] [--stores short,journal,knowledge] [--fallback-fts]
  mem out journal search QUERY
  mem out journal show ID
  mem out journal list --scope SCOPE
  mem out knowledge search QUERY
  mem out knowledge show ID
  mem out knowledge list --scope SCOPE
  mem out knowledge pinned [--scope SCOPE]
  mem out knowledge pin ID
  mem out knowledge unpin ID
  mem out archive list
  mem out archive show ID
  mem out archive restore ID
  mem docs ingest   --title TITLE --body BODY --doc-type TYPE [--version VER] [--feature KEY ...]
  mem docs resolve  --feature KEY | --task-id ID | --topic QUERY [--limit N]
  mem docs chunks   --doc-id ID [--heading H] [--query Q] [--limit N]
  mem docs search   QUERY [--limit N]
  mem docs ack      --task-id ID --doc-id ID [--version VER]
  mem docs stale    --task-id ID
  mem docs contract --feature KEY | --task-id ID
  mem gc short    [--dry-run] [--enable-gc]
  mem summarize   NOTE_ID [--json]
  mem summarize   --ids ID1,ID2,... [--json]

Global (for client-mode):
  --api-url http://127.0.0.1:7766   # if set, CLI calls HTTP API
  --json                           # JSON output

DB flags (in-proc / server):
  --short short.db
  --journal journal.db
  --knowledge knowledge.db
  --archive archive.db

GC flags:
  --dry-run      Show planned operations without executing
  --enable-gc    Enable GC execution (disabled by default)

Recall flags:
  --top-k N         Number of top results (default: 10)
  --range N         Context range around anchor (default: 5)
  --stores STORES   Comma-separated stores: short,journal,knowledge
  --fallback-fts    Fallback to FTS if embedding unavailable
`)
}

// -------------------- common --------------------

type commonFlags struct {
	apiURL string
	json   bool

	short     string
	journal   string
	knowledge string
	archive   string
}

func (c *commonFlags) bind(fs *flag.FlagSet) {
	fs.StringVar(&c.apiURL, "api-url", "", "HTTP API base URL (if set, CLI uses HTTP client)")
	fs.BoolVar(&c.json, "json", false, "output JSON")
	fs.StringVar(&c.short, "short", "short.db", "path to short.db")
	fs.StringVar(&c.journal, "journal", "journal.db", "path to journal.db")
	fs.StringVar(&c.knowledge, "knowledge", "knowledge.db", "path to knowledge.db")
	fs.StringVar(&c.archive, "archive", "archive.db", "path to archive.db")
}

func (c *commonFlags) paths() db.Paths {
	return db.Paths{
		Short:     c.short,
		Journal:   c.journal,
		Knowledge: c.knowledge,
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
	if err := attachOpenAIFromEnv(svc); err != nil {
		_ = svc.Close()
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

func attachOpenAIFromEnv(svc *service.Service) error {
	return svc.ConfigureLLMsFromEnv()
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
		if err := attachOpenAIFromEnv(svc); err != nil {
			_ = svc.Close()
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
		noLLM := fs.Bool("no-llm", false, "skip LLM auto-summary")
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
			NoLLM:       *noLLM,
		})
		if apiErr != nil {
			log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
		}

		if cf.json {
			printJSON(resp)
			return
		}
		fmt.Printf("ok id=%s\n", resp.Note.ID)

	case "journal":
		fs := flag.NewFlagSet("mem in journal", flag.ExitOnError)
		cf := &commonFlags{}
		cf.bind(fs)
		title := fs.String("title", "", "note title")
		body := fs.String("body", "", "note body")
		stdin := fs.Bool("stdin", false, "read body from stdin")
		noLLM := fs.Bool("no-llm", false, "skip LLM auto-summary")
		scope := fs.String("scope", "", "working scope (required)")
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

		if *scope == "" {
			log.Fatal("--scope is required for journal")
		}

		ctx := context.Background()
		client, cleanup, err := makeClient(ctx, cf)
		if err != nil {
			log.Fatal(err)
		}
		defer cleanup()

		resp, apiErr := client.JournalIngest(ctx, api.JournalIngestRequest{
			Title:        *title,
			Body:         b,
			SourceType:   *sourceType,
			Origin:       *origin,
			SourceTrust:  *sourceTrust,
			Sensitivity:  *sensitivity,
			Tags:         tags,
			WorkingScope: *scope,
			NoLLM:        *noLLM,
		})
		if apiErr != nil {
			log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
		}

		if cf.json {
			printJSON(resp)
			return
		}
		fmt.Printf("ok id=%s\n", resp.Note.ID)

	case "knowledge":
		fs := flag.NewFlagSet("mem in knowledge", flag.ExitOnError)
		cf := &commonFlags{}
		cf.bind(fs)
		title := fs.String("title", "", "note title")
		body := fs.String("body", "", "note body")
		stdin := fs.Bool("stdin", false, "read body from stdin")
		noLLM := fs.Bool("no-llm", false, "skip LLM auto-summary")
		scope := fs.String("scope", "", "working scope (required)")
		tags := multiStringFlag{}
		fs.Var(&tags, "tag", "tag (repeatable)")
		pinned := fs.Bool("pinned", false, "pin the note")
		sourceType := fs.String("source-type", "manual", "source type")
		origin := fs.String("origin", "", "origin")
		sourceTrust := fs.String("source-trust", "user_input", "source trust")
		sensitivity := fs.String("sensitivity", "internal", "sensitivity")
		_ = fs.Parse(args[1:])

		b := *body
		if *stdin {
			b = readAllStdin()
		}

		if *scope == "" {
			log.Fatal("--scope is required for knowledge")
		}

		ctx := context.Background()
		client, cleanup, err := makeClient(ctx, cf)
		if err != nil {
			log.Fatal(err)
		}
		defer cleanup()

		resp, apiErr := client.KnowledgeIngest(ctx, api.KnowledgeIngestRequest{
			Title:        *title,
			Body:         b,
			SourceType:   *sourceType,
			Origin:       *origin,
			SourceTrust:  *sourceTrust,
			Sensitivity:  *sensitivity,
			Tags:         tags,
			WorkingScope: *scope,
			IsPinned:     *pinned,
			NoLLM:        *noLLM,
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

		resp, apiErr := searchAcrossStores(ctx, client, query, *topK)
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

		n, apiErr := resolveNoteAcrossStores(ctx, client, id)
		if apiErr != nil {
			log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
		}
		if cf.json {
			printJSON(n)
			return
		}
		printResolvedNote(n)

	case "recall":
		cmdOutRecall(args[1:])

	case "journal":
		cmdOutJournal(args[1:])

	case "knowledge":
		cmdOutKnowledge(args[1:])

	case "archive":
		cmdOutArchive(args[1:])

	default:
		usage()
		os.Exit(2)
	}
}

// -------------------- recall --------------------

func cmdOutRecall(args []string) {
	fs := flag.NewFlagSet("mem out recall", flag.ExitOnError)
	cf := &commonFlags{}
	cf.bind(fs)
	topK := fs.Int("k", 10, "top k results")
	msgRange := fs.Int("range", 5, "message range for context")
	stores := fs.String("stores", "short", "comma-separated stores to search (short,journal,knowledge)")
	fallbackFTS := fs.Bool("fallback-fts", false, "use FTS fallback if embedding unavailable")
	_ = fs.Parse(args)
	query := strings.Join(fs.Args(), " ")

	if query == "" {
		log.Fatal("QUERY is required")
	}

	ctx := context.Background()
	client, cleanup, err := makeClient(ctx, cf)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	storeList := parseRecallStores(*stores)

	resp, apiErr := client.Recall(ctx, api.RecallRequest{
		Query:        query,
		TopK:         *topK,
		MessageRange: *msgRange,
		Stores:       storeList,
		FallbackFTS:  *fallbackFTS,
	})
	if apiErr != nil {
		log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
	}

	if cf.json {
		printJSON(resp)
		return
	}

	for _, nwc := range resp.Results {
		anchor := nwc.Anchor
		fmt.Printf("=== [%s] %s (score: %.4f) ===\n", anchor.Store, anchor.Title, anchor.Score)
		if anchor.Summary != "" {
			fmt.Printf("Summary: %s\n", anchor.Summary)
		}
		fmt.Printf("ID: %s\n", anchor.ID)
		if len(nwc.Before) > 0 {
			fmt.Printf("\n--- Before (%d) ---\n", len(nwc.Before))
			for _, b := range nwc.Before {
				fmt.Printf("  [%s] %s\n", b.ID, b.Title)
			}
		}
		if len(nwc.After) > 0 {
			fmt.Printf("\n--- After (%d) ---\n", len(nwc.After))
			for _, a := range nwc.After {
				fmt.Printf("  [%s] %s\n", a.ID, a.Title)
			}
		}
		fmt.Println()
	}
}

func parseRecallStores(s string) []string {
	if s == "" {
		return []string{"short"}
	}
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(strings.ToLower(p))
		if p == "short" || p == "journal" || p == "knowledge" {
			result = append(result, p)
		}
	}
	if len(result) == 0 {
		return []string{"short"}
	}
	return result
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
			Target:  "short",
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

func parseStores(s string) []string {
	if s == "" {
		return []string{"short"}
	}
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// -------------------- journal commands --------------------

func cmdOutJournal(args []string) {
	if len(args) < 1 {
		usage()
		os.Exit(2)
	}
	switch args[0] {
	case "search":
		fs := flag.NewFlagSet("mem out journal search", flag.ExitOnError)
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

		resp, apiErr := client.JournalSearch(ctx, api.JournalSearchRequest{Query: query, TopK: *topK})
		if apiErr != nil {
			log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
		}
		if cf.json {
			printJSON(resp)
			return
		}
		for _, n := range resp.Notes {
			fmt.Printf("%s\t%s\t[%s]\n", n.ID, n.Title, n.WorkingScope)
		}

	case "show":
		fs := flag.NewFlagSet("mem out journal show", flag.ExitOnError)
		cf := &commonFlags{}
		cf.bind(fs)
		_ = fs.Parse(args[1:])
		if len(fs.Args()) < 1 {
			log.Fatal("ID is required")
		}
		id := fs.Args()[0]

		ctx := context.Background()
		client, cleanup, err := makeClient(ctx, cf)
		if err != nil {
			log.Fatal(err)
		}
		defer cleanup()

		n, apiErr := client.JournalGet(ctx, id)
		if apiErr != nil {
			log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
		}
		if cf.json {
			printJSON(n)
			return
		}
		fmt.Printf("# %s\n\nScope: %s\n\n%s\n", n.Title, n.WorkingScope, n.Body)

	case "list":
		fs := flag.NewFlagSet("mem out journal list", flag.ExitOnError)
		cf := &commonFlags{}
		cf.bind(fs)
		scope := fs.String("scope", "", "working scope (required)")
		limit := fs.Int("limit", 20, "limit")
		_ = fs.Parse(args[1:])

		if *scope == "" {
			log.Fatal("--scope is required")
		}

		ctx := context.Background()
		client, cleanup, err := makeClient(ctx, cf)
		if err != nil {
			log.Fatal(err)
		}
		defer cleanup()

		resp, apiErr := client.JournalListByScope(ctx, api.JournalListByScopeRequest{
			WorkingScope: *scope,
			Limit:        *limit,
		})
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

	default:
		usage()
		os.Exit(2)
	}
}

// -------------------- knowledge commands --------------------

func cmdOutKnowledge(args []string) {
	if len(args) < 1 {
		usage()
		os.Exit(2)
	}
	switch args[0] {
	case "search":
		fs := flag.NewFlagSet("mem out knowledge search", flag.ExitOnError)
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

		resp, apiErr := client.KnowledgeSearch(ctx, api.KnowledgeSearchRequest{Query: query, TopK: *topK})
		if apiErr != nil {
			log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
		}
		if cf.json {
			printJSON(resp)
			return
		}
		for _, n := range resp.Notes {
			pinned := ""
			if n.IsPinned {
				pinned = " [pinned]"
			}
			fmt.Printf("%s\t%s\t[%s]%s\n", n.ID, n.Title, n.WorkingScope, pinned)
		}

	case "show":
		fs := flag.NewFlagSet("mem out knowledge show", flag.ExitOnError)
		cf := &commonFlags{}
		cf.bind(fs)
		_ = fs.Parse(args[1:])
		if len(fs.Args()) < 1 {
			log.Fatal("ID is required")
		}
		id := fs.Args()[0]

		ctx := context.Background()
		client, cleanup, err := makeClient(ctx, cf)
		if err != nil {
			log.Fatal(err)
		}
		defer cleanup()

		n, apiErr := client.KnowledgeGet(ctx, id)
		if apiErr != nil {
			log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
		}
		if cf.json {
			printJSON(n)
			return
		}
		pinned := ""
		if n.IsPinned {
			pinned = " (pinned)"
		}
		fmt.Printf("# %s%s\n\nScope: %s\n\n%s\n", n.Title, pinned, n.WorkingScope, n.Body)

	case "list":
		fs := flag.NewFlagSet("mem out knowledge list", flag.ExitOnError)
		cf := &commonFlags{}
		cf.bind(fs)
		scope := fs.String("scope", "", "working scope (required)")
		limit := fs.Int("limit", 20, "limit")
		_ = fs.Parse(args[1:])

		if *scope == "" {
			log.Fatal("--scope is required")
		}

		ctx := context.Background()
		client, cleanup, err := makeClient(ctx, cf)
		if err != nil {
			log.Fatal(err)
		}
		defer cleanup()

		resp, apiErr := client.KnowledgeListByScope(ctx, api.KnowledgeListByScopeRequest{
			WorkingScope: *scope,
			Limit:        *limit,
		})
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

	case "pinned":
		fs := flag.NewFlagSet("mem out knowledge pinned", flag.ExitOnError)
		cf := &commonFlags{}
		cf.bind(fs)
		scope := fs.String("scope", "", "working scope (optional)")
		limit := fs.Int("limit", 20, "limit")
		_ = fs.Parse(args[1:])

		ctx := context.Background()
		client, cleanup, err := makeClient(ctx, cf)
		if err != nil {
			log.Fatal(err)
		}
		defer cleanup()

		resp, apiErr := client.KnowledgeListPinned(ctx, api.KnowledgeListPinnedRequest{
			WorkingScope: *scope,
			Limit:        *limit,
		})
		if apiErr != nil {
			log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
		}
		if cf.json {
			printJSON(resp)
			return
		}
		for _, n := range resp.Notes {
			fmt.Printf("%s\t%s\t[%s]\n", n.ID, n.Title, n.WorkingScope)
		}

	case "pin":
		fs := flag.NewFlagSet("mem out knowledge pin", flag.ExitOnError)
		cf := &commonFlags{}
		cf.bind(fs)
		_ = fs.Parse(args[1:])
		if len(fs.Args()) < 1 {
			log.Fatal("ID is required")
		}
		id := fs.Args()[0]

		ctx := context.Background()
		client, cleanup, err := makeClient(ctx, cf)
		if err != nil {
			log.Fatal(err)
		}
		defer cleanup()

		_, apiErr := client.KnowledgePin(ctx, id)
		if apiErr != nil {
			log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
		}
		fmt.Println("ok")

	case "unpin":
		fs := flag.NewFlagSet("mem out knowledge unpin", flag.ExitOnError)
		cf := &commonFlags{}
		cf.bind(fs)
		_ = fs.Parse(args[1:])
		if len(fs.Args()) < 1 {
			log.Fatal("ID is required")
		}
		id := fs.Args()[0]

		ctx := context.Background()
		client, cleanup, err := makeClient(ctx, cf)
		if err != nil {
			log.Fatal(err)
		}
		defer cleanup()

		_, apiErr := client.KnowledgeUnpin(ctx, id)
		if apiErr != nil {
			log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
		}
		fmt.Println("ok")

	default:
		usage()
		os.Exit(2)
	}
}

// -------------------- archive commands --------------------

func cmdOutArchive(args []string) {
	if len(args) < 1 {
		usage()
		os.Exit(2)
	}
	switch args[0] {
	case "list":
		fs := flag.NewFlagSet("mem out archive list", flag.ExitOnError)
		cf := &commonFlags{}
		cf.bind(fs)
		limit := fs.Int("limit", 20, "limit")
		_ = fs.Parse(args[1:])

		ctx := context.Background()
		client, cleanup, err := makeClient(ctx, cf)
		if err != nil {
			log.Fatal(err)
		}
		defer cleanup()

		resp, apiErr := client.ArchiveList(ctx, api.ArchiveListRequest{Limit: *limit})
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
		fs := flag.NewFlagSet("mem out archive show", flag.ExitOnError)
		cf := &commonFlags{}
		cf.bind(fs)
		_ = fs.Parse(args[1:])
		if len(fs.Args()) < 1 {
			log.Fatal("ID is required")
		}
		id := fs.Args()[0]

		ctx := context.Background()
		client, cleanup, err := makeClient(ctx, cf)
		if err != nil {
			log.Fatal(err)
		}
		defer cleanup()

		n, apiErr := client.ArchiveGet(ctx, id)
		if apiErr != nil {
			log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
		}
		if cf.json {
			printJSON(n)
			return
		}
		fmt.Printf("# %s\n\n%s\n", n.Title, n.Body)

	case "restore":
		fs := flag.NewFlagSet("mem out archive restore", flag.ExitOnError)
		cf := &commonFlags{}
		cf.bind(fs)
		_ = fs.Parse(args[1:])
		if len(fs.Args()) < 1 {
			log.Fatal("ID is required")
		}
		id := fs.Args()[0]

		ctx := context.Background()
		client, cleanup, err := makeClient(ctx, cf)
		if err != nil {
			log.Fatal(err)
		}
		defer cleanup()

		resp, apiErr := client.ArchiveRestore(ctx, id)
		if apiErr != nil {
			log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
		}
		if cf.json {
			printJSON(resp)
			return
		}
		fmt.Printf("restored: %s\n", resp.Note.ID)

	default:
		usage()
		os.Exit(2)
	}
}

// -------------------- docs commands --------------------

func cmdDocs(args []string) {
	if len(args) < 1 {
		usage()
		os.Exit(2)
	}
	switch args[0] {
	case "ingest":
		cmdDocsIngest(args[1:])
	case "resolve":
		cmdDocsResolve(args[1:])
	case "chunks":
		cmdDocsChunks(args[1:])
	case "search":
		cmdDocsSearch(args[1:])
	case "ack":
		cmdDocsAck(args[1:])
	case "stale":
		cmdDocsStale(args[1:])
	case "contract":
		cmdDocsContract(args[1:])
	default:
		usage()
		os.Exit(2)
	}
}

func cmdDocsIngest(args []string) {
	fs := flag.NewFlagSet("mem docs ingest", flag.ExitOnError)
	cf := &commonFlags{}
	cf.bind(fs)
	title := fs.String("title", "", "document title")
	body := fs.String("body", "", "document body")
	docType := fs.String("doc-type", "", "document type (spec, requirement, design, etc.)")
	version := fs.String("version", "", "document version")
	sourcePath := fs.String("source-path", "", "source file path")
	var features flagsStringSlice
	fs.Var(&features, "feature", "feature keys (can be repeated)")
	var tags flagsStringSlice
	fs.Var(&tags, "tag", "tags (can be repeated)")
	_ = fs.Parse(args)

	if *title == "" || *body == "" || *docType == "" {
		log.Fatal("--title, --body, and --doc-type are required")
	}
	if *version == "" {
		*version = "1"
	}

	ctx := context.Background()
	client, cleanup, err := makeClient(ctx, cf)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	req := api.DocsIngestRequest{
		Title:       *title,
		Body:        *body,
		DocType:     *docType,
		Version:     *version,
		SourcePath:  *sourcePath,
		FeatureKeys: features,
		Tags:        tags,
	}
	resp, apiErr := client.DocsIngest(ctx, req)
	if apiErr != nil {
		log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
	}
	if cf.json {
		printJSON(resp)
		return
	}
	fmt.Printf("ingested: %s (version %s, %d chunks)\n", resp.DocID, resp.Version, resp.ChunkCount)
}

func cmdDocsResolve(args []string) {
	fs := flag.NewFlagSet("mem docs resolve", flag.ExitOnError)
	cf := &commonFlags{}
	cf.bind(fs)
	feature := fs.String("feature", "", "feature key")
	taskID := fs.String("task-id", "", "task ID")
	topic := fs.String("topic", "", "topic query")
	limit := fs.Int("limit", 10, "max results")
	_ = fs.Parse(args)

	if *feature == "" && *taskID == "" && *topic == "" {
		log.Fatal("one of --feature, --task-id, or --topic is required")
	}

	ctx := context.Background()
	client, cleanup, err := makeClient(ctx, cf)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	req := api.DocsResolveRequest{
		Feature: *feature,
		TaskID:  *taskID,
		Topic:   *topic,
		Limit:   *limit,
	}
	resp, apiErr := client.DocsResolve(ctx, req)
	if apiErr != nil {
		log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
	}
	if cf.json {
		printJSON(resp)
		return
	}
	if len(resp.Required) > 0 {
		fmt.Println("Required:")
		for _, e := range resp.Required {
			fmt.Printf("  %s\t%s\t%s\n", e.DocID, e.Version, e.Title)
		}
	}
	if len(resp.Recommended) > 0 {
		fmt.Println("Recommended:")
		for _, e := range resp.Recommended {
			fmt.Printf("  %s\t%s\t%s\n", e.DocID, e.Version, e.Title)
		}
	}
}

func cmdDocsChunks(args []string) {
	fs := flag.NewFlagSet("mem docs chunks", flag.ExitOnError)
	cf := &commonFlags{}
	cf.bind(fs)
	docID := fs.String("doc-id", "", "document ID")
	heading := fs.String("heading", "", "filter by heading")
	query := fs.String("query", "", "filter by query")
	limit := fs.Int("limit", 0, "max chunks")
	_ = fs.Parse(args)

	if *docID == "" {
		log.Fatal("--doc-id is required")
	}

	ctx := context.Background()
	client, cleanup, err := makeClient(ctx, cf)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	req := api.ChunksGetRequest{
		DocID:   *docID,
		Heading: *heading,
		Query:   *query,
		Limit:   *limit,
	}
	resp, apiErr := client.ChunksGet(ctx, req)
	if apiErr != nil {
		log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
	}
	if cf.json {
		printJSON(resp)
		return
	}
	for _, c := range resp.Chunks {
		fmt.Printf("## %s\n%s\n\n", c.Heading, c.Body)
	}
}

func cmdDocsSearch(args []string) {
	fs := flag.NewFlagSet("mem docs search", flag.ExitOnError)
	cf := &commonFlags{}
	cf.bind(fs)
	limit := fs.Int("limit", 10, "max results")
	_ = fs.Parse(args)
	if len(fs.Args()) < 1 {
		log.Fatal("query is required")
	}
	query := fs.Args()[0]

	ctx := context.Background()
	client, cleanup, err := makeClient(ctx, cf)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	req := api.DocsSearchRequest{Query: query, Limit: *limit}
	resp, apiErr := client.DocsSearch(ctx, req)
	if apiErr != nil {
		log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
	}
	if cf.json {
		printJSON(resp)
		return
	}
	for _, e := range resp.Results {
		fmt.Printf("%s\t%s\t%s\n", e.DocID, e.Version, e.Title)
	}
}

func cmdDocsAck(args []string) {
	fs := flag.NewFlagSet("mem docs ack", flag.ExitOnError)
	cf := &commonFlags{}
	cf.bind(fs)
	taskID := fs.String("task-id", "", "task ID")
	docID := fs.String("doc-id", "", "document ID")
	version := fs.String("version", "", "document version")
	_ = fs.Parse(args)

	if *taskID == "" || *docID == "" {
		log.Fatal("--task-id and --doc-id are required")
	}

	ctx := context.Background()
	client, cleanup, err := makeClient(ctx, cf)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	req := api.ReadsAckRequest{
		TaskID:  *taskID,
		DocID:   *docID,
		Version: *version,
		Reader:  "cli",
	}
	resp, apiErr := client.ReadsAck(ctx, req)
	if apiErr != nil {
		log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
	}
	if cf.json {
		printJSON(resp)
		return
	}
	fmt.Printf("acknowledged: %s -> %s (%s)\n", resp.Receipt.TaskID, resp.Receipt.DocID, resp.Receipt.Version)
}

func cmdDocsStale(args []string) {
	fs := flag.NewFlagSet("mem docs stale", flag.ExitOnError)
	cf := &commonFlags{}
	cf.bind(fs)
	taskID := fs.String("task-id", "", "task ID")
	_ = fs.Parse(args)

	if *taskID == "" {
		log.Fatal("--task-id is required")
	}

	ctx := context.Background()
	client, cleanup, err := makeClient(ctx, cf)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	req := api.DocsStaleCheckRequest{TaskID: *taskID}
	resp, apiErr := client.DocsStaleCheck(ctx, req)
	if apiErr != nil {
		log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
	}
	if cf.json {
		printJSON(resp)
		return
	}
	if len(resp.Stale) == 0 {
		fmt.Printf("task %s is fresh\n", resp.TaskID)
		return
	}
	fmt.Printf("task %s is stale:\n", resp.TaskID)
	for _, s := range resp.Stale {
		fmt.Printf("  %s: %s -> %s (%s)\n", s.DocID, s.PreviousVersion, s.CurrentVersion, s.Reason)
	}
}

func cmdDocsContract(args []string) {
	fs := flag.NewFlagSet("mem docs contract", flag.ExitOnError)
	cf := &commonFlags{}
	cf.bind(fs)
	feature := fs.String("feature", "", "feature key")
	taskID := fs.String("task-id", "", "task ID")
	_ = fs.Parse(args)

	if *feature == "" && *taskID == "" {
		log.Fatal("one of --feature or --task-id is required")
	}

	ctx := context.Background()
	client, cleanup, err := makeClient(ctx, cf)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	req := api.ContractsResolveRequest{
		Feature: *feature,
		TaskID:  *taskID,
	}
	resp, apiErr := client.ContractsResolve(ctx, req)
	if apiErr != nil {
		log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
	}
	if cf.json {
		printJSON(resp)
		return
	}
	if len(resp.AcceptanceCriteria) > 0 {
		fmt.Println("Acceptance Criteria:")
		for _, a := range resp.AcceptanceCriteria {
			fmt.Printf("  - %s\n", a)
		}
	}
	if len(resp.ForbiddenPatterns) > 0 {
		fmt.Println("Forbidden Patterns:")
		for _, f := range resp.ForbiddenPatterns {
			fmt.Printf("  - %s\n", f)
		}
	}
	if len(resp.DefinitionOfDone) > 0 {
		fmt.Println("Definition of Done:")
		for _, d := range resp.DefinitionOfDone {
			fmt.Printf("  - %s\n", d)
		}
	}
}

// flagsStringSlice implements flag.Value for repeated string flags
type flagsStringSlice []string

func (f *flagsStringSlice) String() string { return strings.Join(*f, ",") }
func (f *flagsStringSlice) Set(v string) error {
	*f = append(*f, v)
	return nil
}
