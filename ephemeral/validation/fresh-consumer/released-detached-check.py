import os,time,pathlib,subprocess,shutil,json
source=pathlib.Path('/Users/tyler/.codex/worktrees/recorder-template/skgo-template-voice-recorder')
hidden=source.with_name(source.name+'-unavailable')
base=pathlib.Path('/private/tmp/skgo-recorder-consumer.GthvK4')
consumer=base/'released'
runtime=base/'released-runtime';runtime.mkdir(exist_ok=True)
tools=base/'released-tools'
assert not hidden.exists() and not tools.exists()
shutil.copytree(consumer/'native/build/macos-synthetic/DerivedData/Build/Products/Debug/FieldNotes.app',runtime/'FieldNotes.app',dirs_exist_ok=True)
server=None
try:
 (source/'build/toolchain').rename(tools)
 source.rename(hidden)
 assert not source.exists()
 env=os.environ.copy();env.update(GOWORK='off',CI='1',PATH=str(tools/'node_modules/.bin')+':'+env['PATH'])
 with (runtime/'build.log').open('w') as log:
  log.write('Template checkout and all its plugin directories unavailable at original path. Toolchain moved to independent temporary directory.\n');log.flush()
  subprocess.run(['go','build','-o',str(runtime/'FieldNotes-server'),'./cmd'],cwd=consumer,env=env,stdout=log,stderr=log,check=True,timeout=30)
  log.write('Standalone server rebuild succeeded.\n');log.flush()
  subprocess.run(['go','tool','skgo','native','build','--platform','macos','--preset','synthetic'],cwd=consumer,env=env,stdout=log,stderr=log,check=True,timeout=60)
  log.write('Native macOS rebuild succeeded with original template/plugin path absent.\n');log.flush()
 shutil.rmtree(runtime/'FieldNotes.app')
 shutil.copytree(consumer/'native/build/macos-synthetic/DerivedData/Build/Products/Debug/FieldNotes.app',runtime/'FieldNotes.app')
 serverlog=(runtime/'server.log').open('w')
 server=subprocess.Popen([str(runtime/'FieldNotes-server')],cwd=runtime,env=dict(env,ADDR='127.0.0.1:8188'),stdout=serverlog,stderr=serverlog)
 (runtime/'ready').write_text('ready')
 print('Detached release rebuilds passed; server and copied native product ready. Source remains unavailable for up to90 seconds.',flush=True)
 deadline=time.time()+90
 while time.time()<deadline and not (runtime/'done').exists():time.sleep(0.25)
finally:
 if server is not None:server.terminate();server.wait(timeout=10)
 if hidden.exists():hidden.rename(source)
 if tools.exists():tools.rename(source/'build/toolchain')
 print('Restored source/plugin/toolchain paths and stopped detached server.',flush=True)
