# Slogger

Slogger enables services to gain consistent logging formats throughout environments. 
Under the hood it wraps Go's standard library structured logging package `log/slog.`

Slogger exposes functions to initialise new `slog.Handler`'s with appropriate styling for the environment of the service.

On top of styling logger output for services, slogger also provides a local development slog handler which helps to provide "nicer on the eye"
formatting for local development.

## Usage 

Note: Services will need to map an environment variable to the appropriate `slogger.Env` type.

### Basic usage (Returns a new `slog.Logger` ready for DEV env):

```go
    slogHandler := slogger.NewHandler(slogger.EnvDev)
    logger := slog.New(slog.Handler)
```
### Usage with config (Overriding the style):

```go
    dev := slogger.NewHandler(slogger.EnvDev, slogger.WithStyle(slogger.StyleJSON))
    logger := slog.New(slog.Handler)
```

## Running Examples

```sh 
make run-ex LOGGER_STYLE=json
```


### Example Local Development output

```
[12:15:25 GMT] This is an info level log from local
  foo: 1
  bar: 2
```
