fienode - discover identical CoW copies
=======================================

Something analogous to an inode for Copy-on-Write (CoW, `cp --reflink`) files.

Discover when two files on a CoW filesystem share identical physical data.

Why's it called `fienode`? It's analogous to an inode,
but derived from **Fi**le **E**xtents.

#### Warning: Proof-of-concept. Alpha quality software, use at your own risk.

**See Also**:

* [`sharedextents`](https://github.com/pwaller/sharedextents), for comparing physical extents between two files
* [`filefrag -v`](https://en.wikipedia.org/wiki/E2fsprogs), for listing physical extents

## Installation

#### From source

```
go get github.com/pwaller/fienode
```

#### Binary Download

See [releases](https://github.com/pwaller/fienode/releases/).

* [fienode-linux-amd64]( https://github.com/pwaller/fienode/releases/download/v1.0/fienode-linux-amd64) v1.0

## Usage: `fienode <filename>...`

For example:

```
$ fienode foo
foo 8a53b838c97f1f9712e6a77e2bc00dd1922d32e9
$ cp foo bar
$ fienode bar
bar c829d10e6ccc0f90fc41d1492e423e2cabfe2bca
$ cp --reflink foo baz
$ fienode baz
baz 8a53b838c97f1f9712e6a77e2bc00dd1922d32e9

# (note: baz and foo share the same hash)
```

## How it works

The result is the a SHA1 hash of the physical extents of the file.

See
["How to verify a file copy is reflink/CoW?" on unix.stackexchange.com](http://unix.stackexchange.com/a/277033/26224)
for more information.

## Caveats

There may be bugs. This will delete all your data and eat your cat.
When it does, that is your problem. Keep backups, folks.

I just copied a large file `x` to `y` with `cp --reflink x y`. For about a
minute, it gave equal results. Thereafter, BTRFS decided - seemingly at
random - to coalesce two extents in one of the files but not the other.
So fienode then returned different results, even though the majority of the
file was actually shared.

#### License

MIT.
