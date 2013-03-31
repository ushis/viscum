# viscum RSS/ATOM fetching and processing server

viscum is a RSS/ATOM fetching and processing server powered by
[Go](http://golang.org/) and [PostgreSQL](http://www.postgresql.org/).

## Architecture

Lets see how overengineered it is.

                     +--------------------------------------+
                     |                 viscumd              |
                     +--------------------------------------+
                          |               |             |
                          |         +-----------+       |         +----------+
      +--------+          |         |           |<--<--<|-<--<--<-| RSS/ATOM |
      |  SMTP  |<-+  +----------+   |  Fetcher  |   +-------+     +----------+
      +--------+  |  |          |<--|           |<--|       |
                  +--|  Mailer  |   +-----------+   |  RPC  |     +----------+
      +--------+  |  |          |<--------|---------|       |<--->|  viscum  |
      |  Pipe  |<-+  +----------+         |         +-------+     +----------+
      +--------+          |               |             |
                     +--------------------------------------+
                     |                PostgreSQL            |
                     +--------------------------------------+

1. viscumd starts the Mailer, the Fetcher and the RPC in separate
   [Goroutines](http://golang.org/doc/effective_go.html#goroutines)

1. The Fetcher fetches some feeds, stores new entries in the database and
   notifies the Mailer about new stuff in the queue.

1. The Mailer looks into the queue and sends new entries to an external
   program (e.g. mail) or to a SMTP server.

1. The RPC waits for a viscum client to add/rm subscriptions, send control
   instructions or request some info.

## Get it

Clone the repo and build the binaries. You will need a go compiler and
[godag](https://code.google.com/p/godag/) to achieve this goal.

    git clone git://github.com/ushis/viscum.git
    cd viscum
    make

Create a postgres database user, a database and run the initial setup script.

    createuser viscum
    createdb -O viscum viscum
    psql -U viscum viscum < share/postgre.sql

Edit the config ```etc/viscumd.conf```, start the server...

    ./build/viscumd -config=etc/viscumd.conf

...and play with the client.

    ./build/viscum -config=etc/viscum.conf add hello@example.com https://github.com/blog.atom

Arch Linux user can grabit from the
[AUR](https://aur.archlinux.org/packages/viscum-git/).

## TODO

Well, this project is totally premature. It is missing test, has bugs and many
many rough edges.

## LICENSE (MIT)

Copyright (c) 2013 ushi

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
