Use the rsync plugin to deploy files to a server using rsync over ssh. The following parameters are used to configure this plugin:

* `user` - connects as this user
* `host` - connects to this host address
* `port` - connects to this host port
* `source` - source path from which files are copied
* `target` - target path to which files are copied
* `delete` - delete extraneous files from the target dir
* `recursive` - recursively transfer all files
* `include` - include files matching the specified pattern
* `exclude` - exclude files matching the specified pattern
* `filter` - include or exclude files according to filtering rules
* `script` - execute commands on the remote host after files are copied
* `key` - private SSH key for the remote machine

The following secret values can be set to configure the plugin.

* `SSH_KEY` - corresponds to **key**

It is highly recommended to put the **SSH_KEY** into a secret so it is not
exposed to users. This can be done using the drone-cli.

```bash
drone secret add --image=plugins/rsync \
    octocat/hello-world SSH_KEY @path/to/.ssh/id_rsa
```

Then sign the YAML file after all secrets are added.

```bash
drone sign octocat/hello-world
```

See [secrets](http://readme.drone.io/0.5/usage/secrets/) for additional
information on secrets

## Examples

Sample configuration in the `.drone.yml` file:

```yaml
deploy:
  rsync:
    user: root
    host: 127.0.0.1
    port: 22
    source: copy/files/from
    target: send/files/to
    delete: false
    recursive: true
    exclude:
      - "exclude/this/pattern/*"
      - "or/this/one"
    commands:
      - service nginx restart
```

## Common Problems

The below error message may be encountered when rsync is unable to connect to your server. First verify you can connect to your remote server via ssh.

```
rsync: connection unexpectedly closed (0 bytes received so far) [sender]
rsync error: unexplained error (code 255) at io.c(226) [sender=3.1.1]
```
