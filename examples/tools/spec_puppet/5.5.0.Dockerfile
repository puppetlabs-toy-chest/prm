# based upon the work in https://github.com/da-ar/pcv-test/
# this file is a mock up of what a generated file based on the prm-config.yml
# would have to contain in order to work.

# build images:
### docker build -t pdk:puppet-5.5.0-spec_cache -f examples/tools/spec_cache/5.5.0.cache.Dockerfile .
### docker build -t pdk:puppet-5.5.0-puppet_spec -f examples/tools/spec_puppet/5.5.0.Dockerfile .

# cache gems:
### docker run --rm -v ${PWD}:/code -v C:\Users\dave\.pdk\prm\cache:/cache pdk:puppet-5.5.0-spec_cache

# step 1:
### docker run --rm -v ${PWD}:/code -v C:\Users\dave\.pdk\prm\cache:/cache:ro -w /code pdk:puppet-5.5.0-puppet_spec spec_prep

# step 2
### docker run --network none --rm -v ${PWD}:/code -v C:\Users\dave\.pdk\prm\cache:/cache:ro -w /code pdk:puppet-5.5.0-puppet_spec


# gem - gems need ruby - so get puppet...
FROM puppet/puppet-agent:5.5.0

# fix for puppet 5 - this would need baked in to the template
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 4528B6CD9E61EF26

# requires_git: true
RUN apt update
RUN apt install git -y
# build_tools: true
RUN apt install build-essential -y

# All GEMs must
### user must run "bundle cache --no-install --all --all-platforms"
# we could have a command for this in PRM
RUN /opt/puppetlabs/puppet/bin/gem install bundler --no-document

# name: ["puppetlabs_spec_helper", "rspec-puppet-facts"]
#  2.4: ["puppetlabs_spec_helper", "2.15.0" ]
RUN /opt/puppetlabs/puppet/bin/gem install puppetlabs_spec_helper -f --conservative --minimal-deps -v '2.15.0' --no-document
RUN /opt/puppetlabs/puppet/bin/gem install rspec-puppet-facts -f --conservative --minimal-deps --no-document

# copy any file from content over to /tmp
COPY  examples/tools/spec_puppet/content/* /tmp/

# mount the local code
VOLUME [ "/code", "/cache" ]
WORKDIR /code

# Set up the default container command - we need some code to run before
# we can run the container command. This is being done in docker-entrypoint.sh.

# ENV: ['CMD_ENTRY', '/opt/puppetlabs/puppet/bin/rake -f /tmp/Rakefile']
ENV CMD_ENTRY="/opt/puppetlabs/puppet/bin/rake -f /tmp/Rakefile"

# use_entrypoint_script: "docker-entrypoint"
ENTRYPOINT [ "/tmp/docker-entrypoint.sh"]
# OR
# executable: rake
# ENTRYPOINT [ "/opt/puppetlabs/puppet/bin/rake", "-f", "/tmp/Rakefile"]

# default_args: [spec_standalone]
CMD [ "spec_standalone" ]
