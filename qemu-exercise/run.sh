#!/usr/bin/env bash
set -e

KERNEL_VERSION="linux-6.1.74"
BUSYBOX_VERSION="busybox-1_36_1"

####### Step 0: Check requirements
sudo apt update && sudo apt install -y qemu-system-x86 git wget tar fakeroot build-essential ncurses-dev xz-utils libssl-dev bc flex libelf-dev bison grub-common xorriso

# ####### Step 1: Download and extract the Linux kernel (I chose LTS)
if [[ ! -f "${KERNEL_VERSION}.tar.xz" ]]; then
  wget "https://cdn.kernel.org/pub/linux/kernel/v6.x/${KERNEL_VERSION}.tar.xz"
fi
tar xvf "${KERNEL_VERSION}.tar.xz"
rm "${KERNEL_VERSION}.tar.xz"

####### Step 2: Compile the kernel
cd "${KERNEL_VERSION}"
# Create the default configuration. Had problems with tinyconfig so I had to cut short to fit in 2 hours
# NOTE: I assume you are running on an x86_64 system (Ubuntu 20.04 and 22.04 don't support 32-bit anyways, and that is the platform the Exercise says we run on).
# If desired, cross-compilation can be done.
make defconfig
# Build the bare kernel
make -j$(nproc) 2>&1 bzImage | tee kernel-log

####### Step 3: Download, configure and compile Busybox
if [[ ! -f "${BUSYBOX_VERSION}.tar.bz2" ]]; then
  wget "https://git.busybox.net/busybox/snapshot/${BUSYBOX_VERSION}.tar.bz2"
fi
tar xvf "${BUSYBOX_VERSION}.tar.bz2"
rm "${BUSYBOX_VERSION}.tar.bz2"

cd "${BUSYBOX_VERSION}"

# Default busybox config...
make defconfig
# ...but with static library build enabled. I could of course just copy a .config file but I wanted to have nothing but the
# shell script mentioned in the exercise description.
sed -i 's/# CONFIG_STATIC is not set/CONFIG_STATIC=y/' .config
make -j$(nproc) 2>&1 | tee busybox-log

####### Step 4: Create the filesystem
make install
cd _install
mkdir -p dev proc sys
cp ../../../scripts/init.sh ./init
chmod +x init

find . -print0 | cpio --null -ov --format=newc | gzip -9 > ../initramfs.cpio.gz

# Back to Project Root
cd ../../../

###### BONUS! Create a bootable ISO
GREEN='\033[0;32m'
NC='\033[0m' 
while getopts ':b' opt; do
  case "$opt" in
    b)
      echo -e "${GREEN}Creating Bootable ISO...${NC}"
      mkdir -p iso/boot/grub
      cp "$KERNEL_VERSION/arch/x86_64/boot/bzImage" iso/boot/
      cp "$KERNEL_VERSION/${BUSYBOX_VERSION}/initramfs.cpio.gz" iso/boot/
      if [[ -d /sys/firmware/efi ]]; then
        echo -e "${GREEN}Detected EFI system.${NC}"
        cp config/grub_efi.cfg iso/boot/grub/grub.cfg
      else
        echo -e "${GREEN}Detected BIOS system.${NC}"
        cp config/grub_bios.cfg iso/boot/grub/grub.cfg
      fi
      grub-mkrescue -o calinos.iso iso/
      echo -e "${GREEN}Done! Created .iso image calinos.iso${NC}"
      ;;
  esac
done
  
# Finally run the Linux kernel. Hopefully we get a hello world!
qemu-system-x86_64 -kernel "$KERNEL_VERSION/arch/x86_64/boot/bzImage" -initrd "$KERNEL_VERSION/$BUSYBOX_VERSION/initramfs.cpio.gz" -nographic -append 'console=ttyS0' -m 512

