# chart
Quick & smart charting for STDIN

[![Build Status](https://img.shields.io/travis/MarianoGappa/chart.svg)](https://travis-ci.org/MarianoGappa/chart)
[![Coverage Status](https://coveralls.io/repos/github/MarianoGappa/chart/badge.svg?branch=master&nocache=1)](https://coveralls.io/github/MarianoGappa/chart?branch=master)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/MarianoGappa/chart/master/LICENSE)

## Syntax

```
chart [pie|bar|line|scatter|log|' '|';'|','|'\t'|-t %title%|--title %title%|-x %x_axis_label%|-y %y_axis_label%]
```

- `pie`: render a pie chart
- `bar`: render a bar chart
- `line`: render a line chart
- `scatter`: render a scatter plot chart
- `log`: use logarithmmic scale (bar chart only)
- `' '|';'|','|'\t'`: this character separates columns on each line (\t = default)
- `-t|--title`: title for the chart
- `-x`: label for the x axis
- `-y`: label for the y axis

## Installation

```
go get github.com/MarianoGappa/chart
```

## Example use cases

- Pie chart of your most used terminal commands
```
history | awk '{print $2}' | chart
```

![Pie chart of your most used terminal commands](img/pie.png?v=1)

- Bar chart of today's currency value against USD, in logarithmic scale
```
curl -s http://api.fixer.io/latest?base=USD | jq -r ".rates | to_entries| \
    map(\"\(.key)\t\(.value|tostring)\")|.[]" | chart bar log -t "Currency value against USD"
```

![Bar chart of today's currency value against USD, in logarithmic scale](img/bar-log.png?v=1)

- Bar chart of a Github user's lines of code per language (requires setting up an Access Token)
```
USER=???
ACCESS_TOKEN=???
curl -u $USER:$ACCESS_TOKEN -s "https://api.github.com/user/repos" | \
    jq -r 'map(.languages_url) | .[]' | xargs curl -s -u $USER:$ACCESS_TOKEN | \
    jq -r '. as $in| keys[] | [.+ " "]+[$in[.] | tostring] | add' | \
    awk '{arr[$1]+=$2} END {for (i in arr) {print i,arr[i]}}' | \
    awk '{print $2 "\t" $1}' | sort -nr | chart bar
```

![Bar chart of a Github user's lines of code per language (requires setting up an Access Token)](img/bar.png?v=1)

- MySQL query output charting

TODO

## Details

- `chart` is still experimental.
- it infers STDIN format by analysing line format on each line (doesn't infer separator though; defaults to `\t` and accepts user overrides) and computing the winner format.
- it uses the awesome [ChartJS](http://www.chartjs.org/) library to plot the charts.
- when input data is string-only, `chart` infers a "word frequency pie chart" use case.
- should work on Linux/Mac/Windows thanks to [open-golang](https://github.com/skratchdot/open-golang).

## Contribute

Yes, please.
