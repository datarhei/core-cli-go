package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	coreclient "github.com/datarhei/core-client-go/v16"
	coreclientapi "github.com/datarhei/core-client-go/v16/api"
	"github.com/itchyny/gojq"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/mattn/go-isatty"
	"github.com/spf13/viper"
	"github.com/tidwall/pretty"
)

func connectSelectedCore() (coreclient.RestClient, error) {
	selected := viper.GetString("cores.selected")
	if len(globalFlagCore) != 0 {
		selected = globalFlagCore
	}

	return connectCore(selected)
}

func connectCore(name string) (coreclient.RestClient, error) {
	list := viper.GetStringMapString("cores.list")

	core, ok := list[name]
	if !ok {
		return nil, fmt.Errorf("selected core doesn't exist")
	}

	u, err := url.Parse(core)
	if err != nil {
		return nil, fmt.Errorf("invalid data for core: %w", err)
	}

	address := u.Scheme + "://" + u.Host + u.Path
	password, _ := u.User.Password()

	client, err := coreclient.New(coreclient.Config{
		Address:      address,
		Username:     u.User.Username(),
		Password:     password,
		AccessToken:  u.Query().Get("accessToken"),
		RefreshToken: u.Query().Get("refreshToken"),
	})
	if err != nil {
		return nil, fmt.Errorf("can't connect to core at %s: %w", address, err)
	}

	about, err := client.About(true)
	if err != nil {
		return nil, fmt.Errorf("can't fetch details from core at %s: %w", address, err)
	}

	version := about.Version.Number
	corename := about.Name
	coreid := about.ID
	accessToken, refreshToken := client.Tokens()

	query := u.Query()

	query.Set("accessToken", accessToken)
	query.Set("refreshToken", refreshToken)
	query.Set("version", version)
	query.Set("name", corename)
	query.Set("id", coreid)

	u.RawQuery = query.Encode()

	list[name] = u.String()

	viper.Set("cores.list", list)
	viper.WriteConfig()

	//fmt.Fprintln(os.Stderr, client.String())

	return client, nil
}

func getEditor() (string, string, error) {
	editor := viper.GetString("editor")
	if len(editor) == 0 {
		editor = os.Getenv("EDITOR")
	}

	if len(editor) == 0 {
		return "", "", fmt.Errorf("no editor defined")
	}

	path, err := exec.LookPath(editor)
	if err != nil {
		if !errors.Is(err, exec.ErrDot) {
			return "", "", fmt.Errorf("%s: %w", editor, err)
		}
	}

	return editor, path, nil
}

func editData(data []byte) ([]byte, bool, error) {
	editor, _, err := getEditor()
	if err != nil {
		return nil, false, err
	}

	file, err := os.CreateTemp("", "corecli_*")
	if err != nil {
		return nil, false, err
	}

	filename := file.Name()

	defer os.Remove(filename)

	_, err = file.Write(data)
	file.Close()

	if err != nil {
		return nil, false, err
	}

	for {
		editor := exec.Command(editor, filename)
		editor.Stdout = os.Stdout
		editor.Stderr = os.Stderr
		editor.Stdin = os.Stdin
		if err := editor.Run(); err != nil {
			return nil, false, err
		}

		editedData, err := os.ReadFile(filename)
		if err != nil {
			return nil, false, err
		}

		var x interface{}

		if err := json.Unmarshal(editedData, &x); err != nil {
			errorData, err := formatJSONError(editedData, err)
			fmt.Printf("%s\n", errorData)
			fmt.Printf("Invalid JSON: %s\n", err)
			fmt.Printf("Do you want to re-open the editor (Y/n)? ")

			var char rune
			if _, err := fmt.Scanf("%c", &char); err != nil {
				return nil, false, err
			}

			if char == '\n' || char == 'Y' || char == 'y' {
				continue
			}

			return nil, false, fmt.Errorf("invalid JSON: %w", err)
		}

		return editedData, !bytes.Equal(data, editedData), nil
	}
}

func formatJSONError(input []byte, err error) ([]byte, error) {
	if jsonError, ok := err.(*json.SyntaxError); ok {
		line, character, offsetError := lineAndCharacter(input, int(jsonError.Offset))
		if offsetError != nil {
			return input, err
		}

		return markJSONError(input, line-1, character-1), fmt.Errorf("syntax error at line %d, character %d: %w", line, character, err)
	}

	if jsonError, ok := err.(*json.UnmarshalTypeError); ok {
		line, character, offsetError := lineAndCharacter(input, int(jsonError.Offset))
		if offsetError != nil {
			return input, err
		}

		return markJSONError(input, line-1, character-1), fmt.Errorf("expect type '%s' for '%s' at line %d, character %d: %w", jsonError.Type.String(), jsonError.Field, line, character, err)
	}

	return input, err
}

