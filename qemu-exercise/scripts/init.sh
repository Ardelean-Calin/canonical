#!/bin/sh
mount -t devtmpfs none /dev
mount -t proc none /proc
mount -t sysfs none /sys
# Print it green!
echo -e "\033[0;32mHello world - from Calin\033[0m"
exec /bin/sh
