# bene-api
Rewrite of github.com/minkezhang/truffle.

## Adding a New Atom Type

1. Add new metadata message to `/proto/metadata.proto`
1. Add message as a field to `/proto/atom.proto`
1. Create a package in `/db/atom/metadata/${TYPE}/`
1. Implement the `metadata.G` and `metadata.M` interfaces for this new type
1. Explicitly define the `metadata.Equal` function in the package (otherwise
   `cmp.Diff` will incorrectly detect struct differences
1. Update `atom.MergeMetadata`, `atom.Load`, and `atom.Save`
