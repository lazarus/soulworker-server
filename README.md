# SoulWorker Server
My slow attempt at an 'RE' of the soulworker protocol.

No database (yet). Character creation is iffy, some parts are still broken.

Crashes on world join with an "invalid movie" error. (Didn't go through all the packets so some still need to be edited)

How to use:

1. Download and open in GoLand or whatever IDE (or not) you use
2. Compile & Run
3. Enjoy up to character creation

For windows networking mitm/proxy:

Before client launch, run in cmd:
```
netsh int ip add addr 1 KOREAN_SERVER_IP/32 st=ac sk=tr
```
After you're done, run in cmd and launch the client:
```
netsh int ip delete addr 1 KOREAN_SERVER_IP
```
You can get the IP through wireshark or resource monitor or any various tool.

(or potentially can just specify the IP through the command line, I forgot to test before updating this).

To launch the KRSW OnStove client (+skipping Stove Authenticator):
```
SoulWorker.exe Live/KOREAN_SERVER_IP/10000 SkipSGAuthen:yes
```

To bypass XIGNCode3, you can build the dll found [here](https://github.com/austinh115/XignCode3-bypass).

To bypass the OnStove authenticator without command line arguments, you can build the dll found [here](https://github.com/austinh115/OnStove-Client).

(I'm not super good at CPP or RE so it can be done better. Some of the functions don't even end up getting called which isn't what I expected)

NOTE: You won't be able to login as the game client sees the auth code as invalid, will need to play around with the bypass return values until it works.

Changelog:

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