func lineAndCharacter(input []byte, offset int) (line int, character int, err error) {
	lf := byte(0x0A)

	if offset > len(input) || offset < 0 {
		return 0, 0, fmt.Errorf("couldn't find offset %d within the input", offset)
	}

	// Humans tend to count from 1.
	line = 1
	lastLineCharacters := 0

	for i, b := range input {
		if b == lf {
			line++
			lastLineCharacters = character
			character = 0
		}
		character++
		if i == offset {
			break
		}
	}

	// Fix the reported offset because it reflects the consumed bytes from
	// parsing and not the actual position of the error.
	if line == 1 {
		character -= 1
	} else {
		character -= 2
		if character < 0 {
			line -= 1
			character = lastLineCharacters
		}
	}

	return line, character, nil
}

func markJSONError(input []byte, line, character int) []byte {
	lf := byte(0x0A)
	output := bytes.Buffer{}
	lineContext := 10

	lines := bytes.Split(input, []byte{lf})

	nlines := len(lines)
	fromLine := line - lineContext
	fromCut := true
	if fromLine < 0 {
		fromLine = 0
		fromCut = false
	}
	toLine := line + lineContext
	toCut := true
	if toLine >= nlines {
		toLine = nlines - 1
		toCut = false
	}

	if fromCut {
		output.Write([]byte(fmt.Sprintf("... %d previous lines omitted ...\n", fromLine)))
	}

	for i := fromLine; i < toLine; i++ {
		l := lines[i]

		output.Write(l)
		output.WriteByte(lf)

		if i == line {
			m := make([]byte, character+1)
			for i := range m {
				m[i] = '_'
			}
			m[character] = '^'

			output.Write(m)
			output.WriteByte(lf)
		}
	}

	if toCut {
		output.Write([]byte(fmt.Sprintf("... %d following lines omitted ...\n", nlines-toLine)))
	}

	return output.Bytes()
}

func formatJSON(d interface{}, useColor bool) (string, error) {
	data, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	data = pretty.PrettyOptions(data, &pretty.Options{
		Width:    pretty.DefaultOptions.Width,
		Prefix:   pretty.DefaultOptions.Prefix,
		Indent:   pretty.DefaultOptions.Indent,
		SortKeys: true,
	})

	if !useColor {
		return string(data), nil
	}

	data = pretty.Color(data, nil)

	return string(data), nil
}

func writeJSON(w io.Writer, d interface{}, useColor bool) error {
	color := useColor

	if color {
		if w, ok := w.(*os.File); ok {
			if !isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd()) {
				color = false
			}
		} else {
			color = false
		}
	}

	if len(globalFlagJq) != 0 {
		query, err := gojq.Parse(globalFlagJq)
		if err != nil {
			return err
		}

		data, err := json.Marshal(d)
		if err != nil {
			return err
		}

		var input any

		err = json.Unmarshal(data, &input)
		if err != nil {
			return err
		}

		iter := query.Run(input)
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}
			if err, ok := v.(error); ok {
				return err
			}
			switch x := v.(type) {
			case string:
				fmt.Printf("%s\n", x)
			default:
				j, _ := formatJSON(v, color)
				fmt.Printf("%s", j)
			}
		}

		return nil
	}

	data, err := formatJSON(d, color)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, data)

	return nil
}

