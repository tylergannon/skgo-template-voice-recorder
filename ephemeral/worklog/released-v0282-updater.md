# Released v0.28.2 updater qualification

decision: Coordinator authorized one ordinary released updater attempt on one exact fresh copy of the original held edited v0.27.0 consumer. Copy all inputs without source/manifest patches; preserve baseline and prior authored receipt hashes. No development completion flags or host identity overrides. Failure stops the consumer patch chase. Native runtime and managed installer #299 stay separate.

finding: Released v0.28.2 installed and reached staged dependency alignment, pnpm 12.9.1 installation and Kit verification. pnpm reported adding a fresh adapter release to minimumReleaseAgeExclude, then the authored-config guard refused that key. All 120 input entries in both copy and held baseline are unchanged. The app stays v0.27.0; completion/build/runtime were not reached. Stop here rather than guessing another prerequisite patch. Framework owner and coordinator received the exact failure before publication.
