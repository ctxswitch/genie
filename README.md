# dynamo

Event generator

```
~/.dynamo/all.yaml
```

```
dynamo --config=path/to/config.yaml generate event/id
dynamo --config=path/to/config.yaml generate events
dynmao --dir=path/to/dir generate events
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
TODO: Loops?