func formatByteCountBinary(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d  B", b)
	}

	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func processTable(list []coreclientapi.Process, processMap map[string]string, nodes map[string]coreclientapi.ClusterNode) {
	t := table.NewWriter()

	t.AppendHeader(table.Row{"ID", "Domain", "Reference", "Order", "State", "Memory", "CPU", "Runtime", "Node", "Last Log"})

	nodeCount := map[string]uint64{}
	stateCount := map[string]uint64{}
	deployCount := map[string]uint64{}

	for _, p := range list {
		runtime := p.State.Runtime
		if p.State.State != "running" {
			runtime = 0

			if p.State.Reconnect > 0 {
				runtime = -p.State.Reconnect
			}
		}

		order := strings.ToUpper(p.State.Order)
		switch order {
		case "START":
			order = text.FgGreen.Sprint(order)
		case "STOP":
			order = text.Colors{text.FgWhite, text.Faint}.Sprint(order)
		}

		state := strings.ToUpper(p.State.State)
		switch state {
		case "RUNNING":
			state = text.FgGreen.Sprint(state)
		case "FINISHED":
			state = text.Colors{text.FgWhite, text.Faint}.Sprint(state)
		case "DEPLOYING":
			state = text.Colors{text.FgWhite, text.Faint}.Sprint(state)
		case "FAILED":
			state = text.FgRed.Sprint(state)
		case "STARTING":
			state = text.FgCyan.Sprint(state)
		case "FINISHING":
			state = text.FgCyan.Sprint(state)
		case "KILLED":
			state = text.Colors{text.FgRed, text.Faint}.Sprint(state)
		}

		stateCount[state]++

		nodeid := processMap[coreclient.NewProcessID(p.ID, p.Domain).String()]
		if nodeid != p.CoreID {
			nodeid = "(" + nodeid + ")"

			if len(p.CoreID) != 0 {
				nodeid = p.CoreID + " " + nodeid
				nodeid = text.FgYellow.Sprint(nodeid)
				deployCount[text.FgYellow.Sprint("MISDEPLOYED")]++
			} else {
				nodeid = text.FgRed.Sprint(nodeid)
				deployCount[text.FgRed.Sprint("UNDEPLOYED")]++
			}
		} else {
			deployCount[text.FgGreen.Sprint("DEPLOYED")]++
		}

		nodeCount[p.CoreID]++

		cpu := fmt.Sprintf("%.1f%%", p.State.Resources.CPU.Current)
		if p.State.Resources.CPU.IsThrottling {
			cpu = "* " + cpu
			if nodes != nil {
				if n, ok := nodes[p.CoreID]; ok {
					n.Resources.IsThrottling = true
					nodes[p.CoreID] = n
				}
			}
		}

		lastlog := p.State.LastLog
		if len(lastlog) > 58 {
			lastlog = lastlog[:55] + "..."
		}

		t.AppendRow(table.Row{
			p.ID,
			p.Domain,
			p.Reference,
			order,
			state,
			formatByteCountBinary(p.State.Resources.Memory.Current),
			cpu,
			(time.Duration(runtime) * time.Second).String(),
			nodeid,
			lastlog,
		})
	}

	t.SetAutoIndex(true)

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 4, Align: text.AlignRight},
		{Number: 5, Align: text.AlignRight},
		{Number: 6, Align: text.AlignRight},
		{Number: 7, Align: text.AlignRight},
		{Number: 8, Align: text.AlignRight},
	})

	t.SortBy([]table.SortBy{
		{Number: 2, Mode: table.Asc},
		{Number: 1, Mode: table.Asc},
		{Number: 4, Mode: table.Asc},
		{Number: 6, Mode: table.Dsc},
	})

	t.SetStyle(table.StyleLight)

	fmt.Println(t.Render())

	// print process state count

	t = table.NewWriter()

	t.AppendHeader(table.Row{"State", "Count"})

	sum := uint64(0)
	for state, count := range stateCount {
		t.AppendRow(table.Row{
			state,
			count,
		})
		sum += count
	}

	t.AppendFooter(table.Row{
		"",
		sum,
	})

	t.SortBy([]table.SortBy{
		{Number: 1, Mode: table.Asc},
	})

	t.SetStyle(table.StyleLight)

	fmt.Println(t.Render())

	// print deployment count

	t = table.NewWriter()

	t.AppendHeader(table.Row{"Deployment", "Count"})

	sum = uint64(0)
	for state, count := range deployCount {
		t.AppendRow(table.Row{
			state,
			count,
		})
		sum += count
	}

	t.AppendFooter(table.Row{
		"",
		sum,
	})

	t.SetStyle(table.StyleLight)

	fmt.Println(t.Render())

	// print node count

	t = table.NewWriter()

	t.AppendHeader(table.Row{"Node", "Count", "CPU", "Memory", "Throttling"})

	sum = uint64(0)
	for nodeid, count := range nodeCount {
		cpu := "n/a"
		memory := "n/a"
		throttling := false

		if n, ok := nodes[nodeid]; ok {
			if n.Resources.CPULimit != 0 {
				cpu = fmt.Sprintf("%.1f%%", n.Resources.CPU/n.Resources.CPULimit*100)
			}

			if n.Resources.MemLimit != 0 {
				memory = fmt.Sprintf("%.1f%%", float64(n.Resources.Mem)/float64(n.Resources.MemLimit)*100)
			}

			throttling = n.Resources.IsThrottling
		}

		t.AppendRow(table.Row{
			nodeid,
			count,
			cpu,
			memory,
			throttling,
		})
		sum += count
	}

	t.AppendFooter(table.Row{
		"",
		sum,
	})

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 2, Align: text.AlignRight},
		{Number: 3, Align: text.AlignRight},
		{Number: 4, Align: text.AlignRight},
		{Number: 5, Align: text.AlignRight},
	})

	t.SortBy([]table.SortBy{
		{Number: 1, Mode: table.Asc},
	})

	t.SetStyle(table.StyleLight)

	fmt.Println(t.Render())
}

