# Copyright (c) 2021 AccelByte Inc. All Rights Reserved.
# This is licensed software from AccelByte Inc, for limitations
# and restrictions contact your company contract manager.

FROM alpine:3.20

ENV HOME /srv

COPY service $HOME/service

RUN find $HOME -type d -exec 'chmod' '555' '{}' ';' && \
    find $HOME -type f -exec 'chmod' '444' '{}' ';' && \
    find $HOME -type f -exec 'chown' 'root:root' '{}' ';' && \
    chmod 555 $HOME/service

USER nobody

ENTRYPOINT ["/srv/service"]
