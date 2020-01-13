# envsubsty
[![Go Report Card](https://goreportcard.com/badge/github.com/kukaryambik/envsubsty)](https://goreportcard.com/report/github.com/kukaryambik/envsubsty)

The envsubsty converts the specified environment variables in files to their value.

### Usage
```
envsubsty [-wh] [-v 'vars'] [file|directory ...]
```
Or
```
cat file.txt | envsubsty [-v 'vars']
```
Flags:
 - `-h` - Show help message.
 - `-v 'string'` - Comma or space-separated list of variables to convert.
 - `-w` - Write the output to the source file.
