# sw_extractor.py
# written by WRS/x1nixmzeng
# modified by Austin for GOLang and aditional data types
#######################################################################
# usage: 
# 1: save all string references for tb_* files to "soul_worker_res_types.txt"
# 2: label these functions (there are more, but I don't remember the exact addresses. These are his. Renamed with IDA. read_wstring_and_len, read_float32, read_(u)int32, read_(u)int16, read_uint8, byte)
#  a) read_int - 012643A0
#  b) read_byte - 012642A0
#  c) read_short - 01264320
#  d) read_wstring_and_len - 0125E0C0
#  e) read_tb - 0125DBF0
#  f) read and verify checksum - 01286D00
# 3: dump the disassembly output to "soul_worker_asmdump.txt"
# Note (I used Scylla + IDA to dump the ASM) - Austin
# 4: run "sw_extractor.py > structs.txt"
#######################################################################

# read all resource types to dump
fn_resources = 'SW_string_dump.txt'

# read the labelled assembly dump
fn_assembly = 'SW_dump_SCY.asm'

with open(fn_resources, encoding="utf8") as f:
    res_list = f.readlines()

with open(fn_assembly, encoding="utf8") as f:
    asm_list = f.readlines()
	
str_ascii = "; "

for n, line in enumerate(res_list, 1):
	line = line.rstrip()
	addr = line[line.find("\t")+1:line.find("\t ")]
	asciipos = line.find(str_ascii);
	name = line[asciipos+len(str_ascii)+1:len(line)-1]

	found = 0
	
	types = {}
	
	for m, line2 in enumerate(asm_list, 1):

		if found == 0 :
			# look for our address
			if line2.find(addr) != -1:
				found = 1
		elif found == 1:
			# look for the file read
			if line2.find("read_table") != -1:
				found = 2
				print("type " + ("_".join((" ".join(name.split("_"))).title().split(" "))) + " struct {")
		elif found == 2:
			# try to early out
			if line2.find("read_and_verify_checksum") != -1:
				print("}")
				found = 0
				break
			else:
				# look for read calls
				pos = line2.find("read_")
				if pos != -1:
					type = line2[pos+5:len(line2)-1]
					if type not in types:
						types[type] = 0
					types[type] += 1
					cnt = types[type]

					if type != "wstring_and_len":
						print("\t" + type.capitalize() + "_" + str(cnt) + "\t" + type)
					else:
						print("\tLen" + str(cnt) + "\tuint16")
						print("\tString" + str(cnt) + "\tstring")