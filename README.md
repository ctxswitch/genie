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

Add test mode for config testing and exit on config issues.
Allow custom tags/objects (use pongo2 as an example).
Once that is done, we allow tags to be custom functions as well

We need to be able to register encoders/transformers and sinks as well.  Question to all encoders start out by getting json bytes?  Might be an easy way to start.

Rethink configs for resources?


Can we add a minimize filter for text blocks?

Low hanging fruit for filtering text sections
minimize endminimize
replace endreplace // lol, would this replace variable names?  no. it doesn't impact things in
a control.

A bit harder but still decently easy
line control characters
loop endloop
filter endfilter
include
mathmatical expressions

Hard
if endif
macro endmacro
import

Ints and bools - bools come in handy with control structures.
<% let value = 100 %>
<% let value = 101.0 %>
<% let value = true %>
<% let value = false %>

Built in functions
abs()
map() // used with tojson? see attr, map should just be a resource - or can
      // we convert string to map?
select()
unique() // don't know if I know this one.
attr() // this would imply that we had some sort of a map resource.  could
       // this perhaps be resource maps defined in configs instead of using
       // json text in the event?  OOOOh, I like this.
max()
upper()
batch() // groups resources, but it's kind of a duplicate for join
min()
urlencode()
capitalize()
sort()
urlize() // return a map from url
int()
wordcount()
default() // maybe?
reject()
striptags()
wordwrap()
dictsort()
join()
groupattr()
acceptattr()
rejectattr()
sum()
escape()
replace()
length()
reverse()
tojson()
round()
trim()
float()
lower()
safe()
truncate()

<% minimize %>{ "hello": "world" }<% endminimize %>
<% replace old="hello" new="world" %>Hello world<% endreplace %>


Tomorrow:
Get let wired up, redo the whitespace control tokens.
I just finished up a rudimentary delim identifier (should refactor) but there is an
infinite loop

Requirements
- The functions will need to modify text during runtime since there can be dynamic data between the function blocks.  This is a bit of a performance hit but I think we are ok for now.  We can refactor for performance later.  These functions are applied during template execution.
- Treat lstrip and trim as functions that are wrapped around a text block.
- I need to delimit text blocks per line.  It will make it easy to wrap them.

Question? Why wouldn't we just enable trim by default and get rid of lshift?  Seems like we wouldn't really care about whitespace control with what we want to do.  It never really made sense in the templates anyway except that you ended up with values all spread out...
  - to remove this behavior, we'd just skip a single new line after the block close if it exists, then there's nothing else to do.  lshift just seems pretty ineffective at this point.