# SoulWorker Server
My slow attempt at an 'RE' of the soulworker protocol.

No database (yet)

Character creation is iffy, some parts are still broken

Changelog:

* Dec 6, 2018
	* Uploaded with initial edit from 11/12/2018

* July 13, 2019
	* Hardcoded keyTable instead of having it read from file
	* Updated some packets from NA client to KR/BSW client
	* Re-did some structs
	* Added some string utils
	* Note: Crashes on world join with an "invalid movie" error. (Didn't go through all the packets so some still need to be edited)

* July 29, 2019
	* Moved hardcoded keyTable to globals where it's defined
	* Minimal file updates, just slight renaming on a few variables
	* Changed cringe server name
	* Added some tools/tests for `tb_*.res` table extraction/verification
		* Includes `tb_*` string dump and the ASM dump for functions regarding them
		* Includes the final go struct dump for all table structures
		* You can find the original script that I modified in this [Gist](https://gist.github.com/x1nixmzeng/a4a5c419f1cd4bc72cba30d5e647bc4f)

How to use:

1. Download and open in GoLand or whatever IDE (or not) you use
2. Compile & Run
3. Enjoy up to character creation

For windows networking mitm/proxy:

Before client launch, run in cmd:
```
netsh int ip add addr 1 KOREAN_SERVER_IP/32 st=ac sk=tr
```
After you're done, run in cmd:
```
netsh int ip delete addr 1 KOREAN_SERVER_IP
```

You can get the IP through wireshark or resource monitor or any various tool.