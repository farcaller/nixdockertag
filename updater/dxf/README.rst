|Build Status| |Coverage Status| |PyPI version|

Python module and command-line tool for storing and retrieving data in a
Docker registry.

-  Store arbitrary data (blob-store)
-  Content addressable
-  Set up named aliases to blobs
-  Supports Docker registry schemas v1 and v2
-  Works on Python 3.6+

Please note that ``dxf`` does *not* generate Docker container
configuration, so you won’t be able to ``docker pull`` data you store
using ``dxf``. See `this
issue <https://github.com/davedoesdev/dxf/issues/3>`__ for more details.

Command-line example:

.. code:: shell

   dxf push-blob fred/datalogger logger.dat @may15-readings
   dxf pull-blob fred/datalogger @may15-readings

which is the same as:

.. code:: shell

   dxf set-alias fred/datalogger may15-readings $(dxf push-blob fred/datalogger logger.dat)
   dxf pull-blob fred/datalogger $(dxf get-alias fred/datalogger may15-readings)

Module example:

.. code:: python

   from dxf import DXF

   def auth(dxf, response):
       dxf.authenticate('fred', 'somepassword', response=response)

   dxf = DXF('registry-1.docker.io', 'fred/datalogger', auth)

   dgst = dxf.push_blob('logger.dat')
   dxf.set_alias('may15-readings', dgst)

   assert dxf.get_alias('may15-readings') == [dgst]

   for chunk in dxf.pull_blob(dgst):
       sys.stdout.write(chunk)

Usage
-----

The module API is described
`here <http://rawgit.davedoesdev.com/davedoesdev/dxf/master/docs/_build/html/index.html>`__.

The ``dxf`` command-line tool uses the following environment variables:

-  ``DXF_HOST`` - Host where Docker registry is running.
-  ``DXF_INSECURE`` - Set this to ``1`` if you want to connect to the
   registry using ``http`` rather than ``https`` (which is the default).
-  ``DXF_USERNAME`` - Name of user to authenticate as.
-  ``DXF_PASSWORD`` - User’s password.
-  ``DXF_AUTHORIZATION`` - HTTP ``Authorization`` header value.
-  ``DXF_AUTH_HOST`` - If set, always perform token authentication to
   this host, overriding the value returned by the registry.
-  ``DXF_PROGRESS`` - If this is set to ``1``, a progress bar is
   displayed (on standard error) during ``push-blob`` and ``pull-blob``.
   If this is set to ``0``, a progress bar is not displayed. If this is
   set to any other value, a progress bar is only displayed if standard
   error is a terminal.
-  ``DXF_BLOB_INFO`` - Set this to ``1`` if you want ``pull-blob`` to
   prepend each blob with its digest and size (printed in plain text,
   separated by a space and followed by a newline).
-  ``DXF_CHUNK_SIZE`` - Number of bytes ``pull-blob`` should download at
   a time. Defaults to 8192.
-  ``DXF_SKIPTLSVERIFY`` - Set this to ``1`` to skip TLS certificate
   verification.
-  ``DXF_TLSVERIFY`` - Optional path to custom CA bundle to use for TLS
   verification.
-  ``DXF_PLATFORM`` - Optional platform (e.g. ``linux/amd64``) to use
   for multi-arch manifests. If a multi-arch manifest is encountered and
   this is not set then a dict containing entries for each platform will
   be displayed.

You can use the following options with ``dxf``. Supply the name of the
repository you wish to work with in each case as the second argument.

-  ``dxf push-blob <repo> <file> [@alias]``

      Upload a file to the registry and optionally give it a name
      (alias). The blob’s hash is printed to standard output.

   ..

      The hash or the alias can be used to fetch the blob later using
      ``pull-blob``.

-  ``dxf pull-blob <repo> <hash>|<@alias>...``

      Download blobs from the registry to standard output. For each blob
      you can specify its hash, prefixed by ``sha256:`` (remember the
      registry is content-addressable) or an alias you’ve given it
      (using ``push-blob`` or ``set-alias``).

-  ``dxf blob-size <repo> <hash>|<@alias>...``

      Print the size of blobs in the registry. If you specify an alias,
      the sum of all the blobs it points to will be printed.

-  ``dxf mount-blob <repo> <from-repo> <hash> [@alias]``

      Cross mount a blob from another repository and optionally give it
      an alias. Specify the blob by its hash, prefixed by ``sha256:``.

   ..

      This is useful to avoid having to upload a blob to your repository
      if you know it already exists in the registry.

-  ``dxf del-blob <repo> <hash>|<@alias>...``

      Delete blobs from the registry. If you specify an alias the blobs
      it points to will be deleted, not the alias itself. Use
      ``del-alias`` for that.

