``` mermaid
sequenceDiagram
    Master-->>Channel: "There are tasks to do!"
    Channel-->>Worker1: *listens*
    Channel-->>Worker2: *listens*
    Worker1-->>Master: "I am available!"
    Master->>Worker1: task name and context
    Worker1-->+Worker1: store context start task...
    Worker2-->>Master: "I am available!"
    Master-->>Worker2: "There are no tasks to do!"
    Worker1->>-Master: done: task name, prev. context, new context
```