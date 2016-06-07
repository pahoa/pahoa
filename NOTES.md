# notes

tool reads all stories from pivotal and list the ones in start mode

in case user clicks to finish story, tool is going to search for branches that
match user story's id and will create a merge request to develop



## board

- todo                      # list todo, reject
- in development            # start story
- waiting for code-review   # submit for code-review (merge request to develop)
- in code-review            #
- develop in test           # deploy to develop
- waiting for qa            # finish story, merge request to qa
- in qa                     # deploy to qa
- done                      # approve

## actions

- list todo               # pivotal
- start story             # pivotal
- submit for code-review  # gitlab
- deploy to develop       # gitlab + jenkins
- finish story            # pivotal
- delivery story          # pivotal
- deploy to qa            # gitlab + jenkins
- approve                 # pivotal
  - deploy to production  # gitlab + jenkins
- reject                  # pivotal
  - rollback qa           # git

## state machine

- todo -> in development:
  - set pivotal story to `started`
- in development -> todo:
  - set pivotal story to `not started`
- in development -> waiting for code-review:
  - create merge request from feature branch to develop
- waiting for code-review -> in development:
  - close merge request open from feature branch to develop
- waiting for code-review -> in code-review:
  - nothing
- in code-review -> waiting for code-review:
  - nothing
- in code-review -> waiting for qa:
  - accept merge request from feature branch to develop
  - run jenkins build and deploy
  - set pivotal story to `finished`
  - create merge request from feature branch to qa
- waiting for qa -> in development:
  - close merge request from feature branch to qa
  - set pivotal story to `started`
- waiting for qa -> in qa:
  - 

## database

```go
type Card struct {
  Hash            string
}

type Activity struct {
  CardID      int
  From        string
  To          string
  StartedAt   time.Time
  CompletedAt time.Time
}

type 
```

- card
  - `hash`
  - `status`
  - `previous_status`
  - 
- task
  - `card_id`
  - ``

 

## state machine

create a state machine to validate move cards on the board.

each move triggers actions to be executed.

in case user wants to see the detail of a card they will be able to see the
current position of the card and the list of actions in queue, executing and
executed.

