# Turbine

A minimalistic build server

## Usage
```
./turbine -p 666
```

## Api


## Todo
logging, write to file or /var/log/messages.  
log cpu, memory and bandwidth usage for each running service?  
support rest api for stop/start/status for each service.  

## Config
```
[
    {
        "name": "turbine",
        "build": "cd hook && go build",
        "run": "cd hook && ./turbine",
        "time": 0
    }
]
```