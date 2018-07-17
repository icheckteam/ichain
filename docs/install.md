# Install


Ichain can be installed to
`$GOPATH/src/github.com/icheckteam/ichain` like a normal Go program:

```
go get github.com/icheckteam/ichain
```

If the dependencies have been updated with breaking changes, or if
another branch is required, `dep` is used for dependency management.
Thus, assuming you've already run `go get` or otherwise cloned the repo,
the correct way to install is:

```
cd $GOPATH/src/github.com/icheckteam/ichain
make get_tools
make get_vendor_deps
make install
```

Verify that everything is OK by running:

```
ichaind version
```

you should see:

```
0.20.0-dev
```

then with:

```
ichaincli version
```
you should see the same version (or a later one for both).

## Update

Get latest code (you can also `git fetch` only the version desired),
ensure the dependencies are up to date, then recompile.

```
cd $GOPATH/src/github.com/icheckteam/ichain
git fetch -a origin
git checkout VERSION
make get_vendor_deps
make install
```