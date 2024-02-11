#!/bin/bash
# Gendo replicates common code across Git repos.
# Copyright (C) 2022  Ken Shibata <+@nyiyui.ca>
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

if [ "${BASH_SOURCE[0]}" -ef "$0" ]
then
    echo "This script should be sourced, not run."
    exit 1
fi

apply() {
	remote=$1
	config=$2
	commit_message=$3
	1>&2 echo "=== $remote"
	tmpdir=$(mktemp -d --suffix=gen)
	(
		git clone --quiet "$remote" $tmpdir
		./replace -config "$config" -dir $tmpdir
		cd $tmpdir
    cat ./layouts/_default.html
		git add .
		git commit -m "$commit_message"
		git push --quiet
	)
	rm -rf $tmpdir
}
