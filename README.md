# Launchy

Because someone needs to launch the Mac OS X services already.

## Introduction

I hate OS'X's `launchctl`? You have to give it exact filenames. The syntax is annoyingly different from Linux's nice, simple init system and it's overly verbose. It's just not a very developer-friendly tool.

Launchy aims to be that friendly and performant tool by wrapping `launchctl` and providing a few simple operations that you perform all the time:

```
launchy list [-r pattern]
launchy status [-r pattern]
```

where pattern is just a substring that matches the agent's plist filename. Use launchy's `list` or `status` command to view possible completions.

### Examples

#### List

```shell
~/go/src/github.com/kkirsche/launchy git:master ❯❯❯ launchy list                               ✱ ◼
at.obdev.LittleSnitchUIAgent
org.macosforge.xquartz.startx
com.adobe.AAM.Updater-1.0
com.apple.serveralertproxy
homebrew.mxcl.postgresql
org.virtualbox.vboxwebsrv
```

```shell
~/go/src/github.com/kkirsche/launchy git:master ❯❯❯ launchy list -r "LittleSnitch"             ✱ ◼
at.obdev.LittleSnitchUIAgent
```


### Credit
Shamelessly ripped off of https://github.com/eddiezane/lunchy