func dbProcessTable(list []coreclientapi.Process, processMap map[string]string) {
	t := table.NewWriter()

	t.AppendHeader(table.Row{"ID", "Domain", "Reference", "Order", "State", "Memory LMT", "CPU LMT", "Node", "Error"})

	for _, p := range list {
		order := strings.ToUpper(p.State.Order)
		switch order {
		case "START":
			order = text.FgGreen.Sprint(order)
		case "STOP":
			order = text.Colors{text.FgWhite, text.Faint}.Sprint(order)
		}

		state := "DEPLOYED"
		if p.State.State == "failed" {
			state = "FAILED"
		}

		switch state {
		case "DEPLOYED":
			state = text.FgGreen.Sprint(state)
		default:
			state = text.FgRed.Sprint(state)
		}

		nodeid := processMap[coreclient.NewProcessID(p.ID, p.Domain).String()]

		lastlog := p.State.LastLog
		if len(lastlog) > 58 {
			lastlog = lastlog[:55] + "..."
		}

		t.AppendRow(table.Row{
			p.ID,
			p.Domain,
			p.Reference,
			order,
			state,
			formatByteCountBinary(p.State.Resources.Memory.Limit),
			fmt.Sprintf("%.1f%%", p.State.Resources.CPU.Limit),
			nodeid,
			lastlog,
		})
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 2, Align: text.AlignRight},
		{Number: 4, Align: text.AlignRight},
		{Number: 5, Align: text.AlignRight},
		{Number: 6, Align: text.AlignRight},
		{Number: 7, Align: text.AlignRight},
	})

	t.SortBy([]table.SortBy{
		{Number: 2, Mode: table.Asc},
		{Number: 1, Mode: table.Asc},
		{Number: 4, Mode: table.Asc},
	})

	t.SetStyle(table.StyleLight)

	fmt.Println(t.Render())
}

func processIO(p coreclientapi.Process) {
	if p.State == nil || p.State.Progress == nil {
		return
	}

	if len(p.State.Progress.Input) == 0 && len(p.State.Progress.Output) == 0 {
		return
	}

	t := table.NewWriter()

	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	t.SetTitle("Inputs / Outputs")
	t.AppendHeader(table.Row{"", "#", "ID", "Type", "URL", "Specs"}, rowConfigAutoMerge)

	for i, p := range p.State.Progress.Input {
		var specs string
		if p.Type == "audio" {
			specs = fmt.Sprintf("%s %s %dHz", strings.ToUpper(p.Codec), p.Layout, p.Sampling)
		} else {
			specs = fmt.Sprintf("%s %dx%d", strings.ToUpper(p.Codec), p.Width, p.Height)
		}

		t.AppendRow(table.Row{
			"input",
			i,
			p.ID,
			strings.ToUpper(p.Type),
			p.Address,
			specs,
		}, rowConfigAutoMerge)
	}

	for i, p := range p.State.Progress.Output {
		var specs string
		if p.Type == "audio" {
			specs = fmt.Sprintf("%s %s %dHz", strings.ToUpper(p.Codec), p.Layout, p.Sampling)
		} else {
			specs = fmt.Sprintf("%s %dx%d", strings.ToUpper(p.Codec), p.Width, p.Height)
		}

		t.AppendRow(table.Row{
			"output",
			i,
			p.ID,
			strings.ToUpper(p.Type),
			p.Address,
			specs,
		}, rowConfigAutoMerge)
	}

	t.SetStyle(table.StyleLight)

	fmt.Println(t.Render())
}

const (
	CharsetLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharsetNumbers = "1234567890"
	CharsetSymbols = "#@+*%&/<>[]()=?!$.,:;-_"

	CharsetAll = CharsetLetters + CharsetNumbers + CharsetSymbols
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func StringLetters(length int) string {
	return StringWithCharset(length, CharsetLetters)
}

func StringNumbers(length int) string {
	return StringWithCharset(length, CharsetNumbers)
}

func StringAlphanumeric(length int) string {
	return StringWithCharset(length, CharsetLetters+CharsetNumbers)
}

func String(length int) string {
	return StringWithCharset(length, CharsetAll)
}

func loadTemplate(name string) (coreclientapi.ProcessConfig, error) {
	config := coreclientapi.ProcessConfig{}
	templates := viper.GetStringMap("templates")

	data, hasTemplate := templates[name]
	if !hasTemplate {
		return config, fmt.Errorf("the template with name '%s' doesn't exist", name)
	}

	xxx, err := json.Marshal(data)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(xxx, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
