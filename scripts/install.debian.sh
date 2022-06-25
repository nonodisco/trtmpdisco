ARCH=$(dpkg --print-architecture) && curl -s https://api.github.com/repos/nonodisco/trtmpdisco/releases/latest | grep "linux-$ARCH" | cut -d : -f 2,3 | tr -d \" | wget -O "linux-$ARCH" -qi -
