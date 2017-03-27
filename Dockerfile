FROM concourse/buildroot:git

COPY assets/ /opt/resource/
COPY test.sh /opt/resource/
RUN  chmod +x /opt/resource/test.sh

COPY git-phlow /bin



ENV TMPDIR=/gitphlow
RUN mkdir $TMPDIR
