---
title: Getting Started
type: docs
---

# Getting Started

## Installation

tubenoisseur is a [Go](https://go.dev/) project. That means you can install the program directly from the language!

Option A: Via Go

```shell
go install github.com/gbr-ufs/tubenoisseur@latest
```

Option B: Standalone Binaries

Precompiled binaries are available at: https://github.com/gbr-ufs/tubenoisseur/releases/latest

Option C: System Package

Linux users can install tubenoisseur from the comfort of their package manager. Just check if your package manager or distro is here:

[![packaging status](https://repology.org/badge/vertical-allrepos/tubenoisseur.svg)](https://repology.org/project/tubenoisseur/versions)

## Your First Report

Let's try tubenoisseur out by checking out what [The Linux Experiment](https://www.youtube.com/@TheLinuxEXP) is using as his sources.

Open your terminal and run the `tubenoisseur` command. It takes the channel handle (without the "@", as having to type it out every time would be annoying) as its argument:

```shell
tubenoisseur TheLinuxEXP
```

![GIF showing a terminal with the user typing the previous command.](https://vhs.charm.sh/vhs-1F3pkXduOR4eUwyPIAJnjb.gif)

Let's check it out:

![GIF showing the generated report. It is a Markdown file with the channel handle as the first-level header, and under it is the top sources of the channel.](https://vhs.charm.sh/vhs-2y3rGQkYbpn0oxu70MDnt2.gif)

Hm, Tuxedo Computers? PayPal? Liberapay? That seems like a lot of false positives. I think we can do better. Let's try rerunning with the `--exclude` flag and with the `--debug` flag to get some actual output:

```shell
tubenoisseur TheLinuxEXP --debug --exclude squarespace.com,www.tuxedocomputers.com,www.youtube.com,www.patreon.com,paypal.me,liberapay.com,the-linux-experiment.creator-spring.com
```

![GIF showing a terminal with the user typing the previous command. Thanks to the `--debug` flag, the program actually has some output now.](https://vhs.charm.sh/vhs-6RYUSYZ5wWo5hvvVF3JJus.gif)

Let's see it again:

![GIF showing the generated report. The false positives are now gone.](https://vhs.charm.sh/vhs-1keZib0IVYwbOsNvFFGX2m.gif)

Much better. Now, you might be wandering: "Wow, that really was a lot to type. Do I have to do it every time?". Lucky for you, the answer is no! tubenoisseur automatically saves the excluded domains, in a per channel JSON file, used for state management. Let's have a look at it:

![GIF showcasing the domain file. The "channelHandle" and "domainCounts" fields are visible.](https://vhs.charm.sh/vhs-3IvCeNTgPNTVnPcR5ieVUe.gif)
