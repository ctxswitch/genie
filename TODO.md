# TODO

### Commands
* [] Parse the events and event/name args
* [] Add global config args for config locations
* [] Add run-once
* [] Add test-config
* [] Clean up Generate.RunE

### Resources
* [] Allow integer ranges to use gaussian distribution to output numbers between the min and max
* [] New IP resource, build out database (whether it's csv or sqlite) and use some sort of weighted random selection based on the number of ips in a subnet, then random from the min/max of the range.  User would be able to specify region(s).
* [] The timestamp package can now move into the resource since config is handled there.

### Template
* [x] Variable refactor/extraction
* [x] New system for variable scoping
* [x] Introduce scoped variables to templates.  They won't come into play until loops/conditions
* [] Support integers, floats, and boolean values
* [] Expressions support for values other than strings
* [] Add for loop with variable scoping
* [] Add if/elif/else conditionals
* [] Macro support (macro/endmacro)
* [] Snippet support (import)
* [] Custom tags/objects - this I'm not sure of right now, as it'd probably just be useful in a few edge cases.

### Sinks
* [] Sinks can be defined for each event
* [] Implement kafka sink
* [] Implement backoff for HTTP(s) sinks
* [] Add sink overrides from the command.  Not a big deal right now, but we could use that for sending to stdout or another test endpoint while a developer is building/testing the events.
* [x] Maybe (leaning heavily to this): Sinks should be able to pool requests from multiple generators.  This would mean that sinks would be shared across all events.  This would make sense from a resource perspective as to not overwhelm an endpoint, but that also means that the sinks would be independent and we run them independently from the generators.
* [] Make connection pool workers for sinks that share the client (http specifically), stdout doesn't really want or need it.

### Encoding
* [] Create encoders for the events.  Still don't have them planned out, but they would allow custom/builtin encoding i.e. json bytes, msgpack, etc.
* [] statsd encoder
* [] fluentbit specific encoder

### Maps
* [] Line out the purpose of maps within the current structure.  At this point I'm not sure that I need them.
* [] Implement maps.
* [] Implement a map to json string filter.

### Fixes
* [] vars are not allocated automatically before execute, so panic ensues.
* [] Unknown lists are returning empty with no warnings.  Should probably have warnings during compile at least.
* [] Unknown lists are returning empty with no warnings.  Should probably have warnings during compile at least.

Notes:

* Need to start on the user docs soon.
* Review testing for commands.
