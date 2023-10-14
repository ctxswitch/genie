# TODO

Tomorrow, work on import template function, replace filter, and start on the IP resource.  At the least, I'd like to get in a basic valid ip resource up and running.

### Commands
* [] Add test subcommand

### Resources
* [] Allow integer ranges to use gaussian distribution to output numbers between the min and max.  We'll need this for generating events that could be picked up through anomally detection in another project that is being worked on.
* [] The timestamp package can now move into the resource since config is handled there.

### Template
* [] Support integers, floats, and boolean values (that will add the ability to do arithmetic)
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
* TODO

### Notes
* Need to start on the user docs soon.
* Review testing for commands.
