#!/bin/bash
# Collects the Debian packages of a distroless image that need refreshing from the Debian image this runs in.
#
# Usage: sync-packages.sh <base-rootfs> <delta-dir> [extra-package ...]
#
# Every package recorded in <base-rootfs>/var/lib/dpkg/status.d, plus the extra packages, is upgraded here
# to the latest version in the Debian archive. A package is written to <delta-dir> only when the version
# here is strictly newer than the one recorded in <base-rootfs>. Extra packages are also written when the
# versions are equal, which restores the files that distroless strips from a package (e.g. ca-certificates).
# For each package written, <delta-dir> gets every file of the package except man pages, plus its status.d
# metadata, so that image scanners report the version that is actually shipped.
# <delta-dir> is meant to be copied over the root of the distroless image, so paths that would conflict
# with it are skipped: paths whose type (directory, symlink, file) differs from <base-rootfs>, and files
# that distroless customized (e.g. os-release), detected by not matching the checksums recorded in
# <base-rootfs>.
set -euo pipefail

base=$1
delta=$2
shift 2
extras=" $* "
base_status_dir=$base/var/lib/dpkg/status.d

mapfile -t packages < <( (awk '/^Package:/{print $2}' "$base_status_dir"/*; printf '%s\n' "$@") | sort -u)

apt-get update
apt-get upgrade -y
apt-get install -y --no-install-recommends "${packages[@]}"

# Prints the type of a path without following a final symlink
path_type() {
  if [[ -L $1 ]]; then echo symlink
  elif [[ -d $1 ]]; then echo directory
  else echo file
  fi
}

mkdir -p "$delta/var/lib/dpkg/status.d"
work=$(mktemp -d)
updated=()
for pkg in "${packages[@]}"; do
  new_version=$(dpkg-query -W -f='${Version}' "$pkg")
  base_status_file=$(grep -lx "Package: $pkg" "$base_status_dir"/* | grep -v '\.md5sums$' | head -n 1 || true)
  status_file=var/lib/dpkg/status.d/$pkg

  if [[ -n $base_status_file ]]; then
    status_file=var/lib/dpkg/status.d/$(basename "$base_status_file")
    old_version=$(awk '/^Version:/{print $2; exit}' "$base_status_file")
    op=gt
    if [[ $extras == *" $pkg "* ]]; then op=ge; fi
    if ! dpkg --compare-versions "$new_version" "$op" "$old_version"; then
      echo "keep:   $pkg $old_version (available: $new_version)"
      continue
    fi
  else
    old_version=none
  fi
  echo "update: $pkg $old_version -> $new_version"
  updated+=("$status_file")

  # Checksums distroless recorded for the old version, to detect the files it customized
  declare -A base_md5=()
  if [[ -f $base_status_file.md5sums ]]; then
    while read -r hash path; do base_md5[$path]=$hash; done < "$base_status_file.md5sums"
  fi

  # Every path of the package, except man pages and the merged-usr top-level symlinks (/bin, /lib, ...),
  # since distroless keeps those as real directories
  : > "$work/files"
  : > "$work/skipped"
  dpkg-query -L "$pkg" | { grep '^/' || true; } \
    | { grep -Ev '^/(\.|bin|sbin|lib|lib32|lib64|libx32)$' || true; } \
    | { grep -Ev '^/usr/share/man(/|$)' || true; } | sed 's|^/||' > "$work/paths"
  while read -r path; do
    f=/$path
    if [[ ! -e $f && ! -L $f ]]; then continue; fi
    if [[ -e $base/$path || -L $base/$path ]]; then
      base_type=$(path_type "$base/$path")
      debian_type=$(path_type "$f")
      if [[ $base_type != "$debian_type" ]]; then
        echo "skip:   $pkg: $path (base: $base_type, debian: $debian_type)"
        echo "$path" >> "$work/skipped"
        continue
      fi
      if [[ $base_type == file && -n ${base_md5[$path]:-} ]] \
        && [[ $(md5sum < "$base/$path" | cut -d' ' -f1) != "${base_md5[$path]}" ]]; then
        echo "skip:   $pkg: $path (customized in base)"
        echo "$path" >> "$work/skipped"
        continue
      fi
    fi
    echo "$path" >> "$work/files"
  done < "$work/paths"
  tar -C / --no-recursion -cf - -T "$work/files" | tar -C "$delta" -xpf -

  # Files of the old version that the new version no longer ships are left in the image, untracked
  if [[ -f $base_status_file.md5sums ]]; then
    sed 's/^[0-9a-f]*  //' "$base_status_file.md5sums" | grep -vxF -f "$work/paths" | grep -v '^usr/share/man/' \
      | sed "s|^|orphan: $pkg: |" || true
  fi
  unset base_md5

  # distroless metadata has no Status or Conffiles fields
  dpkg-query -s "$pkg" | awk '/^(Status|Conffiles):/{skip=1; next} skip && /^ /{next} {skip=0; print}' > "$delta/$status_file"

  # Use the checksums from the Debian package, so the check below verifies the shipped files against Debian
  md5sums=$(dpkg-query --control-path "$pkg" md5sums)
  if [[ -z $md5sums ]]; then
    echo "error: $pkg has no md5sums to verify its files against" >&2
    exit 1
  fi
  grep -v '^[0-9a-f]*  usr/share/man/' "$md5sums" > "$delta/$status_file.md5sums" || true
  # Skipped paths are not shipped, so they are not checked
  awk 'FILENAME == ARGV[1] {skip[$0] = 1; next} !(substr($0, index($0, "  ") + 2) in skip)' \
    "$work/skipped" "$delta/$status_file.md5sums" > "$work/$(basename "$status_file").check"
done

# The certificate bundle is generated on install, so it is not part of the package file list
if [[ " ${updated[*]} " == *"/ca-certificates "* ]]; then
  mkdir -p "$delta/etc/ssl/certs"
  cp -a /etc/ssl/certs/. "$delta/etc/ssl/certs/"
fi

# Fail if a file that is shipped is missing or does not match the Debian package
cd "$delta"
for status_file in "${updated[@]}"; do
  md5sum --quiet -c "$work/$(basename "$status_file").check"
done
rm -rf "$work"
