Packagelock is an interactive Serversoftware that shows all packages on your servers collected by agents.  

## Backend
the go based Backend provides a JSON-REST API for the frontend.  

## Frontend 
the TypeScript based frontend is a Single Page Application that uses the REST backend to display the data. it uses patternfly components (RedHat UI tool-Kit).


feature creep:
- timed pooling from agents to '/pool/' to get commands like rescans or updates

