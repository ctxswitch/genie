# genie

Event generator (more than just events... request/payload?)


## Global Options

* `--prefix/-p`: set the path prefix, i.e. the directory where all generator files will be kept.  Default value is `./genie.d`.
* `--config/-c`: specify an alternate configuration file instead of the default (`<prefix>/config.yaml`).

### generate

* `generate events`: start generators for all configured events
* `generate event/...`: start generator for a single event
* `generate event ...`: start generator for a single event

#### Options
* `--sink/-s`: override the configured sink.
* `--count=<num>`: output num events and shut down.

### test

test


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

### generate

Starts the generators for one or many resources.

genie generate all (could be the default?)
genie generate raw_logs