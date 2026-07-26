# SPDX-FileCopyrightText: 2026 Gabriel Santos de Souza <gabriel.santosdesouza@dcomp.ufs.br>
#
# SPDX-License-Identifier: GPL-3.0-or-later

Require bat
Require tubenoisseur

Set FontSize 22
Set Height 650
Set Theme "Github"
Set Width 1300

Type "tubenoisseur bashbunni --debug"
Enter
Wait /\$/
Type "bat tubenoisseur/bashbunni/bashbunni.md"
Enter
Sleep 10s

Output demo.gif
