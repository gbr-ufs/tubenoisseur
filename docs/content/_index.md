---
description: YouTuber Reference Scraper.
linkTitle: Home
title: tubenoisseur
type: docs
---

![Relaxed Go Gopher floating in a giant wineglass](gopher.png)

# [tubenoisseur: YouTuber Reference Scraper](https://gbr-ufs.github.io/tubenoisseur/)

[![Codecov](https://codecov.io/gh/gbr-ufs/tubenoisseur/graph/badge.svg?token=RZBZBQG5PN)](https://codecov.io/gh/gbr-ufs/tubenoisseur)
[![License](https://img.shields.io/badge/license_-GPLv3+-822422?logo=GNU&logoColor=black&labelColor=white)](LICENSE)

![GIF showcasing the program being used on the "bashbunni" YouTube channel, and the generated file with the sources](https://vhs.charm.sh/vhs-2i5ntpn2Sm0nBXVdDMJ9fW.gif)

<p align="center">
  <a href="https://vhs.charm.sh">
    <img alt="VHS" src="https://stuff.charm.sh/vhs/badge.svg">
  </a>
</p>

## Introduction

tubenoisseur is a command-line interface (CLI) built with [Cobra](https://cobra.dev/) meant to be coupled with a scheduling system, such as systemd or a cron job, to gather links in the description of videos from a YouTube channel, with the goal of mapping out their sources.

## Philosophies

- Dependency Injection.
- Functional Core, Imperative Shell[^1].
- Parse, don't validate[^2].
- Test Driven Development[^3].
- You Aren't Gonna Need It[^4].

[^1]: Bernhardt, G. (2012) Functional core, imperative shell. Destroy All Software. Available at: https://www.destroyallsoftware.com/screencasts/catalog/functional-core-imperative-shell (Accessed: May 1, 2026).

[^2]: King, A. (2019) Parse, don’t validate. Alexis King’s Blog. Available at: https://lexi-lambda.github.io/blog/2019/11/05/parse-don-t-validate/ (Accessed: September 29, 2025).

[^3]: Beck, K. (2003) Test-driven development: By example. Boston: Addison-Wesley (The Addison-Wesley signature series).

[^4]: Jeffries, R., Anderson, A. and Hendrickson, C. (2001) Extreme programming installed. Boston: Addison-Wesley (The XP series).
