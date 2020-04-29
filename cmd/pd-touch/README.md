# pd-deal

Command-line deal printer

## Sample invocations

* print value of all orders
```
pd-deal -token $PDTOKEN -filter 231 -template '{{.Value}}'
pd-deal -token $PDTOKEN -filter 231 -template '{{.Value | printf "%0.f"}}'
```

* print value of all deals with offers
```
pd-deal -token $PDTOKEN -filter 305 -template '{{.Value}}'
pd-deal -token $PDTOKEN -filter 305 -template '{{.Value | printf "%0.f"}}'
```



