# n0depend

A tool that print task dependency of [n0stack](https://github.com/n0stack/n0stack) as graph.

This tool uses `graph-easy` from [Graph::Easy](https://metacpan.org/release/Graph-Easy) to print graphs.

To install `graph-easy`, execute `cpanm Graph::Easy`.

## Example

```
$ cat sample.yaml
DeleteVM_a:
  ...
  action: DeleteVM
  args:
    name: wow_a

DeleteBlockStorage_a:
  ...
  depends_on:
    - DeleteVM_a

DeleteVM_b:
  ...

DeleteBlockStorage_b:
  ...
  depends_on:
    - DeleteVM_b

GenerateBlockStorage_a:
  ...
  depends_on:
    - DeleteBlockStorage_a

CreateVM_a:
  ...
  depends_on:
    - GenerateBlockStorage_a

GenerateBlockStorage_b:
  ...
  action: GenerateBlockStorage
  depends_on:
    - CreateVM_a
    - DeleteBlockStorage_b

CreateVM_b:
  ...
  depends_on:
    - GenerateBlockStorage_b
$ n0depend sample.yaml
+------------+     +----------------------+     +------------------------+     +------------+     +------------------------+     +------------+
| DeleteVM_a | --> | DeleteBlockStorage_a | --> | GenerateBlockStorage_a | --> | CreateVM_a | --> | GenerateBlockStorage_b | --> | CreateVM_b |
+------------+     +----------------------+     +------------------------+     +------------+     +------------------------+     +------------+
                                                                                                    ^
                                                                                                    |
                                                                                                    |
+------------+     +----------------------+                                                         |
| DeleteVM_b | --> | DeleteBlockStorage_b | --------------------------------------------------------+
+------------+     +----------------------+
```
