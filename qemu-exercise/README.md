To run the first example, simply run 
```
./run.sh -b
```

The script will download the latest LTS kernel as of today, compile it with default options (tinyconfig needed some additional tweaking so I skipped using it), then build a busybox-based filesystem image and finally run the whole thing inside Qemu x86_64.

The optional `-b` flag also creates a **b**ootable ISO `calinos.iso` which you can then test with the following command:

```
qemu-system-x86_64 -boot d -cdrom calinos.iso -m 512
```
