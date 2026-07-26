# SPDX-FileCopyrightText: 2026 Gabriel Santos de Souza <gabriel.santosdesouza@dcomp.ufs.br>
#
# SPDX-License-Identifier: GPL-3.0-or-later

FROM scratch
COPY tubenoisseur /usr/bin/tubenoisseur
ENTRYPOINT ["/usr/bin/tubenoisseur"]
