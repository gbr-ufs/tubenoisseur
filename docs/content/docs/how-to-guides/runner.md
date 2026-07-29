---
title: tubenoisseur Runner
type: docs
---

# tubenoisseur Runner

tubenoisseur was made with scheduled runners in mind. You're supposed to throw it at something like a cron job or a systemd service so it runs automatically. Besides these options, another recommended one is to set up a runner on GitHub by using [GitHub Actions](https://github.com/features/actions). Note that, if you don't want to set this up, but are still interested, there's a [community runner available](../references/tubenoisseur_runner).

Without further ado, let's get into how to set this thing up.


## Steps

1. Create a GitHub repository.
2. Create a git repository.
3. Under `.github/workflows/tubenoisseur.yaml`, write this:
```yaml
name: Tubenoisseur

permissions:
  contents: write

on:
  schedule:
    # CHANGE THIS (<https://crontab.guru/>. Set it to a random time!)
    - cron: 0 10 * * *
  workflow_dispatch:
jobs:
  tubenoisseur:
    runs-on: ubuntu-latest
    strategy:
      max-parallel: 1
      matrix:
        channel:
          # CHANGE THIS (list of channels)
          - BrodieRobertson
          - TheLinuxEXP
    steps:
    - name: Checkout Repository
      uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # Remember to update this!
    - name: Cache
      id: cache
      uses: actions/cache@640a1c2554105b57832a23eea0b4672fc7a790d5 # Remember to update this!
      with:
        key: ${{ runner.os }}-go-tubenoisseur-v1.0.4 # Remember to update this!
        path: ~/go/bin/tubenoisseur
    - name: Add Go Binary Directory to $PATH
      if: ${{ steps.cache.outputs.cache-hit == 'true' }}
      run: echo "$HOME/go/bin" >> $GITHUB_PATH
    - name: Set Up Go
      if: ${{ steps.cache.outputs.cache-hit != 'true' }}
      uses: actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c # Remember to update this!
      with:
        go-version: 1.26.5
    - name: Install Tubenoisseur
      if: ${{ steps.cache.outputs.cache-hit != 'true' }}
      run: go install github.com/gbr-ufs/tubenoisseur@v1.0.4 # Remember to update this!
    - name: Tubenoisseur
      run: tubenoisseur ${{ matrix.channel }} --debug
      env:
        CLICOLOR_FORCE: "1"
    - name: Git Pull
      run: git pull origin ${{ github.ref_name }} --autostash --rebase
    - name: Commit
      uses: stefanzweifel/git-auto-commit-action@04702edda442b2e678b25b537cec683a1493fcb9 # Remember to update this!
      with:
        commit_message: "chore: update ${{ matrix.channel }}"
```
Modify the lines marked as `CHANGE THIS`. The first one is the actual schedule. The second one is the list of channels you want to use tubenoisseur on.
4. Push the changes to your GitHub repository.
