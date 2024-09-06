# genie

Genie is an templatable event generator.  It can currently be used to generate predicatable structured and non structured events that can be used to mimic interesting business and operational measurements for testing, validation, and benchmarking.  The tool currently offers several templatable resources and outputs to stdout, http, and kafka.  Examples can be found int the [examples directory](examples/).

*Note: This tool is still in very early development stages and could use alot of love in several areas, especially around error handling/notification.  Though changes and enhancements are expected, the template syntax should stay relatively consistent.

## Install

### Homebrew

```
brew tap ctxsh/keg
brew install genie
```

### Docker

```
docker run -v $(pwd)/myconfigs:/etc/genie.d --rm ctxsh/genie:latest
```

## Global Options

* `--config/-c`: Specify the location of the configuration files.  It expects a directory and defaults to `./genie.d`.

## Commands

### generate [EVENT] [ARGS...]

The generate command starts up the configured event generators. An optional event can be specified to limit the generator to a single event.  If an event is not specified, the generator will run all configured events.

#### Options

* `--sink/-s`: Specify the sink that the generators will send the events to.  Default is `stdout`
* `--run-once`: Output a single event for each one of the configured generators and shut down.

## Sinks

Genie currently supports three sinks.  In it's early stages, the sinks only provide basic functionality, but in the future we fully expect to add features and functionality (both from a configuration and a reliability standpoint).  New sinks including Nats and TCP/socket client are currently being planned.  The TCP sink will be paired with event encoding to provide protocol support (msgpack and protobuf) for tools like Fluentd, Fluentbit, or Vector.

Sinks are referenced by the sink kind and name in the following format:

```
<kind>.<name>
```

As an example, if you have an http sink that uses the reference name of localhost (see full example below), you would reference the sink on the command line by using `-s http.localhost` when running the generator.

*Note: the sinks do not currently have the ability to specify the number of send workers that will be used to process the events that are sent to it.  This will be rectified in a future release*

### Stdout

The standard out sink outputs all events to stdout.  It is the default if no sink is specified.

### HTTP

The HTTP sink creates a client connection to the configured URL.  Currently though a method can be supplied, we only support the use of POST request using the event template as the payload.

#### Options

##### HTTP

HTTP sinks are defined as a map of the sink name to the sink configuration values.  Currently only `url`, `method`, and `headers` are supported - with the caveat previously stated around the `method` setting.

##### Headers

Headers are defined as a list of names and values.  Template values are allowed as header values to provide flexibility.

```yaml
headers:
  - name: Content-Type
    value: application/json
  - name: X-Request-Id
    value: <<uuid.request_uuid>>
```

#### Full Example

```yaml
sinks:
  http:
    localhost:
      url: http://localhost:8080
      method: post
      headers:
        - name: Content-Type
          value: application/json
        - name: X-Event-Type
          value: logging
        - name: X-Request-Id
          value: <<uuid.request_uuid>>
        - name: User-Agent
          value: <<list.user_agent>>
```

#### Future enhancements

Addition of TLS support, timeouts, retries, authentication, method specific configuration, and configurable connection pools.

### Kafka

Currently the kafka sink supports the basic configuration of brokers and the produce topic and produces asyncronously to the brokers.  Additional features and configurations will be added as the need arises.

#### Full Example

```yaml
sinks:
  kafka:
    analytics:
      topic: events
      brokers:
        - broker-0.analytics.strataviz.svc.cluster.local:9092
        - broker-1.analytics.strataviz.svc.cluster.local:9092
        - broker-2.analytics.strataviz.svc.cluster.local:9092

```

#### Future enhancements

Addition of TLS support, more complete configuration of producer settings, and synchronous produce.

## Event Generator

### Configuration

Event generators are configured as a list of events that include the following configuration options:

* `name`: a unique name that the event generator will be referenced by.
* `generators`: the number of generators/workers to start.  Each generator will send a new event every `intervalSeconds`.  Default is `1` generator.
* `intervalSeconds`: the interval in seconds in which to send a new event.  Default is `1.0` seconds.
* `vars`: a list of globally scoped variables.  Default is an empty list.
* `template`: the path to a template file.  It can either be relative or absolute.  The `template` configuration is mutually exclusive with `raw`.
* `raw`: specify the template string in the configuration instead of in a seperate file.  The `raw` configuration is mutually exclusive with `template`.

#### Global Variables

The global variables are configured as a list of objects which define the name and value of the variable.  Global variables do not support templating since they are passed directly to the template engine.  To create variables that reference resources or other variables,you can use the `let` statement in the template itself.

#### Full Example

```yaml
events:
  - name: planet_logs
    generators: 10
    vars:
      - name: version
        value: v1
    template: planet_logs.tmpl
  - name: nginx_logs
    generators: 1
    raw: >
      <# $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" "$http_x_forwarded_for"'; #>
      <<ipaddr.internal>> - - [<<timestamp.now_common_log>>] "<<list.method>> <<list.path>> HTTP/1.1" <<list.status_code>> <<integer_range.size>> "-" "<<list.user_agent>>" "<<ipaddr.external>>"
```

## Templating

The templating engine is similar to Jinja in both syntax and functionality.  It contains variable and/or expressions that are replaced with values when rendered.  There are are also formatting operators that provide the ability to trim/strip whitespace/newline characters from the 

### Expressions

Expressions can be rendered using the delimiters `<<` and `>>`, with the variable name or resource referenced between the delimiters.

