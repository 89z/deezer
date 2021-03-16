import sys
MAGIC = b'PLTE'

if len(sys.argv) != 2:
   print('android_gw_key.py [path-to-icon2.png]\n')
   exit()

fpath = sys.argv[1]
icon_buf = bytearray(open(fpath, 'rb').read())
magic_idx = icon_buf.find(MAGIC)

if magic_idx < 0:
   print('Magic phrase not found!')
   exit()

crypted = icon_buf[magic_idx:]
out = bytearray(300)

for i in range(len(out)):
   pm = i + len(MAGIC)
   if i >= 90:
      out[i - 90] = crypted[pm] ^ 0x40
   crypted[pm] = int((crypted[pm] + crypted[pm + 1] + crypted[pm + 2]) / 3)

print(out[1:17].decode('utf8'))
