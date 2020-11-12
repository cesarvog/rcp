# rcp

Copy and paste **strings** across internet

---

This program was design to be used in a shell.

Some user cases:

### Copy simple string and pasting

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

<p align="center"><img src="imgs/rcpsimple.gif?raw=true"/></p>

----

