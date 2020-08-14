
# Countdown Word Finder

_NOTE: Work in progress / Experimentation.
Don't make assumptions based on this code, it's a baseline for further iteration.
Furthermore, WHILE THIS NOTE PERSISTS, I MAY FORCE PUSH TO MASTER_

This is a Go demo project, which solves the problem: _"given a string of letters, which
words can be formed by using each letter from the string at most once."_

This word-game is played on the British TV panel show "8 Out Of 10 Cats... Does Countdown",
hence the title of the repo. On the show, the contestants get to pick if they want a random
vowel or a consonant, and do so for nine letters in total. They then compete by trying to
form the longest possible word out of those letters.

Validity of a word is decided by a dictionary, nothing else. There is no language processing.

All code is provided under the [MIT License](LICENSE).

[![Build Status](https://travis-ci.org/eloj/countdown.svg?branch=master)](https://travis-ci.org/eloj/countdown)

## Example

```bash
$ ./words-server &>server.log &
$ curl -s localhost:8080/countdown/v1/words/LESRNOXIS | json_pp
```

```json
{
   "duration" : 0.069932,
   "max_dist" : 4,
   "min_dist" : 1,
   "num_checked" : 42,
   "num_hits" : 5,
   "query" : "lesrnoxis",
   "words" : [
      "ironless",
      "lesions",
      "lesson",
      "snores",
      "roses"
   ]
}
```

## To Do

* Swagger API spec
* Example http client to go with server
* golint
* Docker packaging?
