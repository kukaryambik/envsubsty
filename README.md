# envsubsty
[![Go Report Card](https://goreportcard.com/badge/github.com/kukaryambik/envsubsty)](https://goreportcard.com/report/github.com/kukaryambik/envsubsty)
[![Github Release](https://img.shields.io/github/release/kukaryambik/envsubsty.svg)](https://github.com/kukaryambik/envsubsty/releases)

The envsubsty converts the specified environment variables in files to their value.

Unlike the classic envsubst, this application might convert complex variables like `${FOO:-${BAR:-value}}`.

### Usage
```bash
envsubsty [-hVwe] [-v 'vars'] [file|directory ...]
```
Or
```bash
cat file.txt | envsubsty [-v 'vars']
```
Flags:
 - `-V` - Show version.
 - `-h` - Show help message.
 - `-v 'string'` - Comma or space-separated list of variables to convert.
 - `-w` - Write the output to the source file.
 - `-e` - Convert empty variables.
