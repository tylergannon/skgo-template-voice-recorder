# Source packaging

`allowlist.json` names each canonical source and its generated consumer path.
Prepare the embedded payload before building or checking the plugin:

```sh
go run ./template/stage
go test ./...
```

The stage copies only named sources, keeps a `.txt` suffix so Go does not compile
application payload files as plugin packages, and is ignored by Git. The renderer
removes that suffix and derives imports, display name and probe identity from the
destination project. No application source is maintained under the payload tree.

The SKGo peer owns the host-compatible plugin builder. Use that builder against
the actual standalone or project-tool invocation; matching version text alone
does not establish Go plugin compatibility. Exact invocation instructions await
the shared helper handoff. The native application slot/config handoff must be
integrated before this package exports a usable plugin.
