# Sayit Go
Pronouncing given word using Oxford Advanced Learner's Dictionary website.  

``` bash
# Ubuntu
sudo apt-get install portaudio19-dev libmpg123-dev
make && make install
```

run from command line:
``` bash
sayit check hello
```
or just copy as `sayit` somewhere in your path and use it with goldendict:
- Edit -> Dictionaries -> Programs.
- choose "Audio" in type field.
- in "Command Line" field copy and paste this command:
``` bash
sayit check %GDWORD%
```
- write "sayit" in "Name" field.

