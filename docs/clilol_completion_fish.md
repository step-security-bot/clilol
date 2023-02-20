---
title: "clilol completion fish"
---
## clilol completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	clilol completion fish | source

To load completions for every new session, execute once:

	clilol completion fish > ~/.config/fish/completions/clilol.fish

You will need to start a new shell for this setup to take effect.


```
clilol completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -j, --json     output json
  -s, --silent   be silent
```

### SEE ALSO

* [clilol completion](clilol_completion.md)	 - Generate the autocompletion script for the specified shell
