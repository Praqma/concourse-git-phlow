FROM concourse/buildroot:git

COPY concourse-git-phlow/assets/ /opt/resource/
COPY concourse-git-phlow/git-phlow /bin

ENV TMPDIR=/gitphlow
RUN mkdir $TMPDIR