-  ``dxf set-alias <repo> <alias> <hash>|<file>...``

      Give a name (alias) to a set of blobs. For each blob you can
      either specify its hash (as printed by ``push-blob`` or
      ``get-alias``) or, if you have the blob’s contents on disk, its
      filename (including a path separator to distinguish it from a
      hash).

-  ``dxf get-alias <repo> <alias>...``

      For each alias you specify, print the hashes of all the blobs it
      points to.

-  ``dxf del-alias <repo> <alias>...``

      Delete each specified alias. The blobs they point to won’t be
      deleted (use ``del-blob`` for that), but their hashes will be
      printed.

-  ``dxf get-digest <repo> <alias>...``

      For each alias you specify, print the hash of its configuration
      blob. For an alias created using ``dxf``, this is the hash of the
      first blob it points to. For a Docker image tag, this is the same
      as ``docker inspect alias --format='{{.Id}}'``.

-  ``dxf get-manifest <repo> <alias>...``

      For each alias you specify, print its manifest obtained from the
      registry.

-  ``dxf list-aliases <repo>``

      Print all the aliases defined in the repository.

-  ``dxf list-repos``

      Print the names of all the repositories in the registry. Not all
      versions of the registry support this.

Certificates
------------

If your registry uses SSL with a self-issued certificate, you’ll need to
supply ``dxf`` with a set of trusted certificate authorities.

Set the ``REQUESTS_CA_BUNDLE`` environment variable to the path of a PEM
file containing the trusted certificate authority certificates.

Both the module and command-line tool support ``REQUESTS_CA_BUNDLE``.

Alternatively, you can set the ``DXF_TLSVERIFY`` environment variable
for the command-line tool or pass the ``tlsverify`` option to the
module.

Authentication tokens
---------------------

``dxf`` automatically obtains Docker registry authentication tokens
using your ``DXF_USERNAME`` and ``DXF_PASSWORD``, or
``DXF_AUTHORIZATION``, environment variables as necessary.

However, if you wish to override this then you can use the following
command:

-  ``dxf auth <repo> <action>...``

      Authenticate to the registry using ``DXF_USERNAME`` and
      ``DXF_PASSWORD``, or ``DXF_AUTHORIZATION``, and print the
      resulting token.

   ..

      ``action`` can be ``pull``, ``push`` or ``*``.

If you assign the token to the ``DXF_TOKEN`` environment variable, for
example:

``DXF_TOKEN=$(dxf auth fred/datalogger pull)``

then subsequent ``dxf`` commands will use the token without needing
``DXF_USERNAME`` and ``DXF_PASSWORD``, or ``DXF_AUTHORIZATION``, to be
set.

Note however that the token expires after a few minutes, after which
``dxf`` will exit with ``EACCES``.

Docker Cloud authentication
---------------------------

You can use the
```dockercloud`` <https://github.com/docker/python-dockercloud>`__
library to read authentication information from your Docker
configuration file and pass it to ``dxf``:

.. code:: python

   auth = 'Basic ' + dockercloud.api.auth.load_from_file()
   dxf_obj = dxf.DXF('index.docker.io', repo='myorganization/myimage')
   dxf_obj.authenticate(authorization=auth, actions=['pull'])
   dxf_obj.list_aliases()

Thanks to `cyrilleverrier <https://github.com/cyrilleverrier>`__ for
this tip.

Installation
------------

.. code:: shell

   pip install python-dxf

Licence
-------

`MIT <https://raw.github.com/davedoesdev/dxf/master/LICENCE>`__

Other projects that use DXF
---------------------------

Docker-charon
~~~~~~~~~~~~~

https://github.com/gabrieldemarmiesse/docker-charon

This package allows you to transfer Docker images from one registry to
another. The second one being disconnected from the internet.

Unlike ``docker save`` and ``docker load``, it creates the payload
directly from the registry (it’s faster) and is able to compute diffs to
only take the layers needed, hence reducing the size.

Tests
-----

.. code:: shell

   make test

Lint
----

.. code:: shell

   make lint

Code Coverage
-------------

.. code:: shell

   make coverage

`coverage.py <http://nedbatchelder.com/code/coverage/>`__ results are
available
`here <http://rawgit.davedoesdev.com/davedoesdev/dxf/master/htmlcov/index.html>`__.

Coveralls page is `here <https://coveralls.io/r/davedoesdev/dxf>`__.

.. |Build Status| image:: https://github.com/davedoesdev/dxf/workflows/ci/badge.svg
   :target: https://github.com/davedoesdev/dxf/actions
.. |Coverage Status| image:: https://coveralls.io/repos/davedoesdev/dxf/badge.png?branch=master
   :target: https://coveralls.io/r/davedoesdev/dxf?branch=master
.. |PyPI version| image:: https://badge.fury.io/py/python-dxf.png
   :target: http://badge.fury.io/py/python-dxf
