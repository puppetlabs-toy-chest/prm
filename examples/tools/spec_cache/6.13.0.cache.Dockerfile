# docker build -t pdk:puppet-6.13.0-spec_cache -f examples/tools/spec_cache/6.13.0.cache.Dockerfile .
# docker run --rm -v ${PWD}:/code -v C:\Users\dave\.pdk\prm\cache:/cache pdk:puppet-6.13.0-spec_cache

FROM puppet/puppet-agent:6.13.0
# Git needs to be installed to handle gem sources not in rubygems.org
RUN apt update
RUN apt install git -y
RUN /opt/puppetlabs/puppet/bin/gem install bundler --no-document

COPY examples/tools/spec_cache/content/* /tmp/

# mount the local code
VOLUME [ "/code", "/cache" ]
WORKDIR "/code"

ENTRYPOINT [ "/tmp/cache.sh" ]
