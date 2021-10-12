# SoulWorker Server
My slow attempt at an 'RE' of the SoulWorker protocol.

How to use:

1. Download and open in GoLand or whatever IDE (or not) you use
2. Set up the database with the provided .sql file

To launch the KRSW OnStove client (+skipping Stove Authenticator):
```
SoulWorker.exe Live/127.0.0.1/10000 SkipSGAuthen:yes
```

To bypass XIGNCode3, you can build the dll found [here](https://github.com/Lazarus/XignCode3-bypass).

To bypass the OnStove authenticator without command line arguments, you can build the dll found [here](https://github.com/Lazarus/OnStove-Client).

Changelog:

* Oct 11, 2021
  * Re-base

* May 11, 2020
  * Major refactor
  * Small database implementation
  * Somehow made it ask for 2fa and don't know how to skip, so done for today

* August 2, 2019
  * No server sided updates, only client side discoveries.

* July 29, 2019
  * Moved hardcoded keyTable to globals where it's defined
  * Minimal file updates, just slight renaming on a few variables
  * Changed cringe server name
  * Added some tools/tests for `tb_*.res` table extraction/verification
    * Includes `tb_*` string dump and the ASM dump for functions regarding them
    * Includes the final go struct dump for all table structures
    * You can find the original script that I modified in this [Gist](https://gist.github.com/x1nixmzeng/a4a5c419f1cd4bc72cba30d5e647bc4f)

* July 13, 2019
  * Hardcoded keyTable instead of having it read from file
  * Updated some packets from NA client to KR/BSW client
  * Re-did some structs
  * Added some string utils

* Dec 6, 2018
  * Uploaded with initial edit from 11/12/2018