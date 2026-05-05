# $\color{MidnightBlue}\textit{\textbf{ccnv}}$


![GO](https://img.shields.io/badge/Go-purple?logo=Go)


Currency Converter (ccnv) is a lightweight utility that performs currency conversion between international currencies. This project is a Golang written replica of [*C++-based ccnv utility implementation*](https://github.com/victorkryz/ccnv)

#### What the app does:
Internally, the app uses a REST API to fetch exchange rates from a remote financial service [*Free Currency Exchange Rates*](https://github.com/fawazahmed0/exchange-api).  
The utility prepares and sends an HTTP request to the remote service and parses the obtained JSON response to extract the relevant exchange rate.

### Command line arguments:

```
-l, --list     list all available currencies
-f, --from     currency convert from (usd, eur, bgn, uah, etc.)
-t, --to       currency convert to (usd, eur, bgn, uah, etc.)
-a, --amount   amount (10, 50, 100, etc.)
-v, --version  print version
-h, --help     print usage
```


### Samples of usage:

```
ccnv -l 
ccnv -f eur -t usd
ccnv -f usd -a 10 -t eur
ccnv -f usd -a 25 -t uah
ccnv --from bgn --amount 200 --to uah
```


### How to build:
-------------------------------------------------------------------------

```
    go build
```  
or

```
    go build -o ccnv
```


