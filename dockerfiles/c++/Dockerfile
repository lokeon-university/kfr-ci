FROM debian:stretch-slim

RUN apt-get update && \
apt-get install -y wget git gnupg make locales locales-all valgrind && \
wget -O - https://apt.llvm.org/llvm-snapshot.gpg.key | apt-key add - && \
echo "deb http://apt.llvm.org/stretch/ llvm-toolchain-stretch-7 main" >> /etc/apt/sources.list && \
apt-get update && \
apt-get install -y clang-7 lldb-7 lld-7 clang-tools-7 libclang-common-7-dev libclang-7-dev libclang1-7 clang-format-7 libfuzzer-7-dev libc++-7-dev libc++abi-7-dev libomp-7-dev && \
ln -s /usr/bin/llvm-config-7 /usr/bin/llvm-config && \
ln -s /usr/bin/clang-7 /usr/bin/clang && \
ln -s /usr/bin/clang++-7 /usr/bin/clang++ && \
apt-get remove --purge -y wget gnupg 

COPY docker-entrypoint.sh /usr/local/bin/

RUN ln -s usr/local/bin/docker-entrypoint.sh /

ENTRYPOINT ["docker-entrypoint.sh"]
