<img src="https://img.shields.io/badge/lifecycle-experimental-orange.svg" alt="Life cycle: experimental"> <a href="https://godoc.org/github.com/openbiox/bioctl"><img src="https://godoc.org/github.com/openbiox/bioctl?status.svg" alt="GoDoc"></a>

## bioctl

[bioctl](https://github.com/openbiox/bioctl) is a simple command line tool to facilitate the data analysis.

## Installation

```bash
# windows
wget https://github.com/openbiox/bioctl/releases/download/v0.1.0/bioctl.exe

# osx
wget https://github.com/openbiox/bioctl/releases/download/v0.1.0/bioctl_osx
mv bioctl_osx bioctl
chmod a+x bioctl

# linux
wget https://github.com/openbiox/bioctl/releases/download/v0.1.0/bioctl_linux64
mv bioctl_linux64 bioctl
chmod a+x bioctl

# get latest version
go get -u github.com/openbiox/bioctl
```

## Usage

### Simple parallel tasks

```bash
# concurent 2 tasks with total 8 tasks
echo 'echo $1 $2; sleep ${1}' > job.sh && bioctl par --cmd "sh job.sh" -t 2 --index 1,2,5-10
# concurent 4 tasks with total 8 tasks and env parse
bioctl par --cmd 'sh job.sh {{index}} {{key2}}' -t 4 --index 1,2,5-10 --env "key2:123"

# concurent 4 tasks with total 8 tasks (direct) and env
bioctl par --cmd 'echo {{index}} {{key2}}; sleep {{index}}' -t 4 --index 1,2,5-10 --env "key2:123"

# concurent 4 tasks with total 8 tasks, env
# and not to force add {{index}} at the end of cmd
bioctl par --cmd 'echo {{key2}}; sleep {{index}}' -t 4 --index 1,2,5-10 --env "key2:123" --force-idx false

# pipe usage
echo 'sh job.sh' | bioctl par -t 4 --index 1,2,5-10 -
```

### Format function

```bash
# format json string
echo '{"a":1, "b":123}' | bioctl fmt --json-pretty -

bget api ncbi -q "Galectins control MTOR and AMPK in response to lysosomal damage to induce autophagy OR MTOR-independent autophagy induced by interrupted endoplasmic reticulum-mitochondrial Ca2+ communication: a dead end in cancer cells. OR The PARK10 gene USP24 is a negative regulator of autophagy and ULK1 protein stability OR Coordinate regulation of autophagy and the ubiquitin proteasome system by MTOR." | bget api ncbi --xml2json pubmed - | sed 's;}{;,;g' | bioctl fmt --json-to-slice --indent 4 -
```

### Simple file stat

```bash
# file stat
# equal to wc -l
bioctl fn -l *
```

## Maintainer

- [@Jianfeng](https://github.com/Miachol)

## License

Academic Free License version 3.0

