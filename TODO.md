# Use cases

curl -s http://api.fixer.io/latest?base=USD | jq -r ".rates | to_entries|map(\"\(.key)\t\(.value|tostring)\")|.[]" | ./chart bar log -t "Currency value against USD"
