# Use cases

Bar chart of today's currency value against USD, in logarithmic scale
```
curl -s http://api.fixer.io/latest?base=USD | jq -r ".rates | to_entries| \
    map(\"\(.key)\t\(.value|tostring)\")|.[]" | ./chart bar log -t "Currency value against USD"
```

Bar chart of a Github user's lines of code per language
```
USER=???
ACCESS_TOKEN=???
curl -u $USER:$ACCESS_TOKEN -s "https://api.github.com/user/repos" | \
    jq -r 'map(.languages_url) | .[]' | xargs curl -s -u $USER:$ACCESS_TOKEN | \
    jq -r '. as $in| keys[] | [.+ " "]+[$in[.] | tostring] | add' | \
    awk '{arr[$1]+=$2} END {for (i in arr) {print i,arr[i]}}' | \
    awk '{print $2, $1}' | sort -nr | ./chart ' ' bar
```

Pie chart of the most used terminal commands
```
make && history | awk '{print $2}' | ./chart
```
