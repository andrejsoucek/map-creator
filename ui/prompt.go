package ui

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func YesNoPrompt(label string, def bool) bool {
    choices := "Y/n"
    if !def {
        choices = "y/N"
    }

    r := bufio.NewReader(os.Stdin)
    var s string

    for {
        fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
        s, _ = r.ReadString('\n')
        s = strings.TrimSpace(s)
        if s == "" {
            return def
        }
        s = strings.ToLower(s)
        if s == "y" || s == "yes" {
            return true
        }
        if s == "n" || s == "no" {
            return false
        }
    }
}