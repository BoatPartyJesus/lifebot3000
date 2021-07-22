package naturallanguageprocessing

import (
	"fmt"
	"regexp"

	"github.com/jdkato/prose/v2"
)

type INlp interface {
	IntentMap(input string) NLPResult
}

func IntentMap(input string) NLPResult {

	reducedInput := aliasReduce(input)

	doc, _ := prose.NewDocument(reducedInput)

	userIdDetection, _ := regexp.Compile(`<@\w{11}>\W*`)

	var mentionAction string
	mentionTarget := ""

	for _, tok := range doc.Tokens() {
		//fmt.Println(tok.Text, tok.Tag, tok.Label)
		if tok.Tag[:1] == "V" {
			fmt.Println("Action:", tok.Text)
			mentionAction = tok.Text
		}

		if userIdDetection.MatchString(tok.Text) {
			cleanUp, _ := regexp.Compile(`\W`)
			uncleanUser := tok.Text
			mentionTarget = cleanUp.ReplaceAllLiteralString(uncleanUser, "")
		}

		if tok.Tag[:1] == "J" && mentionTarget == "" {
			if tok.Text == "recent" {
				mentionTarget = tok.Text
			}
		}
	}

	return NLPResult{
		mentionAction,
		mentionTarget,
	}
}

type NLPResult struct {
	Action string
	Target string
}

func aliasReduce(text string) string {
	aliasMap := map[string][]string{
		"":       {`\ba\b`, `\bthe\b`, `\bto\b`, `\bof\b`, `\bon\b`, `\boff\b`, `'.*\b`, "please", `users?`, "list"},
		"recent": {`recent.*list`, `recent.*names`, "recents"},
		"me":     {"my name", "myself"},
		"give":   {"display", "show", `\bis\b`, "gimme"},
		"add":    {"include", "put"},
		"remove": {"take", "delete", "forget"},
		"reset":  {"clear", "empty", "blat", "nuke"},
	}

	for target, aliases := range aliasMap {
		buildTargets := ""
		for ind, tgt := range aliases {
			divider := ""

			if ind > 0 {
				divider = "|"
			}

			buildTargets += fmt.Sprintf("%s%s", divider, tgt)
		}

		targetReplacement := regexp.MustCompile(buildTargets)

		text = targetReplacement.ReplaceAllLiteralString(text, target)
	}

	return text
}
