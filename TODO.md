# TODO
* [] Implement resources
  * [x] Lists
  * [x] Integer range
  * [x] Random string
  * [x] UUID
  * [x] Timestamp
  * [] Map
* [x] Implement stdout

* [x] Implement config parsing
* [] Test mode for config testing

* [] Allow integer ranges to use gaussian distribution to output numbers between the min and max

* [] Implement map to json transformer

* [] Sinks can be defined for each event
* [] Implement http/s sinks
* [] Implement kafka sink

* [] Support integers, floats, and boolean values
* [] Expressions support for values other than strings
* [] Add for loop with variable scoping
* [] Add if/elif/else conditionals

* [] Macro support (macro/endmacro)
* [] Snippet support (import)

* [] Custom tags/objects
* [] We need to be able to register encoders/transformers and sinks.  Question to all encoders start out by getting json bytes?  Might be an easy way to start.

### Up next
* [] Configs are completely refactored, but now we aren't outputting.  I'm thinking that we keep sinks as a string in the event and then lookup in the manager.  Maybe decouple resources the same way?  Not sure yet though.  At the minimum I need to clean up the event parser.
* [] Maybe revisit Events again, not sure.  More I look at it, it's feeling like it's too coupled.
* [] vars are not allocated automatically before execute, so panic ensues.
---
* [] Parse the events and event/name args
* [] Add global config values like config location
* [] Clean up Generate.RunE
* [] Add sink overrides
* [] Run once impl
* [] Unknown lists are returning empty with no warnings.  Should probably have warnings during compile at least.

Notes:
* Sinks should be able to pool requests from multiple generators.  This would mean that sinks would be shared across all events. 
* Figure out where the configs will be parsed.  I think that configs will probably be parsed in Root.Execute, or we move the bulk of the config parsing and resource construction out of the root/generate where it is now.  Consider some sort of a builder package that takes a config and returns resources, sinks, templates, and all.
* Need to start on the user docs soon.
* Review testing for commands.
* For the new IP resource, build out database (whether it's csv or sqlite) and use some sort of weighted random selection based on the number of ips in a subnet, then random from the min/max of the range.