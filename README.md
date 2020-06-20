## aliengame
This is a game about mad aliens taking over the world. They mysteriously appear in various cities, walk to neighboor towns and fight with each other because they are so mad.
Once they start fighting, the whole city teleports to another world in the universe and disappears from our big blue bubble.

### Installation
```
$ go get -u github.com/ilgooz/aliengame/interface/alienctl
```

### Usage
```
$ alienctl --help
```

### Game Logic 
* A world is created with cities by the given map.
* N number of aliens are spawned at random cities.
* After each resume _(iteration)_ in the world, all aliens walks to the neigbor cities that have a direct path to the current city and fight with each other if there are multiple aliens in the city.
* After fought, the city and aliens on the city is removed from the game. Also any other cities that are neighbor of the gone city updated to destroy paths _(directions)_ to the gone city.
* World is continously resumed until no aliens left or each living alien has walked _10000_ times.

#### Details 
* Multiple aliens might spawn in the same city but they won't fight until the first world.Resume() _(iteration)_.
* The game runs in a single _goroutine_ because it is assumed that all aliens move at the same time and they fight at the same time. This behavior is chosen to reduce implementation complexity.
* A city hosts N _(>=0)_ number of aliens at a time.

### Project Stucture
```
.
├── aliengame                       -> source code of the game
│   ├── aliengame.go                     
│   ├── aliengame_test.go               
│   ├── event.go                    
│   ├── event_test.go                    
│   ├── map.go               
│   └── map_test.go                 
├── go.mod
├── go.sum
├── interface                       -> network/user interfaces to expose the game
│   └── alienctl                    -> cli for the game
│       ├── cmd                     -> reusable cmd for the game
│       │   ├── game.go
│       │   └── game_test.go
│       ├── main.go
│       └── main_test.go
├── README.md
├── mapdata                         -> premade maps for the game
│   └── 0.aliengame
└── x                               -> util like but generic, reusable packages
    └── compass                     -> a compass package for city directions
        ├── compass.go
        └── compass_test.go
```
