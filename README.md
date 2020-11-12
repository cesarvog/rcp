# rcp

Copy and paste **strings** across internet

---

This program was design to be used in a shell.

Some user cases:

### Copying simple string and pasting

Machine A
```
rcp -c I need to transfer this text to another machine
```

Machine B
```
rcp -p
```

---

<p align="center"><img src="imgs/rcpsimple.gif?raw=true"/></p>

---

### Copying text files

Machine A
```
cat ~/package.json | rcp -c
```

Machine B
```
rcp -p >> backup/package.json
```

---

<p align="center"><img src="imgs/rcpfile.gif?raw=true"/></p>

----

### Copying binaries files

rcp was not design to work with bin files, but design to be used with other commands...
so it's possible to work with base64 command to copy bin files as follow


Machine A
```
base64 /usr/bin/cat | rcp -c
```

Machine B
```
rcp -p | base64 -d >> catcopy
```


## Configuring rcp

Before using -c and -p arguments the use of --configure is needed 
this arg will generate ~/.rcp.properties

you'll need to inform a **secret**, the secret must be a string that just you must known and big enough to no one can guesses that
you can use [https://generate.plus/en/base64](https://generate.plus/en/base64) to generate a random base64 string to use as secret

ex:

```
rcp --configure KLw6YjPniBGfcnq9pCg+81bs8idkt0+/gkbpfrPzBFs=
```

You can use **--help** for more help if this doc was not clear



