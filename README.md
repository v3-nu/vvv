# clycli
Command-line utility for various actions

## Commands
- clycli (alias: cly, cy, cc)
  - kubectl (alias: k)
    - getcontexts (alias: getc)
    - setcontext (alias: setc, con)
    - setnamespace (alias: setns, setn, ns)
    - patchfinalizers (alias: patchfin)
  - install
    - kubectl


## TODO:
- Change alias functionality
  - Aliasrc: (ex)
    alias clycli="$(which clycli) alias exec [aliasname] [args]"
  - Run alias with utils.AliasCommand with %s available
