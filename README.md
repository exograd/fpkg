# fpkg
Fpkg is a small program to interact with FreeBSD packages. While the official
[`pkg` program](https://github.com/freebsd/pkg) works perfectly well for
package management, it is not as practical to build packages. Among other
things:

- It forces the user to manually build the file and directory list and their
  checksum and permissions.
- It requires writing a shell script to create a user and/or group.
- It detects the ABI by analyzing ELF files in the local system, a method that
  fails on Alpine Linux (https://github.com/freebsd/pkg/issues/2065).
- It is not packaged out of FreeBSD, making it annoying to cross-build FreeBSD
  packages. Building it manually requires multiple dependencies.

## Package building
Fpkg uses a simple YAML configuration file describing the package to build.

Example:
```yaml
name: "example"
version: "1.0.0"
short_description: "example package"
website_uri: "https://github.com/exograd/example"
maintainer: "Nicolas Martyanoff <nicolas@n16f.net>"
file_owner: "root"
file_group: "wheel"
files:
  - path: "/var/lib/example"
    mode: "600"
  - path_regexp: "/var/www/example/.*"
    owner: "www"
    group: "www"
```

You can then run fpkg:
```
fpkg build -c example.yaml example/
```

Where `example/` is the directory containing the set of files to include in
the package.

Fpkg automatically builds the file and directory index, including the
checksum, permissions, and the owner and group set in the manifest.

The path of the resulting `.pkg` file is printed on `stdout`; this way a
script running fpkg can easily find and copy the package archive to a remote
repository.
