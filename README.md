# sshsrv

`sshsrv` is a simple program to lookup and connect to an SSH endpoint
via DNS SRV records.

## Why SRV Records?

Consider the scenario where there are multiple hosts NATted behind a
single IPv4 address.  There are two usual methods for handling
connecting to these internal hosts from outside the network:

1. Port-forward each internal host's SSH port to some high port on the
   gateway device; or
1. Forward port 22 to a single, internal "jump-host" and either manually
   connect to other internal hosts or proxy (using the `-W` option to
   `ssh(1)`) from there.

The first option requires you to remember _which_ external port you've
assigned to which internal host.  The second option requires an extra
hop between you and your internal hosts.  Using SRV records instead
means that you can connect directly to your internal machines just like
the first option above, and provides the added benefit of a mechanism to
look up which port you've assigned to each internal server.

## Example

Say you have a home network with multiple hosts, but like most of the
world you only have one external IPv4 address.  Set up a SRV record like
the following:

```
$ dig +short _ssh._tcp.myhost.mydomain.com SRV
1 1 22029 gateway.mydomain.com.
```

Using `sshsrv`, when you want to connect to `myhost.mydomain.com` it
will look up this SRV record, which will tell `sshsrv` to actually
connect to `gateway.mydomain.com` on port 22029 instead.  (If no SRV
record exists, `sshsrv` will simply pass the hostname directly to
`ssh(1)`.)

## Usage

`sshsrv` strives to accept the same options that the `ssh(1)` command
does, so that it can be a drop-in replacement for the `ssh(1)`
command. If there is a valid flag that isn't supported, please open an
issue or PR!

## Installation

You can download the latest pre-built releases [here][releases].

Alternatively, you can install it with Go 1.17 or higher like so:

```
$ go install github.com/Crosse/sshsrv@latest
```

## Why not just submit a patch to OpenSSH?
Using a wrapper allows the user to use whatever version of `ssh(1)` is
installed on their system, instead of being an OpenSSH-only addition.

## Questions, Comments, Suggestions?
Submit a pull request!

[releases]: https://github.com/Crosse/sshsrv/releases/latest