```
<< myvar >>
<< list.mylist >>
```

#### Resources

##### Float Range

The float range resource renders random floats between two values.  The following configuration values are supported:

| Name | Description | Default |
| ---- | ----------- | ------- |
| min | The minimum number of the range. In the case of normalized and exponential distributions, it acts as a floor and clamps any generated value less than the minimum. | 0.0 |
| max | The maximum number of the range.  Like minimum, for normalized and exponential distributions, it clamps any generated value greater than the maximum. | 10.0 |
| distribution | The distribution of the generated numbers.  This can be `uniform`, `exponential`, and `normal`. | uniform |
| stddev | When distribution is set to normal, this represents the standard deviation. | (max-min)/8 |
| mean | When distribution is set to normal, this represents the mean of the values generated. | (max-min/2) |
| rate | When distribution is set to exponential, this represents the rate of occurrences in a given time interval | 1.0 |
| format | The expected format for the string representation of the floating point number.  Supports: <ul><li>`none`: no exponent</li><li>`decimal`: decimal representation 'e'</li><li>`decimal_capitalize`: decimal representation 'E'</li><li>`large`: 'e' for large exponents</li><li>`large_capitalize`: 'E' for large exponents</li><li>`binary`: a binary exponent</li><li>`hex`: hexidecimal fraction and binary exponent</li></ul> | none |
| precision | The number of digits (excuding the exponent) that will be printed | 5 |  


##### Integer Range

The integer range resource renders a random value between 2 numbers.  The following configuration values are supported:

| Name | Description | Default |
| ---- | ----------- | ------- |
| min | the mininimum number of the range | 0 |
| max | the maximum number of the range | 10 |
| pad | left pad `0`'s to the number when it is accessed by the template | 0 |
| distribution | the distribution of the generated integers (uniform, normal). | uniform |
| stddev | ***(normal distribution)*** sets the size of a standard deviation for a range that is normally distributed | (max-min)/10 |
| mean | ***(normal distribution)*** sets the mean for the range that is normally distributed | (max-min)/2 |

##### IP Address

The IP address resource renders a random IP address with the specified CIDR ranges.  The following configuration values are supported:

| Name | Description | Default |
| ---- | ----------- | ------- |
| cidrs | a list of IP ranges in valid CIDR format | 192.168.0.0/16 |
| uniques | create a list of ip addresses up to the number of specified unique values to use when rendering | 0 (all IP address within the range) |

##### List

The list resource renders a random value from the provided list.

##### Random String

The random string resource renders random strings based off of the provided character lists.

| Name | Description | Default |
| ---- | ----------- | ------- |
| size | the length of the random string | 10 |
| chars | the characters that will be used in generation (can be any combination of characters provided or the `alphanum`, `alpha`, `numeric`, or `hex` shortcuts) | `alphanum` |
| uniques | create a list of random strings up to the number of specified unique values to use when rendering | 0 |

##### Timestamp

The timestamp resource renders formatted timestamps.

| Name | Description | Default |
| ---- | ----------- | ------- |
| format | the format of the rendered timestamp (unix, unixnano, rfc3339, rfc3339nano, rfc1123, rfc1123z, or common_log) | rfc3339 |
| timestamp | a specific timestamp value to be rendered | now |

##### UUID

The UUID resource renders uuid1 or uuid4 formatted unique identifiers.

| Name | Description | Default |
| ---- | ----------- | ------- |
| type| the type of uuid to generate (uuid1, uuid4) | uuid4 |
| uniques | limit the number of generated uuids to the number provided | 0 |

#### Filters

***Work in progress.***

Filters can be added to variables using the `|` pipe character.

```
<< log_level|capitalize >>
```

### Assignments

To assign values to variables inside code blocks you can use the `let` statement.  Because we do not yet support scoped statements such as conditionals, loops, or functions the scope of the variable is available across the entire template.  In the future, the variable in the outer scope will be available to any inner scopes inside blocks or loops.  Attempts to set the variables will result in shadowing and will not be retained outside of the inner scope.

To assign a resource or variable to a new variable, the `let` statment is used:

```
<% let planet = list.planets %>
<% let ocean = "pacific" %>
<% let body_of_water = ocean %>

<< planet >> << ocean >> << body_of_water >>
```

## Example Template

```
<% let planet = list.planets %>
{
  "message": "<<list.messages>>",
  "aws_az": "<<list.aws_region>><<list.aws_zone>>",
  "aws_region": "<<list.aws_region>>",
  "container": {
    "id":"<<random_string.container_id>>",
    "name":"<<list.left_names>>-<<list.right_names>>"
  },
  "host": "<<planet>>-main-<<random_string.host_id>>",
  "k8s": {
    "arch":"<<list.arch>>",
    "node":"ip-<<ipaddr.internal>>.ec2.internal",
    "pod":"<<planet>>-main-<<random_string.ten>>-<<random_string.five>>"
  },
  "log_level": "<<list.log_levels>>",
  "service": "<<planet>>",
  "service_instance": "<<environment>>",
  "ts":"<<timestamp.now_unix>>",
  "<<planet>>":{
    "request_id":"<<uuid.request_uuid>>",
    "duration":"<<integer_range.duration>>",
    "method":"<<list.method>>",
    "path":"<<list.path>>",
    "user_agent":"<<list.user_agent>>"
  }
}
```
