# Service base

This repository contains base components for our backend services.

# Usage

```go
// Any http.Handler compatible struct
mux := http.NewServerMux()
// Create http component
http := application.NewHttpComponent(router)
// Compose new application
app := application.NewApplications(http)
// Run blocks until application exits
err :=  app.Run()
if err ! nil{
    return err
}
```

# Structure

Each compoonent is composed of 3 functions, that will be invoked by the `Application`
- Startup
- Run
- Close

### Lifecycle
```mermaid
---
title: Lifecycle
---
stateDiagram-v2
    state "Application starting" as AS
    state "Application running" as AR
    state is_c_running_state <<choice>>
    state is_a_running_state <<choice>>
    
    [*] --> AS: Application.Run()
    AS --> Component
    state Component {
        [*] --> Starting: Component.Startup()
        Starting --> Running: Component.Run()
        Running --> is_c_running_state
        is_c_running_state --> Running
        is_c_running_state --> Closing: Error
        Closing --> Done: Component.Close()
    }
    AS --> AR: All components has started
    AR --> is_a_running_state
    is_a_running_state --> AR
    is_a_running_state --> Closing: Signal to shutdown
    Done --> [*]
```
