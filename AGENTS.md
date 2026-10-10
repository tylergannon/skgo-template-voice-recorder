# SKGo voice recorder template

Read VISION.md and README.md first. Work in an isolated Git worktree; keep
actionable operating notes under ephemeral/worklog/ and checkpoint completed
changes. Squash merge reviewed PRs and synchronize the root checkout afterward.

`recorder/` owns canonical reusable application sources. Package an explicit
reviewed allowlist from those files, not a second hand-maintained template tree.
The original Voice Notes application is a separate working baseline: do not
modify or replace it as part of extraction.

Use the actual shared SKGo templateapi and native-app foundation. Do not invent
a local SDK, generic native host or competing project bootstrap. Parameterize
consumer names/imports/bundle identities; exclude developer IDs, credentials,
toolchains, build products, absolute checkout paths and validation artifacts.
Keep required third-party licenses. Application-specific native bridge messages
and background recorder behavior belong here; host recovery belongs to SKGo.

Normal regression is synthetic, without microphone permissions or system route
changes. Real adapter checks are separate and deadline bounded. Use normal
computer controls for screenshots. Never work around broken computer use when
it is needed: repair it or report the blocker. Stop and explain after two
unchanged failures or five minutes without progress; repair the setup rather
than repeat an uncertain end-to-end test.

Delegate independent extraction and review with explicit file ownership; account
for others' edits and do not revert them. Preserve the user-requested independent
Claude Opus 5.5 validation for native behavior. Distinguish compile/tests from
running consumer evidence, and original-app evidence from fresh-consumer proof.

Record shared-SKGo gaps in ephemeral/skgo-feedback.md with version, expectation,
observation and reproduction. Coordinate directly with the authorized framework
peer and coordinating chat; do not create duplicate monitors or delegation layers.
