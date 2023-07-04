# genie

Event generator

```
~/.genie/all.yaml
```

```
genie --config=path/to/config.yaml generate event/id
genie --config=path/to/config.yaml generate events
genie --dir=path/to/dir generate events
```

```yaml
resources:
  lists:
    - name: foo
      values:
        - one
        - two
        - three
```

<% import this.json.tmpl | minimize %>
Built in functions|filters

let
for
if

urlencode // url encode the string - must be url
upper (capitalize) // uppercase all
lower // lowercase all
join // calls the resource multiple times to create an array of values then joins those with the selected delimiter.
env // get environment variable
escape // escape characters in string
default // use the default string as the value if the variable or resource does not exist - don't know if I need this one.
minimize

// maybe
replace
length
truncate
reverse
wordwrap
trim
