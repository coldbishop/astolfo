|go_report|

astolfo
=======

``astolfo`` is a password generator written in Go.

``astolfo`` relies on a full name and a unique name (usually the name of a
website) given by the user to generate a password. The resulting password is
never stored anywhere, thus making it harder for potential attackers to steal
it.

Install
-------

See releases_ for pre-built binaries.

If you want to build from source:

::

    go get -u github.com/coldbishop/astolfo

.. _releases: https://github.com/coldbishop/astolfo/releases

Usage
-----

Standard usage:

::

    astolfo "Your Name" example.com

Generate a 32 length password consisting of only uppercase and lowercase letters:

::

    astolfo --length 32 -U -l "John Bishop" twitter.com

Generate a second password with the same params as above:

::

    astolfo --counter 2 --length 32 -U -l "John Bishop" twitter.com

Generate a four-digit PIN number:

::

    astolfo -L 4 --digit "Caroline" "A local ATM"

Unicode (UTF-8 only) is supported.

::

    astolfo "猫宮ひなた" pixiv.net

Run ``astolfo --help`` or ``man 1 astolfo`` for more information.

Feature
-------

- Strong security courtesy of Argon2_ and BLAKE2_ algorithms
- Easily generate multiple passwords using the same name and site name
- Cross platform (GNU/Linux, macOS, Windows)
- Portable, single binary without any runtime dependencies
- Unicode support (UTF-8) for user name, site name, and even the master password
- Automatic copy of the generated password to clipboard

.. _Argon2: https://www.argon2.com
.. _BLAKE2: https://blake2.net

Non-feature
-----------

- Store passwords anywhere (on disk or in the cloud)
- Synchronization

To-do
-----

- Diceware generator
- Text-based user interface (TUI)
- Graphical user interface (GUI)

License
-------

``astolfo`` is licensed under the zlib License. See `LICENSE.rst`_.

.. _LICENSE.rst: https://github.com/coldbishop/astolfo/blob/master/LICENSE.rst
.. |go_report| image:: https://goreportcard.com/badge/github.com/coldbishop/astolfo
   :target: https://goreportcard.com/report/github.com/coldbishop/astolfo

